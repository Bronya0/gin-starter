package kafka

import "github.com/segmentio/kafka-go"

type Kafka struct {
	Address string
}

func NewKafkaProducer(Address string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(Address),
		RequiredAcks: kafka.RequireOne, // ack模式
		BatchSize:    512000,           // 批次体积
		BatchBytes:   5000000,          //请求最大体积
		Compression:  kafka.Gzip,       // 压缩
	}
}
