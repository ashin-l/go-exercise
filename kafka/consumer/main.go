package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	addrs := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	fmt.Println(addrs)
	//addrs := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("tp33", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Println("Consumed message offset", msg.Offset)
			fmt.Println(string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	fmt.Println("Consumed:", consumed)
}
