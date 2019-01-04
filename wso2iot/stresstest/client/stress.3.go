package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
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

var (
	server   string
	djson    string
	other    []byte
	prefix   string
	pubtopic string
	infoLog  *log.Logger
	//clients          map[string]MQTT.Client
	successClientNum int
	perCount         int
	createInterval   int
	maxClient        int
	pubTimes         int
	pubInterval      int
	msgTimeOut       int
	msgQos           int
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
	finishClientTime time.Duration
	totaltime        time.Duration
)

func main() {
	conf := new(config.Config)
	conf.InitConfig("client.conf")
	server = conf.Read("base", "server")
	if server == "" {
		fmt.Println("server 不能为空！")
		os.Exit(1)
	}
	prefix = conf.Read("base", "prefix")
	if prefix == "" {
		prefix = "Pub"
	}
	djson = conf.Read("base", "payload")
	pubtopic = conf.Read("base", "pubtopic")
	if pubtopic == "" {
		fmt.Println("pubtopic 不能为空！")
		os.Exit(1)
	}
	msgSize, _ := strconv.Atoi(conf.Read("base", "msgSize"))
	other = make([]byte, msgSize)
	for i := range other {
		other[i] = 'd'
	}
	payload := fmt.Sprintf(djson, prefix+"_stressdeviceid_1000", time.Now().UnixNano()/1e6, other)
	payloadSize := len(payload)
	fmt.Println("消息体大小（字节）：", payloadSize)
	perCount, _ = strconv.Atoi(conf.Read("base", "perCount"))
	createInterval, _ = strconv.Atoi(conf.Read("base", "createInterval"))
	maxClient, _ = strconv.Atoi(conf.Read("base", "maxClient"))
	pubTimes, _ = strconv.Atoi(conf.Read("base", "pubTimes"))
	pubInterval, _ = strconv.Atoi(conf.Read("base", "pubInterval"))
	msgTimeOut, _ = strconv.Atoi(conf.Read("base", "msgTimeOut"))
	msgQos, _ = strconv.Atoi(conf.Read("base", "msgQos"))
	fileName := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
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

	//clients = make(map[string]MQTT.Client, maxClient)
	//states := make(map[string]chan bool)
	start := time.Now()
	ticker := time.NewTicker(time.Duration(createInterval) * time.Millisecond)
	for i := 0; i < maxClient; i++ {
		<-ticker.C
		wg.Add(1)
		deviceId := prefix + "_stressDeviceId_" + strconv.Itoa(i)
		fmt.Println(deviceId)
		go createClient(deviceId)
	}
	ticker.Stop()
	//wg.Wait()
	finishClientTime = time.Now().Sub(start)
	//fmt.Printf("设备创建完成，用时：%v\n", finishClientTime)
	//start = time.Now()
	//for k := range clients {
	//	wg.Add(1)
	//	go publishInterval(k)
	//}
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
	fmt.Println("======================================================")
	fmt.Println("成功客户端: ", successClientNum)
	fmt.Println("失败客户端: ", failClientNum)
	fmt.Println("已发消息客户端: ", pubClientNum)
	fmt.Println("成功消息: ", successMsgNum)
	fmt.Println("失败消息: ", failMsgNum)
	fmt.Printf("创建设备用时:%v\n", finishClientTime)
	fmt.Println("发送消息超时（秒）: ", msgTimeOut)
	fmt.Println("发送消息间隔（毫秒）: ", pubInterval)
	fmt.Println("每个客户端发送消息次数: ", pubTimes)
	fmt.Println("发送消息Qos: ", msgQos)
	fmt.Println("消息体大小（字节）: ", payloadSize)
	fmt.Println("断开连接数: ", lostNum)
	fmt.Printf("发送总时间:%v\n", totaltime)
	fmt.Println("======================================================")

	infoLog.Println("======================================================")
	infoLog.Println("成功客户端: ", successClientNum)
	infoLog.Println("失败客户端: ", failClientNum)
	infoLog.Println("已发消息客户端: ", pubClientNum)
	infoLog.Println("成功消息: ", successMsgNum)
	infoLog.Println("失败消息: ", failMsgNum)
	infoLog.Println("每秒创建连接数: ", perCount)
	infoLog.Printf("创建设备用时:%v\n", finishClientTime)
	infoLog.Println("发送消息超时（秒）: ", msgTimeOut)
	infoLog.Println("发送消息间隔（毫秒）: ", pubInterval)
	infoLog.Println("每个连接发送消息次数: ", pubTimes)
	infoLog.Println("发送消息Qos: ", msgQos)
	infoLog.Println("消息体大小（字节）: ", payloadSize)
	infoLog.Println("断开连接数: ", lostNum)
	infoLog.Printf("发送总时间:%v\n", totaltime)
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

func createClient(deviceId string) {
	//create and start a client using the above ClientOptions
	defer wg.Done()
	opts := MQTT.NewClientOptions().AddBroker(server)
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)
	opts.SetConnectionLostHandler(fonLost)
	opts.SetKeepAlive(300)
	opts.SetClientID(deviceId)
	client := MQTT.NewClient(opts)
	token := client.Connect()
	if token.Error() != nil {
		fmt.Println(token.Error())
		infoLog.Println(token.Error(), deviceId)
		muFailClient.Lock()
		failClientNum++
		muFailClient.Unlock()
		return
	}
	if !token.WaitTimeout(30 * time.Second) {
		fmt.Println("create client time out!")
		infoLog.Println("create client time out!", deviceId)
		muFailClient.Lock()
		failClientNum++
		muFailClient.Unlock()
		return
	}
	muSuccessClient.Lock()
	successClientNum++
	//clients[deviceId] = client
	muSuccessClient.Unlock()
	wg.Add(1)
	go publishInterval(deviceId, client)
}

