package pkg

import (
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaProducer struct {
	mu        sync.RWMutex
	producers map[string]*kafka.Producer
}

func NewKafkaProducer() *kafkaProducer {
	return &kafkaProducer{
		producers: make(map[string]*kafka.Producer),
	}
}

func getOrCreateProducer(
	cfg kafka.ConfigMap,
	mu *sync.RWMutex,
	producers map[string]*kafka.Producer,
) (*kafka.Producer, error) {
	key := hashConfig(cfg)
	mu.RLock()
	producer, ok := producers[key]
	mu.RUnlock()

	if ok {
		return producer, nil
	}

	prod, err := kafka.NewProducer(&cfg)
	if err != nil {
		return nil, err
	}

	producers[key] = prod
	return prod, nil
}

func produceEvents(producer *kafka.Producer) {
	go func() {
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

			case kafka.Stats:
				fmt.Printf("📊 Kafka stats: %s\n", ev)

			case kafka.LogEvent:
				fmt.Printf("🪵 Kafka log [%s] %s: %s\n", ev.Timestamp, ev.Name, ev.Message)

			default:
				fmt.Printf("🤷 Unknown event: %v\n", ev)
			}
		}
	}()
}

func (kp *kafkaProducer) Producer(cfg map[ConfigKeyKafka]string, msg ...*kafka.Message) error {
	configMap := kafka.ConfigMap{}
	for key, value := range cfg {
		if err := configMap.SetKey(string(key), value); err != nil {
			return err
		}
	}

	producer, err := getOrCreateProducer(configMap, &kp.mu, kp.producers)
	if err != nil {
		return err
	}

	defer producer.Close()
	produceEvents(producer)

	for _, v := range msg {
		producer.Produce(v, nil)
	}

	producer.Flush(15 * 1000)

	return nil
}

func (p *kafkaProducer) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, prod := range p.producers {
		prod.Close()
	}
	p.producers = map[string]*kafka.Producer{}
}
