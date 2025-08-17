package unit_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/kev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafkaProducerPool_Producer_SetConfigMap_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kafkaAdapter := kev.NewKafka()
	pool := kev.NewKafkaProducerPool(kafkaAdapter)
	cfg := map[kev.ConfigKeyKafka]string{
		"abogoboga": "yeamplow", // invalid config map
	}

	producer, err := pool.Producer(cfg)
	assert.Error(t, err)
	assert.Nil(t, producer)
}

func TestKafkaProducerPool_Producer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)

	pool := kev.NewKafkaProducerPool(mockKafka)

	cfg := map[kev.ConfigKeyKafka]string{
		"bootstrap.servers": "localhost:9092",
	}

	mockKafka.EXPECT().NewProducer(gomock.Any()).Return(nil, errors.New("kafka connection failed")).Times(1)

	producer, err := pool.Producer(cfg)
	assert.Error(t, err)
	assert.Nil(t, producer)
	assert.Contains(t, err.Error(), "failed to create kafka producer")
}

func TestKafkaProducer_FatalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)
	mockProducer := mocks.NewMockKevProducer(ctrl)
	eventChan := make(chan kafka.Event, 1)

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	// Set up expectations in order
	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).
		Times(1)

	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		Times(1) // Only called once when producer is created

	// Close will be called when fatal error is handled
	mockProducer.EXPECT().
		Close().
		Times(1)

	// Create pool
	pool := kev.NewKafkaProducerPool(mockKafka)

	// Get producer first to trigger the event handler goroutine
	producer, err := pool.Producer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	// Give a moment for the event handler goroutine to start
	time.Sleep(10 * time.Millisecond)

	// Now send the fatal error
	fatalErr := kafka.NewError(kafka.ErrFatal, "Fatal error", true)
	eventChan <- fatalErr

	// Close the channel to signal end of events
	close(eventChan)

	// Wait for the event handler to process the fatal error and cleanup
	time.Sleep(100 * time.Millisecond)

	// Verify pool stats show the producer was removed
	stats := pool.GetPoolStats()
	assert.Equal(t, 0, stats["active_producers"])
}

// Alternative test that verifies the cleanup more explicitly
func TestKafkaProducer_FatalErrorWithCleanupVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)
	mockProducer := mocks.NewMockKevProducer(ctrl)
	eventChan := make(chan kafka.Event, 1)

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	mockProducer.EXPECT().
		Events().
		Return(eventChan)

	// Use a channel to signal when Close is called
	closeCalled := make(chan bool, 1)
	mockProducer.EXPECT().
		Close().
		Do(func() {
			closeCalled <- true
		})

	pool := kev.NewKafkaProducerPool(mockKafka)

	// Get producer
	producer, err := pool.Producer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	// Verify we have 1 active producer
	stats := pool.GetPoolStats()
	assert.Equal(t, 1, stats["active_producers"])

	// Send fatal error
	fatalErr := kafka.NewError(kafka.ErrFatal, "Fatal error", true)
	eventChan <- fatalErr
	close(eventChan)

	// Wait for Close to be called
	select {
	case <-closeCalled:
		// Close was called as expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Close was not called within timeout")
	}

	// Give a moment for cleanup to complete
	time.Sleep(50 * time.Millisecond)

	// Verify producer was removed from pool
	stats = pool.GetPoolStats()
	assert.Equal(t, 0, stats["active_producers"])
}

func TestKafkaProducerPool_SendMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventChan := make(chan kafka.Event, 1)
	// Mock producer
	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	mockProducer.EXPECT().
		Flush(gomock.Any()).
		Return(0).
		Times(1)

	mockKafka := mocks.NewMockKafka(ctrl)
	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	pool := kev.NewKafkaProducerPool(mockKafka)

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	tpc := "something"
	msg := kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &tpc,
			Partition: kev.KafkaPartitionAny,
		},
		Value: []byte("Something"),
	}

	err := pool.SendMessage(cfg, &msg)
	assert.NoError(t, err)
	close(eventChan)
}

type dummyEvent struct{}

