package config

const (
	BrokerAddress = "localhost:9092"
)

type KafkaMessage struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

