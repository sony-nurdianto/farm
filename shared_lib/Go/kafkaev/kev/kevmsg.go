package kev

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type (
	KafkaOffset        = kafka.Offset
	KafkaTimestampType = kafka.TimestampType
)

const (
	OffsetBeginning   KafkaOffset = kafka.OffsetBeginning
	OffsetEnd         KafkaOffset = kafka.OffsetEnd
	OffsetStored      KafkaOffset = kafka.OffsetStored
	OffsetInvalid     KafkaOffset = kafka.OffsetInvalid
	KafkaPartitionAny             = kafka.PartitionAny
)

type KafkaHeader struct {
	Key   string
	Value []byte
}

func (kh KafkaHeader) Factory() kafka.Header {
	out := kafka.Header{
		Key:   kh.Key,
		Value: kh.Value,
	}
	return out
}

type KafkaTopicPartition struct {
	Topic       *string
	Partition   int32
	Offset      KafkaOffset
	Metadata    *string
	Error       error
	LeaderEpoch *int32
}

func (ktp KafkaTopicPartition) Factory() kafka.TopicPartition {
	out := kafka.TopicPartition{
		Topic:       ktp.Topic,
		Partition:   ktp.Partition,
		Offset:      ktp.Offset,
		Metadata:    ktp.Metadata,
		Error:       ktp.Error,
		LeaderEpoch: ktp.LeaderEpoch,
	}

	return out
}

type MessageKafka struct {
	TopicPartition KafkaTopicPartition
	Key            []byte
	Value          []byte
	Timestamp      time.Time
	TimestampType  KafkaTimestampType
	Opaque         any
	Headers        []KafkaHeader
}

func (mk MessageKafka) Factory() kafka.Message {
	header := make([]kafka.Header, 0, len(mk.Headers))
	for _, v := range mk.Headers {
		header = append(header, v.Factory())
	}

	out := kafka.Message{
		TopicPartition: mk.TopicPartition.Factory(),
		Key:            mk.Key,
		Value:          mk.Value,
		Timestamp:      mk.Timestamp,
		TimestampType:  mk.TimestampType,
		Opaque:         mk.Opaque,
		Headers:        header,
	}

	return out
}
