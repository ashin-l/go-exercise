package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-syncdb" {
		Syncdb()
		os.Exit(0)
	}
	err := InitDB()
	if err != nil {
		fmt.Println("initdb error:", err)
		return
	}

	/*
	datas, err := Getdatas(0, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(datas)
	*/


	addrs := []string{"192.168.152.44:9092"}
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("topic_stress", 0, sarama.OffsetNewest)
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

	fmt.Println("start...")
	consumed := 0
	data := &DVdata{}
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			json.Unmarshal(msg.Value, data)
			Insert(data)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	fmt.Println("Consumed:", consumed)
}
