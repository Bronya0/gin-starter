package kafka

type Connect struct {
	bootstrapServers []string
	enableTls        bool
	enableSasl       bool
	skipTLSVerify    bool
	tlsCertFile      string
	tlsKeyFile       string
	tlsCaFile        string
	saslUser         string
	saslPwd          string
	saslMechanism    string
}

type Kafka interface {
	Produce(topic string, key, value string, partition, num int, headers []map[string]string) error
}
