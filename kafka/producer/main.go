package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	addrs := []string{"192.168.152.48:9092", "192.168.152.48:9093", "192.168.152.48:9094"}
	//addrs := []string{"192.168.152.48:9092"}
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin(addrs, config)
	if err != nil {
		fmt.Println(err)
	}
	err = admin.CreateTopic("tp33", &sarama.TopicDetail{NumPartitions: 1, ReplicationFactor: 3}, false)
	if err != nil {
		fmt.Println(err)
	}

	err = admin.Close()
	if err != nil {
		fmt.Println(err)
	}

	producer, err := sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "tp33", Value: sarama.StringEncoder("testing 123")}
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
