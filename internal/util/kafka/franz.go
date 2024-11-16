package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"os"
	"strings"
	"sync"
)

type FranzContext struct {
	bootstrapServers []string
	config           []kgo.Opt
	kac              *kadm.Client
	client           *kgo.Client
	mutex            sync.Mutex
}

func New(connect *Connect) (*FranzContext, error) {
	var config []kgo.Opt
	// TLS配置
	if connect.enableTls {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: connect.skipTLSVerify, // 开发环境可以设置为true
		}
		// 如果需要证书认证
		if connect.tlsCertFile != "" && connect.tlsKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(connect.tlsCertFile, connect.tlsKeyFile)
			if err != nil {
				return nil, err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		// 如果需要CA证书
		if connect.tlsCaFile != "" {
			caCert, err := os.ReadFile(connect.tlsCaFile)
			if err != nil {
				return nil, err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}
		config = append(config, kgo.DialTLSConfig(tlsConfig))
	}

	// SASL配置
	if connect.enableSasl {
		// SASL机制设置
		switch strings.ToUpper(connect.saslMechanism) {
		case "PLAIN":
			config = append(config, kgo.SASL(plain.Auth{User: connect.saslUser, Pass: connect.saslPwd}.AsMechanism()))
		case "SCRAM-SHA-256":
			config = append(config, kgo.SASL(scram.Auth{User: connect.saslUser, Pass: connect.saslPwd}.AsSha256Mechanism()))
		case "SCRAM-SHA-512":
			config = append(config, kgo.SASL(scram.Auth{User: connect.saslUser, Pass: connect.saslPwd}.AsSha512Mechanism()))
		default:
			return nil, errors.New(fmt.Sprintf("不支持的SASL机制: %s", connect.saslMechanism))
		}
	}

	config = append(config, kgo.SeedBrokers(connect.bootstrapServers...))

	cl, err := kgo.NewClient(config...)
	if err != nil {
		return nil, err
	}
	return &FranzContext{
		bootstrapServers: connect.bootstrapServers,
		config:           config,
		kac:              kadm.NewClient(cl),
		client:           cl,
	}, nil
}

func (c *FranzContext) Produce(topic string, key, value string, partition, num int, headers []map[string]string) error {
	if c.kac == nil {
		return errors.New("请先初始化连接")
	}
	ctx := context.Background()
	headers2 := make([]kgo.RecordHeader, len(headers))
	for i := 0; i < len(headers); i++ {
		headers2[i] = kgo.RecordHeader{
			Key:   headers[i]["key"],
			Value: []byte(headers[i]["value"]),
		}
	}
	for i := 0; i < num; i++ {
		c.client.Produce(ctx, &kgo.Record{
			Topic:     topic,
			Value:     []byte(value),
			Key:       []byte(key),
			Headers:   headers2,
			Partition: int32(partition),
		}, nil)
	}
	return nil
}
