package main

import (
	"fmt"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	//opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.21:31883")
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("mytest")
	opts.SetUsername("A1_TEST_TOKEN")

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "v1/devices/me/telemetry"

	djson := `{
				"clienttime":%d,
				"value":%d
				}`

	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
	}()
	rand.Seed(3)
	for range ticker.C {
		payload := fmt.Sprintf(djson, time.Now().UnixNano()/1e6, rand.Intn(40)+114)
		fmt.Println(payload)
		token := c.Publish(topic, 1, false, payload)
		token.Wait()
	}
}