func publish(client MQTT.Client, deviceId string) {
	defer wg.Done()
	//defer client.Disconnect(0)
	topic := fmt.Sprintf(pubtopic, deviceId)
	payload := fmt.Sprintf(djson, deviceId, time.Now().UnixNano()/1e6, "test")
	token := client.Publish(topic, byte(msgQos), false, payload)
	if token.Error() != nil {
		fmt.Println(token.Error())
		infoLog.Println(token.Error(), deviceId)
		muFailMsg.Lock()
		failMsgNum++
		muFailMsg.Unlock()
		return
	}
	if !token.WaitTimeout(time.Duration(msgTimeOut) * time.Second) {
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

func publishInterval(deviceId string, client MQTT.Client) {
	defer wg.Done()
	//defer client.Disconnect(0)
	muPub.Lock()
	pubClientNum++
	muPub.Unlock()
	topic := fmt.Sprintf(pubtopic, deviceId)
	ticker := time.NewTicker(time.Duration(pubInterval) * time.Millisecond)
	i := 0
	for {
		if i >= pubTimes {
			break
		}
		i++
		payload := fmt.Sprintf(djson, deviceId, time.Now().UnixNano()/1e6, other)
		token := client.Publish(topic, byte(msgQos), false, payload)
		if token.Error() != nil {
			fmt.Println(token.Error())
			infoLog.Println(token.Error(), deviceId)
			if client.IsConnected() {
				fmt.Println("client connected 1")
				infoLog.Println("client connected 1")
			}
			if client.IsConnectionOpen() {
				fmt.Println("client connected 2")
				infoLog.Println("client connected 2")
			}
			muFailMsg.Lock()
			failMsgNum++
			muFailMsg.Unlock()
		} else if !token.WaitTimeout(time.Duration(msgTimeOut) * time.Second) {
			fmt.Println("publish msg time out!")
			infoLog.Println("publish msg time out!", deviceId)
			muFailMsg.Lock()
			failMsgNum++
			muFailMsg.Unlock()
		} else {
			muSuccessMsg.Lock()
			successMsgNum++
			muSuccessMsg.Unlock()
		}
		<-ticker.C
	}
	ticker.Stop()
}
