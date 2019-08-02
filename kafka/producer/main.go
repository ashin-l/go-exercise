package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	//addrs := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	//addrs := []string{"192.168.152.44:9092", "192.168.152.44:9093", "192.168.152.44:9094"}
	addrs := []string{"192.168.152.21:9092", "192.168.152.22:9092", "192.168.152.23:9092"}
	config := sarama.NewConfig()
	config.Version = sarama.V2_3_0_0
	config.ClientID = "test33"
	admin, err := sarama.NewClusterAdmin(addrs, config)
	if err != nil {
		fmt.Println(err)
	}
	err = admin.CreateTopic("tp33", &sarama.TopicDetail{NumPartitions: 2, ReplicationFactor: 1}, false)
	if err != nil {
		fmt.Println("xxx")
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

	index := 1
	for {
		payload := "test_" + strconv.Itoa(index)
		msg := &sarama.ProducerMessage{Topic: "tp33", Value: sarama.StringEncoder(payload)}
		index++
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("failed to send message: ", err)
		} else {
			fmt.Printf("message sent to partition %d at offset %d\n", partition, offset)
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
