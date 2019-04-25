package main

import (
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.48:1883")
	opts.SetClientID("cli-pub")
	opts.SetUsername("ugnzrLLxYS6Ed4sTiKA9")
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subtopic := "v1/devices/me/attributes"
	if token := c.Subscribe(subtopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	defer func() {
		if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		fmt.Println("in defer")
		c.Disconnect(250)
	}()

	time.Sleep(10 * time.Second)
}
