package producer

import (
	"context"
	"time"

	"github.com/nickolation/grpc-task-hezzl/configs"
	"github.com/nickolation/grpc-task-hezzl/model"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type KafkaHouseClusterProducer interface {
	PushMessage(ctx context.Context, data []byte) error
	PushClickhouseUserLog(ctx context.Context, user *model.User) error
}

type KafkaHouseProducer struct {
	kafkaWriter *kafka.Writer
	Config      *configs.GrpcGonfigManager
}

const _defaultDialKafkaConnectTimeout = 10 * time.Second

func NewKafkaHouseClusterManager(cfg *configs.GrpcGonfigManager) *KafkaHouseProducer {
	p := new(KafkaHouseProducer)
	p.Config = cfg
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  p.Config.GetKafkaAddrs(),
		Topic:    p.Config.GetKafkaClusterTopic(),
		Balancer: &kafka.LeastBytes{},
		Dialer: &kafka.Dialer{
			Timeout:  _defaultDialKafkaConnectTimeout,
			ClientID: p.Config.GetKafkaCliendId(),
		},
		WriteTimeout:     _defaultDialKafkaConnectTimeout,
		ReadTimeout:      _defaultDialKafkaConnectTimeout,
		CompressionCodec: snappy.NewCompressionCodec(),
	})

	p.kafkaWriter = w
	return p
}

var _defaultKafkaMessageKey []byte

func (khp *KafkaHouseProducer) PushMessage(ctx context.Context, data []byte) error {
	m := kafka.Message{
		Key:   _defaultKafkaMessageKey,
		Value: data,
		Time:  time.Now(),
	}

	return khp.kafkaWriter.WriteMessages(ctx, m)
}
