package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
)

func Test_producer(t *testing.T) {
	// 创建一个writer 向topic-A发送消息
	producer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		RequiredAcks: kafka.RequireOne, // ack模式
		Async:        false,            // 异步or同步
		BatchSize:    512000,           // 批次体积
		BatchBytes:   5000000,          //请求最大体积
		Compression:  kafka.Gzip,       // lz4压缩
	}

	var msgs []kafka.Message
	for i := 0; i < 1; i++ {
		msgs = append(msgs, kafka.Message{
			Topic: "demo",
			Value: []byte("Hello World!"),
		})
	}
	err := producer.WriteMessages(context.Background(),
		msgs...,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := producer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
