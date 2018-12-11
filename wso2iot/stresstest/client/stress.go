package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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

var fonLost MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	opt := client.OptionsReader()
	popt := &opt
	fmt.Println("connection lost!", err.Error(), popt.ClientID())
	infoLog.Println("connection lost!", err.Error(), popt.ClientID())
	muLost.Lock()
	lostNum++
	muLost.Unlock()
}

type DbWorker struct {
	//mysql data source name
	Dsn string
}

const (
	perCount    = 2000
	createTime  = 2
	maxClient   = 3000
	pubTimes    = 60
	pubInterval = 3
	msgTimeOut  = 30
	msgQos      = 0
	djson       = `{
                "event": {
                    "metaData": {
                        "owner": "admin",
                        "deviceId": "%s",
                    },
                    "payloadData": {
						"client_time": %d,
						"pmsensor": 37,
						"other" : "%s",
                    }
                }
            }`
)

var (
	infoLog          *log.Logger
	clients          map[string]MQTT.Client
	successClientNum int
	failClientNum    int
	successMsgNum    int
	failMsgNum       int
	lostNum          int
	pubClientNum     int
	muSuccessClient  sync.Mutex
	muFailClient     sync.Mutex
	muSuccessMsg     sync.Mutex
	muFailMsg        sync.Mutex
	muLost           sync.Mutex
	muPub            sync.Mutex
	wg               sync.WaitGroup
	totaltime        time.Duration
)

func main() {
	start := time.Now()
	fileName := "logs/stress.log"
	logFile, err := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.SetFlags(log.Lshortfile)
	infoLog.Println("A info message here")
	////配置一个日志格式的前缀
	//infoLog.SetPrefix("[Info]")
	//infoLog.Println("A Info Message here ")
	////配置log的Flag参数
	//infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)
	//infoLog.Println("A different prefix")

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	fmt.Println("Begin ...")

	//clients = make(map[string]MQTT.Client)
	//states := make(map[string]chan bool)
	ticker := time.NewTicker(1000 * time.Millisecond)
	stime := 1000000 / perCount * time.Microsecond
	i := 0
	for i != createTime {
		i++
		go createClient(strconv.Itoa(i), stime)
		<-ticker.C
	}
	ticker.Stop()
	/*
		ticker = time.NewTicker(100 * time.Millisecond)
		for _, state := range states {
			state <- true
			<-ticker.C
		}
		ticker.Stop()
	*/
	go func(start time.Time) {
		wg.Wait()
		fmt.Println("All finish!")
		totaltime = time.Now().Sub(start)
	}(start)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	sig := <-sigChan
	fmt.Println("Got signal:", sig)
	infoLog.Println("======================================================")
	infoLog.Println("successClientNum: ", successClientNum)
	infoLog.Println("failClientNum: ", failClientNum)
	infoLog.Println("pubClientNum: ", pubClientNum)
	infoLog.Println("successMsgNum: ", successMsgNum)
	infoLog.Println("failMsgNum: ", failMsgNum)
	infoLog.Println("每秒创建连接数: ", perCount)
	infoLog.Println("创建间隔: ", stime)
	infoLog.Println("创建时间（秒）: ", createTime)
	infoLog.Println("发送消息超时（秒）: ", msgTimeOut)
	infoLog.Println("发送消息间隔（秒）: ", pubInterval)
	infoLog.Println("每个连接发送消息次数: ", pubTimes)
	infoLog.Println("发送消息Qos: ", msgQos)
	infoLog.Println("断开连接数: ", lostNum)
	infoLog.Printf("执行总时间:%v\n", totaltime)
	infoLog.Println("======================================================")
	logFile.Close()

	/*
		cs := make(chan os.Signal)
		signal.Notify(cs)
		s := <-cs
		fmt.Println("Got signal:", s) //Got signal: terminated
		defer fmt.Println("msgnum: ", msgnum)
	*/
}

func createClient(prefix string, stime time.Duration) {
	wg.Add(1)
	defer wg.Done()
	//create and start a client using the above ClientOptions
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.152.48:1886")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)
	opts.SetConnectionLostHandler(fonLost)
	prefix += "_stressDeviceId_"
	for i := 0; i != perCount; i++ {
		muSuccessClient.Lock()
		if successClientNum >= maxClient {
			muSuccessClient.Unlock()
			return
		}
		deviceId := prefix + strconv.Itoa(i)
		opts.SetClientID("Pub:" + deviceId)
		client := MQTT.NewClient(opts)
		token := client.Connect()
		if token.Error() != nil {
			fmt.Println(token.Error())
			infoLog.Println(token.Error(), deviceId)
			muFailClient.Lock()
			failClientNum++
			muFailClient.Unlock()
			continue
		}
		if !token.WaitTimeout(30 * time.Second) {
			fmt.Println("create client time out!")
			infoLog.Println("create client time out!", deviceId)
			muFailClient.Lock()
			failClientNum++
			muFailClient.Unlock()
			continue
		}
		wg.Add(1)
		//go publish(client, deviceId)
		go publishInterval(client, deviceId)
		successClientNum++
		muSuccessClient.Unlock()
		time.Sleep(stime)
	}
}

func publish(client MQTT.Client, deviceId string) {
	defer wg.Done()
	//defer client.Disconnect(0)
	topic := "carbon.super/stresstest/" + deviceId + "/test"
	payload := fmt.Sprintf(djson, deviceId, time.Now().UnixNano()/1e6, "test")
	token := client.Publish(topic, msgQos, false, payload)
	if token.Error() != nil {
		fmt.Println(token.Error())
		infoLog.Println(token.Error(), deviceId)
		muFailMsg.Lock()
		failMsgNum++
		muFailMsg.Unlock()
		return
	}
	if !token.WaitTimeout(msgTimeOut * time.Second) {
		fmt.Println("publish msg time out!")
		infoLog.Println("publish msg time out!")
		muFailMsg.Lock()
		failMsgNum++
		muFailMsg.Unlock()
		return
	}
	muSuccessMsg.Lock()
	successMsgNum++
	muSuccessMsg.Unlock()
}

func publishInterval(client MQTT.Client, deviceId string) {
	defer wg.Done()
	defer client.Disconnect(0)
	muPub.Lock()
	pubClientNum++
	muPub.Unlock()
	topic := "carbon.super/stresstest/" + deviceId + "/test"
	ticker := time.NewTicker(time.Duration(pubInterval) * time.Second)
	i := 0
	for {
		if i >= pubTimes {
			break
		}
		i++
		<-ticker.C
		payload := fmt.Sprintf(djson, deviceId, time.Now().UnixNano()/1e6, "test")
		token := client.Publish(topic, msgQos, false, payload)
		if token.Error() != nil {
			fmt.Println(token.Error())
			infoLog.Println(token.Error(), deviceId)
			muFailMsg.Lock()
			failMsgNum++
			muFailMsg.Unlock()
			continue
		}
		if !token.WaitTimeout(msgTimeOut * time.Second) {
			fmt.Println("publish msg time out!")
			infoLog.Println("publish msg time out!", deviceId)
			muFailMsg.Lock()
			failMsgNum++
			muFailMsg.Unlock()
			continue
		}
		muSuccessMsg.Lock()
		successMsgNum++
		muSuccessMsg.Unlock()
	}
	ticker.Stop()
}

func publishPM(deviceId string, state chan bool) {
	for {
		select {
		case isOpen := <-state:
			if isOpen {
				//publishInterval(5, deviceId, state)
			}
		}
	}
}
