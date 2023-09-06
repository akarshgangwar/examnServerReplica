package kafka

import (
	"context"
	"encoding/json"
	"examn_go/config"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)


func Produce(ctx context.Context, broker, topic string, data config.KafkaMessage) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer w.Close()
	// Check if the topic exists
	exists, err := topicExists(broker, topic)
	if err != nil {
		return err
	}
	if !exists {
		// Create the topic if it doesn't exist
		err := createTopic(broker, topic)
		if err != nil {
			return err
		}
		
	}
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = w.WriteMessages(ctx, kafka.Message{
				Key:   []byte(fmt.Sprintf("key-%d", 0)),
				Value: jsonBytes,
			})
			if err != nil {
				return err
			}
			fmt.Printf("JSON data saved to Kafka topic '%s': %s\n", topic, jsonBytes)
			time.Sleep(time.Second)
			return nil
}

func topicExists(broker string, topic string) (bool, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}

func createTopic(broker string, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	return conn.CreateTopics(topicConfig)
}


//producer can be called in this way
	// func dummyCall() {
	// 	ctx := context.Background()
	// 	topicName := "Student1_Exam1"
	
	// 	jsonData := config.KafkaMessage{
	// 		Payload: "2_1.5",
	// 		Type:    "saved_answer",
	// 	}
	
	// 	err := Produce(ctx,config.BrokerAddress, topicName, jsonData)
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 	}
		
	// }
