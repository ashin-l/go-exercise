package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ashin-l/go-exercise/wso2iot/login"
)

var (
	waitGroup = new(sync.WaitGroup)
	infoLog   *log.Logger
)

func main() {
	fmt.Println("Begin ---")
	fileName := "info.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.Println("A info message here")
	now := time.Now()
	host := "https://192.168.152.48:9443"
	client := login.Login(host)
	// step 5 循环创建设备
	//https://192.168.160.19:9443/devicemgt/api/devices/sketch/download?deviceName=env113&latitude=35.1559455&longitude=109.90908&deviceType=EnvMonitor&sketchType=EnvMonitor
	strurl := host + "/devicemgt/api/devices/sketch/download?deviceName=envmonitor%s&latitude=%f&longitude=%f&deviceType=EnvMonitor&sketchType=EnvMonitor"
	waitGroup.Add(1000)
	fmt.Println("step 5 -------  开始注册设备")
	for i := 0; i != 1000; i++ {
		go create(i, 28.459033, 86.748047, strurl, client)
		time.Sleep(20 * time.Millisecond)
	}
	waitGroup.Wait()
	fmt.Printf("注册总时间:%v\n", time.Now().Sub(now))
}

func create(i int, latitude, longitude float64, strurl string, client *http.Client) {
	name := strconv.Itoa(i)
	if i%1000 == 0 {
		latitude += 0.04 * float64(i/1000)
	}
	longitude += 0.07 * float64(i)
	createURL := fmt.Sprintf(strurl, name, latitude, longitude)
	infoLog.Println(createURL)
	resp, err := client.Get(createURL)
	if err != nil {
		// handle error
		infoLog.Panicln(err)
		fmt.Println(err)
	}
	resp.Body.Close()
	infoLog.Println(i, resp.StatusCode)
	waitGroup.Done()
}
