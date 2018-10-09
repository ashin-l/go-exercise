package main

import (
	"fmt"
	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"os"
    //"strconv"
    "math/rand"
	"time"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.48:1886")
    opts.SetClientID("admin:EnvMonitor")
	opts.SetDefaultPublishHandler(f)
    //opts.SetUsername("admin")
    //opts.SetPassword("")
    opts.SetCleanSession(true)
    //opts.SetProtocolVersion(4)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	//if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	//for i := 0; i < 5; i++ {
	//	text := fmt.Sprintf("this is msg #%d!", i)
	//	token := c.Publish("carbon.super/EnvMonitor/1h7zt68y6xh6c", 0, true, text)
	//	token.Wait()
	//}
    deviceOwner := "admin"
    //deviceType := "connectedcup"
    deviceId := "1iv9ty0z4utko"
    //payload := " {\"event\": {\"metaData\": {\"owner\": \"" + deviceOwner +
    //            "\", \"type\": \"coffeelevel\",\"deviceId\": " +
    //            "\"" + deviceId + "\",\"timestamp\": " + strconv.FormatInt(t.Unix(), 10) +
    //            "},\"payloadData\": { \"coffeelevel\": " + strconv.Itoa(pmvalue) + ", \"temperature\": 0} }}"
    rand.Seed(37)
    djson := " {\"event\": {\"metaData\": {\"owner\": \"" + deviceOwner +
                "\", \"type\": \"coffeelevel\",\"deviceId\": " +
                "\"" + deviceId + "\",\"timestamp\": %d"   +
                "},\"payloadData\": { \"coffeelevel\": %d" + ", \"temperature\": 0} }}"

    ticker := time.NewTicker(5 * time.Second)
    for _ = range ticker.C {
        t := time.Now()
        payload := fmt.Sprintf(djson, t.Unix(), rand.Intn(40) + 10)
        fmt.Println(payload)
        token := c.Publish("carbon.super/connectedcup/1iv9ty0z4utko", 0, true, payload)
        token.Wait()
    }

	//time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	//if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}

	c.Disconnect(250)
}