func (dummyEvent) String() string {
	return "Unknown"
}

func TestKafkaProducerPool_HandleEvents_AllCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	eventChan := make(chan kafka.Event, 3)

	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).
		Times(1)

	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	mockProducer.EXPECT().
		KafkaProducer().
		Return(&kafka.Producer{}).
		AnyTimes()

	mockProducer.EXPECT().
		Close().AnyTimes() // akan dipanggil jika ada error fatal

	pool := kev.NewKafkaProducerPool(mockKafka)

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	// Trigger pembuatan producer (dan start handleEvents goroutine)
	producer, err := pool.Producer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	// Simulasi semua jenis event
	// 1. Delivered message
	topic := "test-topic"
	eventChan <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 0,
			Offset:    10,
			Error:     nil,
		},
	}

	// 2. Failed delivery
	eventChan <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    11,
			Error:     kafka.NewError(kafka.ErrUnknown, "fail", false),
		},
	}

	// 4. Stats event
	eventChan <- kafka.Stats{}

	// 5. Log event
	eventChan <- kafka.LogEvent{
		Name:      "logger",
		Tag:       "tag",
		Message:   "log msg",
		Timestamp: time.Now(),
	}

	// 6. Unknown event
	eventChan <- dummyEvent{}

	// Tunggu goroutine selesai (karena fatal error memicu return dan close producer)
	time.Sleep(50 * time.Millisecond)

	close(eventChan)
}

func TestKafkaProducerPool_GetPoolStats_WithEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)

	// Setup 2 mock producer
	mockProducer1 := mocks.NewMockKevProducer(ctrl)
	mockProducer2 := mocks.NewMockKevProducer(ctrl)

	mockProducer1.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockProducer1.EXPECT().Flush(gomock.Any()).Return(0).AnyTimes()
	mockProducer1.EXPECT().Events().Return(make(chan kafka.Event)).AnyTimes()
	mockProducer1.EXPECT().KafkaProducer().Return(&kafka.Producer{}).AnyTimes()

	mockProducer2.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockProducer2.EXPECT().Flush(gomock.Any()).Return(0).AnyTimes()
	mockProducer2.EXPECT().Events().Return(make(chan kafka.Event)).AnyTimes()
	mockProducer2.EXPECT().KafkaProducer().Return(&kafka.Producer{}).AnyTimes()

	// Expect dua kali NewProducer dipanggil dengan config berbeda
	gomock.InOrder(
		mockKafka.EXPECT().NewProducer(gomock.Any()).Return(mockProducer1, nil),
		mockKafka.EXPECT().NewProducer(gomock.Any()).Return(mockProducer2, nil),
	)

	pool := kev.NewKafkaProducerPool(mockKafka)

	// Kirim dua message dengan config berbeda (agar hash berbeda)
	cfg1 := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}
	cfg2 := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9093", // beda port = beda hash
	}

	topic := "test-topic"
	msg := &kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &topic,
			Partition: kev.KafkaPartitionAny,
		},
		Value: []byte("test"),
	}

	// Trigger producer creation
	err1 := pool.SendMessage(cfg1, msg)
	err2 := pool.SendMessage(cfg2, msg)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Sekarang test GetPoolStats
	stats := pool.GetPoolStats()

	assert.Equal(t, 2, stats["active_producers"])
	keys := stats["producer_keys"].([]string)
	assert.Len(t, keys, 2)
	for _, key := range keys {
		assert.True(t, strings.HasSuffix(key, "..."))
		assert.GreaterOrEqual(t, len(key), 11) // 8 + "..."
	}
}

