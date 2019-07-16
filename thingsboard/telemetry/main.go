package main

import (
	"fmt"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.22:31883")
	opts.SetClientID("cli-pub")
	opts.SetUsername("1u5U666PihsxIrECjuHy")

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "v1/devices/me/telemetry"

	djson := `{
				"ts":%d,
				"values":{
					"temperature":%d
					}
				}`

	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		ticker.Stop()
	}()
	rand.Seed(37)
	for _ = range ticker.C {
		payload := fmt.Sprintf(djson, time.Now().UnixNano()/1e6, rand.Intn(40)+114)
		fmt.Println(payload)
		token := c.Publish(topic, 1, false, payload)
		token.Wait()
	}
}
