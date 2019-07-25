package main

import (
	"fmt"
	"time"

	//import the Paho Go MQTT library
	"os"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"strconv"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.44:1883")
	opts.SetClientID("mysub")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

  //sub, err := conn.Subscribe("/topic/Consumer/A/VirtualTopic/test", stomp.AckClient)
  //sub, err := conn.Subscribe("topic.Consumer.A.VirtualTopic.test", stomp.AckClient)
  	//topic := "VirtualTopic.test"
  	//topic := "Consumer.A.VirtualTopic.test"
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	time.Sleep(time.Minute)
	c.Disconnect(250)
}
