package kev

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KevConsumer interface {
	Assign(partitions []kafka.TopicPartition) (err error)
	Assignment() ([]kafka.TopicPartition, error)
	AssignmentLost() bool
	Close() error
	Commit() ([]kafka.TopicPartition, error)
	CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error)
	CommitOffsets(offsets []kafka.TopicPartition) ([]kafka.TopicPartition, error)
	Committed(partitions []kafka.TopicPartition, timeoutMs int) ([]kafka.TopicPartition, error)
	GetConsumerGroupMetadata() (*kafka.ConsumerGroupMetadata, error)
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
	GetRebalanceProtocol() string
	GetWatermarkOffsets(topic string, partition int32) (int64, int64, error)
	IncrementalAssign(partitions []kafka.TopicPartition) error
	IncrementalUnassign(partitions []kafka.TopicPartition) error
	IsClosed() bool
	Logs() chan kafka.LogEvent
	OffsetsForTimes(times []kafka.TopicPartition, timeoutMs int) ([]kafka.TopicPartition, error)
	Pause(partitions []kafka.TopicPartition) (err error)
	Poll(timeoutMs int) (event kafka.Event)
	Position(partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error)
	QueryWatermarkOffsets(topic string, partition int32, timeoutMs int) (int64, int64, error)
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Resume(partitions []kafka.TopicPartition) (err error)
	Seek(partition kafka.TopicPartition, ignoredTimeoutMs int) error
	SeekPartitions(partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error)
	SetOAuthBearerToken(oauthBearerToken kafka.OAuthBearerToken) error
	SetOAuthBearerTokenFailure(errstr string) error
	SetSaslCredentials(username string, password string) error
	StoreMessage(m *kafka.Message) ([]kafka.TopicPartition, error)
	StoreOffsets(offsets []kafka.TopicPartition) ([]kafka.TopicPartition, error)
	String() string
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	Subscription() (topics []string, err error)
	Unassign() (err error)
	Unsubscribe() (err error)
}

type consumer struct {
	kcons *kafka.Consumer
}

func NewConsumer(kcons *kafka.Consumer) consumer {
	return consumer{
		kcons: kcons,
	}
}

func (c consumer) KafkaConsumer() *kafka.Consumer {
	return c.kcons
}

func (c consumer) Assign(partitions []kafka.TopicPartition) (err error) {
	return c.kcons.Assign(partitions)
}

func (c consumer) Assignment() ([]kafka.TopicPartition, error) {
	return c.kcons.Assignment()
}

func (c consumer) AssignmentLost() bool {
	return c.kcons.AssignmentLost()
}

func (c consumer) Close() error {
	return c.Close()
}

func (c consumer) Commit() ([]kafka.TopicPartition, error) {
	return c.kcons.Commit()
}

func (c consumer) CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error) {
	return c.kcons.CommitMessage(m)
}

func (c consumer) CommitOffsets(offsets []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	return c.kcons.CommitOffsets(offsets)
}

func (c consumer) Committed(partitions []kafka.TopicPartition, timeoutMs int) ([]kafka.TopicPartition, error) {
	return c.kcons.Committed(partitions, timeoutMs)
}

func (c consumer) GetConsumerGroupMetadata() (*kafka.ConsumerGroupMetadata, error) {
	return c.kcons.GetConsumerGroupMetadata()
}

func (c consumer) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	return c.kcons.GetMetadata(topic, allTopics, timeoutMs)
}

func (c consumer) GetRebalanceProtocol() string {
	return c.kcons.GetRebalanceProtocol()
}

func (c consumer) GetWatermarkOffsets(topic string, partition int32) (int64, int64, error) {
	return c.kcons.GetWatermarkOffsets(topic, partition)
}

func (c consumer) IncrementalAssign(partitions []kafka.TopicPartition) error {
	return c.kcons.IncrementalAssign(partitions)
}

func (c consumer) IncrementalUnassign(partitions []kafka.TopicPartition) error {
	return c.kcons.IncrementalUnassign(partitions)
}

func (c consumer) IsClosed() bool {
	return c.kcons.IsClosed()
}

func (c consumer) Logs() chan kafka.LogEvent {
	return c.kcons.Logs()
}

func (c consumer) OffsetsForTimes(times []kafka.TopicPartition, timeoutMs int) ([]kafka.TopicPartition, error) {
	return c.kcons.OffsetsForTimes(times, timeoutMs)
}

func (c consumer) Pause(partitions []kafka.TopicPartition) (err error) {
	return c.kcons.Pause(partitions)
}

func (c consumer) Poll(timeoutMs int) (event kafka.Event) {
	return c.kcons.Poll(timeoutMs)
}

func (c consumer) Position(partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	return c.kcons.Position(partitions)
}

func (c consumer) QueryWatermarkOffsets(topic string, partition int32, timeoutMs int) (int64, int64, error) {
	return c.kcons.QueryWatermarkOffsets(topic, partition, timeoutMs)
}

func (c consumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	return c.kcons.ReadMessage(timeout)
}

func (c consumer) Resume(partitions []kafka.TopicPartition) (err error) {
	return c.kcons.Resume(partitions)
}

func (c consumer) Seek(partition kafka.TopicPartition, ignoredTimeoutMs int) error {
	return c.kcons.Seek(partition, ignoredTimeoutMs)
}

func (c consumer) SeekPartitions(partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	return c.kcons.SeekPartitions(partitions)
}

func (c consumer) SetOAuthBearerToken(oauthBearerToken kafka.OAuthBearerToken) error {
	return c.kcons.SetOAuthBearerToken(oauthBearerToken)
}

func (c consumer) SetOAuthBearerTokenFailure(errstr string) error {
	return c.kcons.SetOAuthBearerTokenFailure(errstr)
}

func (c consumer) SetSaslCredentials(username string, password string) error {
	return c.kcons.SetSaslCredentials(username, password)
}

func (c consumer) StoreMessage(m *kafka.Message) ([]kafka.TopicPartition, error) {
	return c.kcons.StoreMessage(m)
}

func (c consumer) StoreOffsets(offsets []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	return c.kcons.StoreOffsets(offsets)
}

func (c consumer) String() string {
	return c.kcons.String()
}

func (c consumer) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error {
	return c.kcons.Subscribe(topic, rebalanceCb)
}

func (c consumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	return c.kcons.SubscribeTopics(topics, rebalanceCb)
}

func (c consumer) Subscription() (topics []string, err error) {
	return c.kcons.Subscription()
}

func (c consumer) Unassign() (err error) {
	return c.kcons.Unassign()
}

func (c consumer) Unsubscribe() (err error) {
	return c.kcons.Unsubscribe()
}
