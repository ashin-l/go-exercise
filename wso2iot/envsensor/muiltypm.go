package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ashin-l/go-exercise/conf"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func cmdHandler(msg MQTT.Message, states map[string][2]chan bool) {
	topics := strings.Split(msg.Topic(), "/")
	if len(topics) != 4 {
		panic("error length!")
	}
	var isOpen bool
	if string(msg.Payload()) == "on" {
		isOpen = true
		fmt.Printf("DeviceId: %s, open!\n", topics[2])
	} else {
		isOpen = false
		fmt.Printf("DeviceId: %s, close!\n", topics[2])
	}
	if v, ok := states[topics[2]]; ok {
		v[0] <- isOpen
		v[1] <- isOpen
	}
}

func opHandler(msg MQTT.Message, states map[string][2]chan bool) {
	var data = make(map[string]interface{})
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		panic(err)
	}
	var isOpen bool
	if v, ok := data["isOpen"]; ok {
		isOpen = v.(bool)
		if isOpen {
			fmt.Println("Open all devices!")
			if v, ok := data["deviceIds"]; ok {
				ids := v.([]interface{})
				ticker := time.NewTicker(500 * time.Millisecond)
				for _, id := range ids {
					<-ticker.C
					if v, ok := states[id.(string)]; ok {
						v[0] <- isOpen
						v[1] <- isOpen
					}
				}
				ticker.Stop()
			}
		} else {
			fmt.Println("Close all devices!")
			if v, ok := data["deviceIds"]; ok {
				ids := v.([]interface{})
				for _, id := range ids {
					if v, ok := states[id.(string)]; ok {
						v[0] <- isOpen
						v[1] <- isOpen
					}
				}
			}
		}
	} else {
		return
	}
	if v, ok := data["deviceIds"]; ok {
		ids := v.([]interface{})
		for _, id := range ids {
			if v, ok := states[id.(string)]; ok {
				v[0] <- isOpen
				v[1] <- isOpen
			}
		}
	}
}

type DbWorker struct {
	//mysql data source name
	Dsn string
}

const (
	deviceOwner = "admin"
	cmdTopic    = "carbon.super/EnvMonitor/+/command"
	opTopic     = "carbon.super/EnvMonitor/operation"
	djson       = `{
                "event": {
                    "metaData": {
                        "owner": "%s",
                        "deviceId": "%s",
                        "type": "%s",
                        "timestamp": %d 
                    },
                    "payloadData": {
                        "pmsensor": %d,
                        "humiditysensor": %d
                    }
                }
            }`
)

var (
	infoLog *log.Logger
	client  MQTT.Client
)

func main() {
	now := time.Now()
	fileName := "info.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.Println("A debug message here")
	infoLog.Printf("haha  %d", 3)
	////配置一个日志格式的前缀
	//infoLog.SetPrefix("[Info]")
	//infoLog.Println("A Info Message here ")
	////配置log的Flag参数
	//infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)
	//infoLog.Println("A different prefix")

	dbw := DbWorker{
		Dsn: "lqc:111@tcp(192.168.152.48:3306)/EnvMonitorDM_DB",
	}
	fmt.Println(dbw.Dsn)
	db, err := sql.Open("mysql", dbw.Dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT  EnvMonitor_DEVICE_ID From EnvMonitor_DEVICE`)
	defer rows.Close()
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
	}
	var dids []string
	did := ""
	for rows.Next() {
		err := rows.Scan(&did)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		if did == "" {
			fmt.Println("DeviceId is nil")
		} else {
			dids = append(dids, did)
		}
	}

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
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	states := make(map[string][2]chan bool)
	for _, deviceId := range dids {
		states[deviceId] = [2]chan bool{make(chan bool), make(chan bool)}
		go publishPM(deviceId, states[deviceId][0])
		go publishHumidity(deviceId, states[deviceId][1])
	}
	count := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	for id, state := range states {
		<-ticker.C
		count++
		fmt.Printf("#%d,DeviceId: %s\n", count, id)
		state[0] <- true
		state[1] <- true
	}
	ticker.Stop()
	fmt.Printf("打开总时间:%v\n", time.Now().Sub(now))
	fmt.Printf("Done! %d devices, time: %s\n", len(states), time.Now())

	token := client.Subscribe(cmdTopic, 2, func(client MQTT.Client, msg MQTT.Message) {
		cmdHandler(msg, states)
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	token = client.Subscribe(opTopic, 2, func(client MQTT.Client, msg MQTT.Message) {
		opHandler(msg, states)
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	var nc chan struct{} //nil channel
	<-nc
	client.Disconnect(250)
}

func publish(interval int, deviceId string, sensortype string, state chan bool) {
	topic := "carbon.super/envmonitor/" + deviceId + "/sensorval"
	mtime := time.Now().UnixNano() / 1e6
	rand.Seed(mtime)
	pmval, hmval := 0, 0
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			mtime += int64(interval * 1000)
			if sensortype == "pmsensor" {
				pmval = rand.Intn(40) + 10
			} else {
				hmval = rand.Intn(100)
			}
			payload := fmt.Sprintf(djson, deviceOwner, deviceId, sensortype, mtime, pmval, hmval)
			//infoLog.Printf("PMSensor: DeviceId: %s, time: %d\n", deviceId, mtime)
			_ = client.Publish(topic, 1, false, payload)
			//token := client.Publish(topic, 1, false, payload)
			//token.Wait()
		case isOpen := <-state:
			if !isOpen {
				ticker.Stop()
				return
			}
		}
	}
}

func publishPM(deviceId string, state chan bool) {
	for {
		select {
		case isOpen := <-state:
			if isOpen {
				publish(20, deviceId, "pmsensor", state)
			}
		}
	}
}

func publishHumidity(deviceId string, state chan bool) {
	for {
		select {
		case isOpen := <-state:
			if isOpen {
				publish(30, deviceId, "humiditysensor", state)
			}
		}
	}
}
