package kev

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

//go:generate mockgen -source=producer_def.go -destination=../test/mocks/kev/mock_producer.go -package=mocks
type KevProducer interface {
	KafkaProducer() *kafka.Producer
	Events() chan kafka.Event
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Flush(timeoutMs int) int
	Close()
}

type producer struct {
	kprod *kafka.Producer
}

func NewProducer(kprod *kafka.Producer) producer {
	return producer{
		kprod: kprod,
	}
}

func (p producer) KafkaProducer() *kafka.Producer {
	return p.kprod
}

func (p producer) Events() chan kafka.Event {
	return p.kprod.Events()
}

func (p producer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	return p.kprod.Produce(msg, deliveryChan)
}

func (p producer) Flush(timeoutMs int) int {
	return p.kprod.Flush(timeoutMs)
}

func (p producer) Close() {
	p.kprod.Close()
}
