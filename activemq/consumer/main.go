package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jjeffery/stomp"
)

func main() {
	conn, err := stomp.Dial("tcp", "192.168.152.44:61613")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	destination := "Consumer.A.VirtualTopic.mytopic"
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 3; i++ {
		id := "A" + strconv.Itoa(i)
		fmt.Println("clientID:", id)
		//destination := fmt.Sprintf(topictpl, id)
		go subscribe(ctx, conn, destination, id)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	sig := <-sigChan
	fmt.Println("Got signal:", sig)
	cancel()
	conn.Disconnect()
	fmt.Println("Desconnect")
}

func subscribe(ctx context.Context, conn *stomp.Conn, destination, id string) {
	sub, err := conn.Subscribe(destination, stomp.AckClient)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer sub.Unsubscribe()
	defer fmt.Printf("=======================   close client: %s\n", id)
	chclose := make(chan bool)
	if id == "A2" {
		go func() {
			<-time.After(30 * time.Second)
			chclose <- true
		}()
	}

	for {
		select {
		case <-ctx.Done():
			err := sub.Unsubscribe()
			if err != nil {
				fmt.Println(id, err)
			}
			return
		case m := <-sub.C:
			if m.Err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("clientID:%s ,msg body:%s\n\n", id, m.Body)
			//m.Conn.Ack(m)
		case <-chclose:
			return
		}
	}
}
