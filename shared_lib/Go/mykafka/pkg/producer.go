package pkg

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducerPool struct {
	mu        sync.RWMutex
	producers map[string]*pooledProducer
}

type pooledProducer struct {
	producer  *kafka.Producer
	lastUsed  time.Time
	eventOnce sync.Once
}

func NewKafkaProducerPool() *KafkaProducerPool {
	pool := &KafkaProducerPool{
		producers: make(map[string]*pooledProducer),
	}

	// Start cleanup goroutine untuk remove idle producers
	go pool.cleanupRoutine()

	return pool
}

// hashConfig membuat hash dari configuration untuk digunakan sebagai key
func hashConfig(cfg kafka.ConfigMap) string {
	var keys []string
	for key := range cfg {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var configStr strings.Builder
	for _, key := range keys {
		val, _ := cfg.Get(key, "")
		configStr.WriteString(fmt.Sprintf("%s=%v;", key, val))
	}

	hash := sha256.Sum256([]byte(configStr.String()))
	return fmt.Sprintf("%x", hash)
}

func (kp *KafkaProducerPool) getOrCreateProducer(cfg kafka.ConfigMap) (*pooledProducer, error) {
	key := hashConfig(cfg)

	// First check with read lock
	kp.mu.RLock()
	if pooled, exists := kp.producers[key]; exists {
		pooled.lastUsed = time.Now()
		kp.mu.RUnlock()
		return pooled, nil
	}
	kp.mu.RUnlock()

	// Need to create new producer - acquire write lock
	kp.mu.Lock()
	defer kp.mu.Unlock()

	// Double-check pattern - another goroutine might have created it
	if pooled, exists := kp.producers[key]; exists {
		pooled.lastUsed = time.Now()
		return pooled, nil
	}

	// Create new producer
	producer, err := kafka.NewProducer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	pooled := &pooledProducer{
		producer: producer,
		lastUsed: time.Now(),
	}

	// Start event handler only once per producer
	pooled.eventOnce.Do(func() {
		go kp.handleEvents(producer, key)
	})

	kp.producers[key] = pooled
	return pooled, nil
}

func (kp *KafkaProducerPool) handleEvents(producer *kafka.Producer, key string) {
	defer func() {
		// Clean up when event loop exits
		kp.mu.Lock()
		delete(kp.producers, key)
		kp.mu.Unlock()
	}()

	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("❌ Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("✅ Delivered to: %v\n", ev.TopicPartition)
			}
		case kafka.Error:
			fmt.Printf("🔥 Kafka error: %v\n", ev)
			if ev.IsFatal() {
				// Fatal error - close producer
				producer.Close()
				return
			}
		case kafka.Stats:
			fmt.Printf("📊 Kafka stats: %s\n", ev)
		case kafka.LogEvent:
			fmt.Printf("🪵 Kafka log [%s] %s: %s\n", ev.Timestamp, ev.Name, ev.Message)
		default:
			fmt.Printf("🤷 Unknown event: %v\n", ev)
		}
	}
}

func (kp *KafkaProducerPool) Produce(cfg map[ConfigKeyKafka]string, msgs ...*MessageKafka) error {
	if len(msgs) == 0 {
		return nil
	}

	// Convert config
	configMap := kafka.ConfigMap{}
	for key, value := range cfg {
		if err := configMap.SetKey(string(key), value); err != nil {
			return fmt.Errorf("failed to set config key %s: %w", key, err)
		}
	}

	// Get or create producer from pool
	pooled, err := kp.getOrCreateProducer(configMap)
	if err != nil {
		return err
	}

	producer := pooled.producer

	// Produce messages
	for _, msg := range msgs {
		msgFactory := msg.Factory()

		// Gunakan delivery channel untuk async handling
		deliveryChan := make(chan kafka.Event)
		err := producer.Produce(&msgFactory, deliveryChan)
		if err != nil {
			return fmt.Errorf("failed to produce message: %w", err)
		}

		// Optional: Handle delivery report immediately
		go func() {
			e := <-deliveryChan
			if m, ok := e.(*kafka.Message); ok {
				if m.TopicPartition.Error != nil {
					fmt.Printf("❌ Message delivery failed: %v\n", m.TopicPartition)
				}
			}
			close(deliveryChan)
		}()
	}

	// Flush with timeout
	remaining := producer.Flush(5000) // 5 second timeout
	if remaining > 0 {
		return fmt.Errorf("failed to flush all messages, %d messages remaining", remaining)
	}

	return nil
}

// cleanupRoutine removes idle producers periodically
func (kp *KafkaProducerPool) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		kp.cleanupIdleProducers(10 * time.Minute)
	}
}

func (kp *KafkaProducerPool) cleanupIdleProducers(maxIdleTime time.Duration) {
	kp.mu.Lock()
	defer kp.mu.Unlock()

	now := time.Now()
	for key, pooled := range kp.producers {
		if now.Sub(pooled.lastUsed) > maxIdleTime {
			pooled.producer.Close()
			delete(kp.producers, key)
			fmt.Printf("🧹 Cleaned up idle producer: %s\n", key)
		}
	}
}

func (kp *KafkaProducerPool) Close() {
	kp.mu.Lock()
	defer kp.mu.Unlock()

	for key, pooled := range kp.producers {
		pooled.producer.Close()
		delete(kp.producers, key)
	}

	fmt.Println("🛑 Kafka producer pool closed")
}

// GetPoolStats returns statistics about the pool
func (kp *KafkaProducerPool) GetPoolStats() map[string]any {
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	return map[string]any{
		"active_producers": len(kp.producers),
		"producer_keys": func() []string {
			keys := make([]string, 0, len(kp.producers))
			for k := range kp.producers {
				keys = append(keys, k[:8]+"...") // Show first 8 chars of hash
			}
			return keys
		}(),
	}
}
