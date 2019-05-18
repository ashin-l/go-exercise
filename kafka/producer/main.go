package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	//addrs := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	addrs := []string{"192.168.152.12:9092"}
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

	msg := &sarama.ProducerMessage{Topic: "demo_kafka_topic_cxf", Value: sarama.StringEncoder("testing 123")}
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
