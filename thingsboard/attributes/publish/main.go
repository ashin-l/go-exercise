package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.48:1883")
	opts.SetClientID("cli-pub")
	opts.SetUsername("ugnzrLLxYS6Ed4sTiKA9")

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "v1/devices/me/attributes"
	payload := `{"attribute1":"value1", "attribute2":true, "attribute3":42.0, "attribute4":73}`
	token := c.Publish(topic, 1, false, payload)
	token.Wait()
}
