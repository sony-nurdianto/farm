package kev

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumerPool struct {
	mu        sync.RWMutex
	consumers map[string]*pooledConsumer
	kafka     Kafka
}

type pooledConsumer struct {
	consumer  KevConsumer
	lastUsed  time.Time
	eventOnce sync.Once
}

func NewKafkaConsumerPool(kk Kafka) *KafkaConsumerPool {
	pool := &KafkaConsumerPool{
		kafka:     kk,
		consumers: make(map[string]*pooledConsumer),
	}

	return pool
}

func (kc *KafkaConsumerPool) getOrCreateProducer(cfg kafka.ConfigMap) (*pooledConsumer, error) {
	key := hashConfig(cfg)

	// First check with read lock
	kc.mu.RLock()
	if pooled, exists := kc.consumers[key]; exists {
		pooled.lastUsed = time.Now()
		kc.mu.RUnlock()
		return pooled, nil
	}
	kc.mu.RUnlock()

	// Need to create new producer - acquire write lock
	kc.mu.Lock()
	defer kc.mu.Unlock()

	// Double-check pattern - another goroutine might have created it
	if pooled, exists := kc.consumers[key]; exists {
		pooled.lastUsed = time.Now()
		return pooled, nil
	}

	// Create new producer

	producer, err := kc.kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	pooled := &pooledConsumer{
		consumer: producer,
		lastUsed: time.Now(),
	}

	kc.consumers[key] = pooled
	return pooled, nil
}

func RebalanceCb(c KevConsumer, e kafka.Event) error {
	switch ev := e.(type) {
	case kafka.AssignedPartitions:
		log.Printf("🟢 Partitions assigned: %d partitions", len(ev.Partitions))
		for _, partition := range ev.Partitions {
			log.Printf("  → %s[%d]", *partition.Topic, partition.Partition)
		}

		if committed, err := c.Committed(ev.Partitions, 5000); err == nil {
			for _, tp := range committed {
				if tp.Offset >= 0 {
					log.Printf("  📍 %s[%d] committed offset: %d", *tp.Topic, tp.Partition, tp.Offset)
				}
			}
		}

		if err := c.Assign(ev.Partitions); err != nil {
			log.Printf("❌ Failed to assign partitions: %v", err)
			return err
		}

		log.Printf("✅ Assignment completed successfully")

	case kafka.RevokedPartitions:
		log.Printf("🔴 Partitions revoked: %d partitions", len(ev.Partitions))
		for _, partition := range ev.Partitions {
			log.Printf("  ← %s[%d]", *partition.Topic, partition.Partition)
		}

		if len(ev.Partitions) > 0 {
			log.Printf("💾 Committing offsets before revocation...")
			if _, err := c.Commit(); err != nil {
				log.Printf("⚠️ Failed to commit offsets: %v", err)
			} else {
				log.Printf("✅ Offsets committed successfully")
			}
		}

		if err := c.Unassign(); err != nil {
			log.Printf("❌ Failed to unassign partitions: %v", err)
			return err
		}

		log.Printf("✅ Revocation completed successfully")

	case kafka.PartitionEOF:
		log.Printf("📄 Reached end of partition %s[%d] at offset %d",
			*ev.Topic, ev.Partition, ev.Offset)

	case kafka.Error:
		if ev.Code() == kafka.ErrAllBrokersDown {
			log.Printf("🚨 Critical: All brokers are down!")
			return ev
		}
		log.Printf("⚠️ Kafka error: %v", ev)

	case kafka.OAuthBearerTokenRefresh:
		log.Printf("🔑 OAuth token refresh required")

	default:
		log.Printf("🔵 Other rebalance event: %T - %v", ev, ev)
	}

	return nil
}

func (kc *KafkaConsumerPool) Close() {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	for key, pooled := range kc.consumers {
		pooled.consumer.Close()
		delete(kc.consumers, key)
	}

	fmt.Println("🛑 Kafka producer pool closed")
}

func (kp *KafkaConsumerPool) GetPoolStats() map[string]any {
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	return map[string]any{
		"active_consumers": len(kp.consumers),
		"consumers_keys": func() []string {
			keys := make([]string, 0, len(kp.consumers))
			for k := range kp.consumers {
				keys = append(keys, k[:8]+"...") // Show first 8 chars of hash
			}
			return keys
		}(),
	}
}

func (kc *KafkaConsumerPool) setConfigMapKey(cfg map[ConfigKeyKafka]string) kafka.ConfigMap {
	configMap := make(kafka.ConfigMap)
	for key, value := range cfg {
		configMap.SetKey(string(key), value)
	}

	return configMap
}

func (kc *KafkaConsumerPool) Consumer(cfg map[ConfigKeyKafka]string) (KevConsumer, error) {
	configMap := kc.setConfigMapKey(cfg)

	pooled, err := kc.getOrCreateProducer(configMap)
	if err != nil {
		return nil, err
	}

	return pooled.consumer, nil
}
