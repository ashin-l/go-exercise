package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
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
                        "humiditysensor": %d,
                        "other": "%s"
                    }
                }
            }`
)

var (
	infoLog *log.Logger
	clients map[string]MQTT.Client
	msgnum  int
	m       sync.Mutex
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

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	fmt.Println("Begin ...")
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.48:1886")
	opts.SetClientID("admin:EnvMonitor:MySub1")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions

	clients = make(map[string]MQTT.Client)
	states := make(map[string][2]chan bool)
	for _, deviceId := range dids {
		opts.SetClientID("admin:EnvMonitor:Pub:" + deviceId)
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			panic(token.Error())
		}
		clients[deviceId] = client
		states[deviceId] = [2]chan bool{make(chan bool), make(chan bool)}
		go publishPM(deviceId, states[deviceId][0])
	}
	count := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	for _, state := range states {
		<-ticker.C
		count++
		state[0] <- true
	}
	ticker.Stop()
	fmt.Printf("打开总时间:%v\n", time.Now().Sub(now))
	fmt.Printf("Done! %d devices, time: %s\n", len(states), time.Now())

	cs := make(chan os.Signal)
	signal.Notify(cs)
	// Block until a signal is received.
	s := <-cs
	fmt.Println("Got signal:", s) //Got signal: terminated
	defer fmt.Println("msgnum: ", msgnum)
}

func publish(interval int, deviceId string, sensortype string, state chan bool) {
	topic := "carbon.super/envmonitor/" + deviceId + "/sensorval"
	pmval, hmval := 0, 0
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			if sensortype == "pmsensor" {
				pmval = rand.Intn(40) + 10
			} else {
				hmval = rand.Intn(100)
			}
			payload := fmt.Sprintf(djson, deviceOwner, deviceId, sensortype, time.Now().UnixNano()/1e6, pmval, hmval, "test")
			fmt.Println(payload)
			//infoLog.Printf("PMSensor: DeviceId: %s, time: %d\n", deviceId, mtime)
			token := clients[deviceId].Publish(topic, 1, false, payload)
			//token := client.Publish(topic, 0, true, payload)
			token.Wait()
			m.Lock()
			msgnum++
			m.Unlock()
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
