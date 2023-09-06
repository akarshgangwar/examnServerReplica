package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context, broker,topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received from kafka: ", string(msg.Value))
	}
}

//consumer can be called in this way
	// func dummyCall() {
	// 	ctx := context.Background()
	// 	topicName := "Student1_Exam1"
	// Consume(ctx,config.BrokerAddress,topicName)
	// }
