package unit_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sony-nurdianto/farm/shared_lib/Go/mykafka/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
	events chan kafka.Event
}

func (m *MockProducer) KKProducer() *kafka.Producer {
	args := m.Called()
	return args.Get(0).(*kafka.Producer)
}

func (m *MockProducer) Events() chan kafka.Event {
	args := m.Called()
	return args.Get(0).(chan kafka.Event)
}

func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	args := m.Called(msg, deliveryChan)
	return args.Error(0)
}

func (m *MockProducer) Flush(timeoutMs int) int {
	args := m.Called(timeoutMs)
	return args.Int(0)
}

func (m *MockProducer) Close() {
	m.Called()
}

// MockKafka implements the Kafka interface for testing
type MockKafka struct {
	mock.Mock
}

func (m *MockKafka) NewProducer(conf *kafka.ConfigMap) (pkg.Producer, error) {
	args := m.Called(conf)
	return args.Get(0).(pkg.Producer), args.Error(1)
}

type dummyEvent struct{}

func (dummyEvent) String() string {
	return "unknow"
}

func TestGetOrCreateProducer(t *testing.T) {
	mockKafka := new(MockKafka)
	mockProd := new(MockProducer)
	mockProd.events = make(chan kafka.Event, 4)

	// Setup ekspektasi
	cfg := kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

	mockKafka.On("NewProducer", &cfg).Return(mockProd, nil)
	mockProd.On("KKProducer").Return(&kafka.Producer{})
	mockProd.On("Events").Return(mockProd.events).Maybe()
	mockProd.On("Close").Run(func(args mock.Arguments) {
		close(mockProd.events)
	}).Return()

	pool := pkg.NewKafkaProducerPool(mockKafka)

	mockProd.events <- kafka.Stats{}

	mockProd.events <- dummyEvent{}
	// Kirim kafka.LogEvent
	mockProd.events <- kafka.LogEvent{
		Timestamp: time.Now(),
		Name:      "test-log",
		Message:   "this is a log message",
	}

	// trigger agar event handler jalan dan selesai
	mockProd.events <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &[]string{"test-topic"}[0],
			Partition: 0,
			Offset:    1,
		},
	}

	// Panggil producer
	p1, err := pool.Producer(map[pkg.ConfigKeyKafka]string{
		pkg.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	})
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Cleanup
	pool.Close()

	// Assertion akhir (opsional)
	mockKafka.AssertExpectations(t)
	mockProd.AssertExpectations(t)
}

func TestGetOrCreateProducer_PartitionError(t *testing.T) {
	mockKafka := new(MockKafka)
	mockProd := new(MockProducer)
	mockProd.events = make(chan kafka.Event, 4)

	// Setup ekspektasi
	cfg := kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

	mockKafka.On("NewProducer", &cfg).Return(mockProd, nil)
	mockProd.On("KKProducer").Return(&kafka.Producer{})
	mockProd.On("Events").Return(mockProd.events).Maybe()
	mockProd.On("Close").Run(func(args mock.Arguments) {
		close(mockProd.events)
	}).Return()

	pool := pkg.NewKafkaProducerPool(mockKafka)

	mockProd.events <- kafka.Stats{}

	// Kirim kafka.LogEvent
	mockProd.events <- kafka.LogEvent{
		Timestamp: time.Now(),
		Name:      "test-log",
		Message:   "this is a log message",
	}

	errKafka := kafka.NewError(kafka.ErrUnknownPartition, "UnknowPartition", false)

	// trigger agar event handler jalan dan selesai
	mockProd.events <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &[]string{"test-topic"}[0],
			Partition: 0,
			Offset:    1,
			Error:     errKafka,
		},
	}

	// Panggil producer
	p1, err := pool.Producer(map[pkg.ConfigKeyKafka]string{
		pkg.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	})
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Cleanup
	pool.Close()

	// Assertion akhir (opsional)
	mockKafka.AssertExpectations(t)
	mockProd.AssertExpectations(t)
}

func TestGetOrCreateProducer_FatalError(t *testing.T) {
	mockKafka := new(MockKafka)
	mockProd := new(MockProducer)
	mockProd.events = make(chan kafka.Event, 3)

	// Setup ekspektasi
	cfg := kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

	mockKafka.On("NewProducer", &cfg).Return(mockProd, nil)
	mockProd.On("KKProducer").Return(&kafka.Producer{})
	mockProd.On("Events").Return(mockProd.events).Maybe()
	mockProd.On("Close").Run(func(args mock.Arguments) {
		fmt.Println("✅ MOCK Close called (fatal)")
		close(mockProd.events)
	}).Return()

	pool := pkg.NewKafkaProducerPool(mockKafka)

	errKafka := kafka.NewError(kafka.ErrFatal, "Fatality", true)
	mockProd.events <- errKafka

	p1, err := pool.Producer(map[pkg.ConfigKeyKafka]string{
		pkg.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	})
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// beri waktu goroutine handleEvents berjalan
	time.Sleep(50 * time.Millisecond)

	pool.Close()

	mockKafka.AssertExpectations(t)
	mockProd.AssertExpectations(t)

	mockProd.AssertCalled(t, "Close") // pastikan Close dipanggil
}

func TestKafkaProducerPool_Producer_ConfigMapError(t *testing.T) {
	mockKafka := new(MockKafka)
	mockProducer := new(MockProducer)

	cfg := map[pkg.ConfigKeyKafka]string{
		"": "localhost:9092",
	}

	expectedConfigMap := kafka.ConfigMap{"": "localhost:9092"}
	mockKafka.On("NewProducer", &expectedConfigMap).Return(mockProducer, errors.New("invalid config key"))

	pool := pkg.NewKafkaProducerPool(mockKafka)

	_, err := pool.Producer(cfg)
	if err == nil {
		t.Error("Expected Producer to return error for invalid config")
	}

	mockKafka.AssertExpectations(t)
}