// func TestKafkaProducerPool_CleanupIdleProducers_Indirect(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockProducer := mocks.NewMockKevProducer(ctrl)
// 	mockKafka := mocks.NewMockKafka(ctrl)
//
// 	// Channel events yang tidak perlu isi apa-apa cukup dibuat buffered
// 	eventChan := make(chan kafka.Event, 1)
//
// 	// Setup expectations
// 	mockProducer.EXPECT().KafkaProducer().Return(&kafka.Producer{}).AnyTimes()
// 	mockProducer.EXPECT().Events().Return(eventChan).AnyTimes()
// 	mockProducer.EXPECT().Close().AnyTimes()
//
// 	mockKafka.EXPECT().NewProducer(gomock.Any()).Return(mockProducer, nil).AnyTimes()
//
// 	// Buat pool dengan interval dan maxIdleTime kecil agar cepat cleanup
// 	pool := kev.NewKafkaProducerPool(mockKafka, &kev.CleanUpOpts{
// 		Interval:    10 * time.Millisecond,
// 		MaxIdleTime: 10 * time.Millisecond,
// 	})
//
// 	cfg := map[kev.ConfigKeyKafka]string{
// 		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
// 	}
//
// 	// Panggil Producer untuk bikin pooledProducer
// 	_, err := pool.Producer(cfg)
// 	assert.NoError(t, err)
//
// 	// Set lastUsed ke waktu lama supaya dianggap idle
//
// 	// Tunggu agar cleanupRoutine jalan dan panggil Close()
// 	time.Sleep(50 * time.Millisecond)
// }

func TestKafkaProducerPool_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	// Setup ekspektasi untuk NewProducer, Close, Events, KafkaProducer karena biasanya dipanggil saat Producer()
	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).
		Times(1)

	mockProducer.EXPECT().
		Close().
		Times(1)

	mockProducer.EXPECT().
		Events().
		Return(make(chan kafka.Event)).
		AnyTimes()

	mockProducer.EXPECT().
		KafkaProducer().
		Return(&kafka.Producer{}).
		AnyTimes()

	pool := kev.NewKafkaProducerPool(mockKafka)

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	// Panggil Producer untuk membuat pooled producer supaya ada di map
	_, err := pool.Producer(cfg)
	assert.NoError(t, err)

	// Pastikan ada producer sebelum Close

	// Panggil Close
	pool.Close()

	// Pastikan producer sudah dihapus setelah Close
}

func TestSendMessage_EmptyMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)
	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	// Test with no messages
	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS: "localhost:9092",
	}) // No messages passed

	assert.NoError(t, err)
}

func TestKafkaProducerPool_SendMessage_NewProducerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafka := mocks.NewMockKafka(ctrl)

	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	cfg := map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}

	expectedErr := errors.New("mocked producer creation error")

	// Expect NewProducer to return error
	mockKafka.
		EXPECT().
		NewProducer(gomock.Any()).
		Return(nil, expectedErr)

	msg := &kev.MessageKafka{}

	err := pool.SendMessage(cfg, msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create kafka producer")
	require.ErrorIs(t, err, expectedErr)
}

func TestKafkaProducerPool_SendMessage_ProduceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	eventChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	mockProducer.EXPECT().
		Close().
		Times(1)

	// Produce returns error, SendMessage langsung return error tanpa flush
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("produce error")).
		Times(1)

	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	tpc := "test-topic"
	msg := &kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic: &tpc,
		},
		Value: []byte("test message"),
	}

	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}, msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to produce message")
}

func TestKafkaProducerPool_SendMessage_DeliveryChanClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	eventChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	mockProducer.EXPECT().
		Flush(gomock.Any()).
		Return(0).
		AnyTimes()

	mockProducer.EXPECT().
		Close().AnyTimes()

	// Produce mock: close channel langsung supaya goroutine menerima closed channel
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ *kafka.Message, deliveryChan chan kafka.Event) error {
			close(deliveryChan) // langsung close channel
			return nil
		}).
		Times(1)

	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	tpc := "test-topic"
	msg := &kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic: &tpc,
		},
		Value: []byte("test message"),
	}

	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}, msg)

	require.NoError(t, err)
}

