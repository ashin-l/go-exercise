package fetcher

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	"github.com/Shopify/sarama"
)

func InitKafka() error {
	common.Logger.Info("init kafka...")
	consumer, err := sarama.NewConsumer(common.AppConf.KafkaAddrs, nil)
	if err != nil {
		fmt.Println("init kafka error:", err)
		return err
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
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
	return nil
}
