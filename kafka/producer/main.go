package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder("testing 123")}
	for {
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("failed to send message: ", err)
		} else {
			fmt.Printf("message sent to partition %d at offset %d\n", partition, offset)
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
