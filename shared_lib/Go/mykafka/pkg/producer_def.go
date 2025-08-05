package pkg

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Producer interface {
	KKProducer() *kafka.Producer
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

func (p producer) KKProducer() *kafka.Producer {
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
