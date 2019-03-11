package fetcher

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func Fetch(topic string) (chan string, error) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := make(chan string)
	go func() {
		for {
			msg := <-partitionConsumer.Messages()
			fmt.Println("Consumed message offset", msg.Offset, msg.Topic)
			out <- string(msg.Value)
		}
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return out, nil
}
