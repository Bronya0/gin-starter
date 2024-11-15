package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
)

func Test_consumer(t *testing.T) {

	// 创建一个消费者，指定GroupID，从 topic-A 消费消息
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "consumer-group-id", // 指定消费者组id
		Topic:    "demo",
		MaxBytes: 10e6, // 10M
	})
	ctx := context.Background()
	// 接收消息
	for {
		m, err := consumer.FetchMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("消息的 topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		// 手动提交
		if err := consumer.CommitMessages(ctx, m); err != nil {
			log.Fatal("提交失败:", err)
		}
	}

	// 程序退出前关闭Reader
	if err := consumer.Close(); err != nil {
		log.Fatal("关闭消费者失败:", err)
	}
}
