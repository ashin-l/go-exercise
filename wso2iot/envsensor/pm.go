package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ashin-l/go-exercise/conf"

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
	myConfig := new(config.Config)
	myConfig.InitConfig("./deviceConfig.properties")
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	fmt.Println("Begin ...")
	opts := MQTT.NewClientOptions().AddBroker(myConfig.Read("Device-Configurations", "mqtt-ep"))
	opts.SetClientID("admin:EnvMonitor")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	deviceOwner := myConfig.Read("Device-Configurations", "owner")
	deviceType := "pmsensor"
	deviceId := myConfig.Read("Device-Configurations", "deviceId")
	//topic := "carbon.super/envmonitor/" + deviceId
	//topic := "carbon.super/EnvMonitor/" + deviceId + "/command"
	topic := "carbon.super/envmonitor/" + deviceId + "/sensorval"
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	djson := `{
	               "event": {
	                   "metaData": {
	                       "owner": "%s",
	                       "deviceId": "%s",
	                       "type": "%s",
	                       "timestamp": %d
	                   },
	                   "payloadData": {
	                       "pmsensor": %d,
						   "humiditysensor": %d,
						   "other": "%s"
	                   }
	               }
	           }`
	ticker := time.NewTicker(5 * time.Second)
	rand.Seed(37)
	for _ = range ticker.C {
		payload := fmt.Sprintf(djson, deviceOwner, deviceId, deviceType, time.Now().UnixNano()/1e6, rand.Intn(40)+10, 55, "aaaa")
		//fmt.Println(payload)
		token := c.Publish(topic, 1, false, payload)
		token.Wait()
	}

	var cn chan struct{}
	<-cn
	c.Disconnect(250)
}
time.Now().Format("2006-01-02 15:04:05.000")