package main

import (
	"fmt"
	"time"

	//import the Paho Go MQTT library
	"os"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"strconv"
	"github.com/astaxie/beego/config"
)

type AppConfig struct {
	Addr   string
	Topic string
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		return
	}

	AppConf := AppConfig{}
	AppConf.Addr = conf.String("addr")
	if len(AppConf.Addr) == 0 {
		fmt.Println("配置文件出错，addr 不能为空")
		os.Exit(0)
	}
	AppConf.Topic = conf.String("topic")
	if len(AppConf.Topic) == 0 {
		fmt.Println("配置文件出错，topic 不能为空")
	}


	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(AppConf.Addr)
	opts.SetClientID("mytest")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	if token := c.Subscribe(AppConf.Topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	djson := `{
		"payloadData": {
			"timeStamp": %d,
			"value": %d
		}
	}
	`

	value := 1
	ticker := time.NewTicker(5 * time.Second)
	mtime := time.Now().UnixNano() / 1e6
	for range ticker.C {
		mtime += 5000
		payload := fmt.Sprintf(djson, mtime, value)
		token := c.Publish(AppConf.Topic, 0, false, payload)
		token.Wait()
		value++
	}

	c.Disconnect(250)
}
