package main

import (
	"fmt"
	"log"
	"os"
	"os/signal" //import the Paho Go MQTT library
	"strconv"
	"time"

	"github.com/ashin-l/go-exercise/conf"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"time"
)

var (
	msgnum    int
	totalnum  int
	totaltime time.Duration
	start     time.Time
	infoLog   *log.Logger
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	if msgnum == 0 {
		start = time.Now()
	}
	msgnum++
	fmt.Printf("TOPIC: %s, MSGNUM: %d\n", msg.Topic(), msgnum)
	if msgnum >= totalnum {
		fmt.Println("Finish!")
		totaltime = time.Now().Sub(start)
		fmt.Println(start)
		fmt.Println("接受消息总时间： ", totaltime)
		infoLog.Println("接受消息总时间： ", totaltime)
	}
}

var fonLost MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	opt := client.OptionsReader()
	popt := &opt
	fmt.Println("connection lost!", err.Error(), popt.ClientID())
	infoLog.Println("connection lost!", err.Error(), popt.ClientID())
}

func main() {
	conf := new(config.Config)
	fileName := "logs/sub.log"
	logFile, err := os.Create(fileName)
	conf.InitConfig("client.conf")
	server := conf.Read("base", "server")
	if server == "" {
		fmt.Println("server 不能为空！")
		os.Exit(1)
	}
	topic := conf.Read("base", "subtopic")
	if topic == "" {
		fmt.Println("subtopic 不能为空！")
		os.Exit(1)
	}
	subid := conf.Read("base", "subid")
	if subid == "" {
		subid = "gosub"
	}
	maxClient, _ := strconv.Atoi(conf.Read("base", "maxClient"))
	pubTimes, _ := strconv.Atoi(conf.Read("base", "pubTimes"))
	totalnum = maxClient * pubTimes
	defer logFile.Close()
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.Println("A info message here")
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(server)
	opts.SetClientID(subid)
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)
	opts.SetConnectionLostHandler(fonLost)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("connected ok!")

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	//if token := c.Subscribe("carbon.super/test/aaa", 0, nil); token.Wait() && token.Error() != nil {
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	cs := make(chan os.Signal)
	signal.Notify(cs)
	s := <-cs
	fmt.Println("Got signal:", s) //Got signal: terminated
	if msgnum != 0 && totaltime == 0 {
		totaltime = time.Now().Sub(start)
		fmt.Println("接受消息总时间： ", totaltime)
		infoLog.Println("接受消息总时间： ", totaltime)
	}
	defer c.Disconnect(100)
}