func TestKafkaProducerPool_SendMessage_DeliveryFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	eventChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	mockProducer.EXPECT().
		Close().
		AnyTimes()

	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ *kafka.Message, deliveryChan chan kafka.Event) error {
			// Simulate message delivery with error
			deliveryChan <- &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Error: kafka.NewError(kafka.ErrUnknown, "delivery failed", false),
				},
			}
			close(deliveryChan)
			return nil
		}).
		Times(1)

	mockProducer.EXPECT().
		Flush(gomock.Any()).
		Return(0).
		Times(1)

	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	tpc := "test-topic"
	msg := &kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic: &tpc,
		},
		Value: []byte("test message"),
	}

	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}, msg)

	require.NoError(t, err)

	// Tambahkan sleep singkat untuk beri waktu goroutine proses dan print log
	time.Sleep(20 * time.Millisecond)
}

func TestKafkaProducerPool_SendMessage_FlushRemaining(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKafka := mocks.NewMockKafka(ctrl)

	// Expect NewProducer dipanggil dan sukses
	mockKafka.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	// Expect Events() dipanggil (handleEvents)
	eventChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().
		Events().
		Return(eventChan).
		AnyTimes()

	// Expect Produce dipanggil dan sukses
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockProducer.EXPECT().
		Close().
		AnyTimes()

	// Simulasikan Flush yang return sisa pesan > 0 (misal 2)
	mockProducer.EXPECT().
		Flush(gomock.Any()).
		Return(2).
		Times(1)

	pool := kev.NewKafkaProducerPool(mockKafka)
	defer pool.Close()

	tpc := "test-topic"
	msg := &kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic: &tpc,
		},
		Value: []byte("test message"),
	}

	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
		kev.ConfigKeyKafka("bootstrap.servers"): "localhost:9092",
	}, msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to flush all messages")
	require.Contains(t, err.Error(), "2 messages remaining")
}

// func TestSendMessage_DeliveryChannelSafety(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockKafka := mocks.NewMockKafka(ctrl)
// 	mockProducer := mocks.NewMockKevProducer(ctrl)
// 	eventChan := make(chan kafka.Event, 1)
//
// 	// Setup expectations
// 	mockKafka.EXPECT().
// 		NewProducer(gomock.Any()).
// 		Return(mockProducer, nil).
// 		Times(1)
//
// 	mockProducer.EXPECT().
// 		Events().
// 		Return(eventChan).
// 		AnyTimes()
//
// 	// Add Close expectation
// 	mockProducer.EXPECT().
// 		Close().
// 		AnyTimes()
//
// 	// Simulate multiple delivery events on the same channel
// 	mockProducer.EXPECT().
// 		Produce(gomock.Any(), gomock.Any()).
// 		Do(func(msg *kafka.Message, deliveryChan chan kafka.Event) {
// 			go func() {
// 				// First delivery event
// 				deliveryChan <- &kafka.Message{
// 					TopicPartition: kafka.TopicPartition{
// 						Topic:     msg.TopicPartition.Topic,
// 						Partition: 0,
// 						Error:     nil,
// 					},
// 				}
//
// 				// Second delivery event (shouldn't happen but tests safety)
// 				deliveryChan <- &kafka.Message{
// 					TopicPartition: kafka.TopicPartition{
// 						Topic:     msg.TopicPartition.Topic,
// 						Partition: 0,
// 						Error:     errors.New("unexpected second delivery"),
// 					},
// 				}
// 			}()
// 		}).
// 		Return(nil).
// 		Times(1)
//
// 	mockProducer.EXPECT().
// 		Flush(5000).
// 		Return(0).
// 		Times(1)
//
// 	pool := kev.NewKafkaProducerPool(mockKafka, nil)
// 	defer pool.Close()
//
// 	topic := "test-topic"
// 	msg := &kev.MessageKafka{
// 		TopicPartition: kev.KafkaTopicPartition{
// 			Topic:     &topic,
// 			Partition: 0,
// 		},
// 		Value: []byte("test-value"),
// 	}
//
// 	err := pool.SendMessage(map[kev.ConfigKeyKafka]string{
// 		kev.BOOTSTRAP_SERVERS: "localhost:9092",
// 	}, msg)
//
// 	assert.NoError(t, err)
// 	time.Sleep(100 * time.Millisecond) // Allow delivery goroutine to complete
// }
