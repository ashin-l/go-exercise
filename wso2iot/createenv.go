package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ashin-l/go-exercise/wso2iot/login"
)

var waitGroup = new(sync.WaitGroup)

func main() {
	fmt.Println("Begin ---")
	now := time.Now()
	host := "https://192.168.152.48:9443"
	client := login.Login(host)
	// step 5 循环创建设备
	//https://192.168.160.19:9443/devicemgt/api/devices/sketch/download?deviceName=env113&latitude=35.1559455&longitude=109.90908&deviceType=EnvMonitor&sketchType=EnvMonitor
	strurl := host + "/devicemgt/api/devices/sketch/download?deviceName=env%s&latitude=%f&longitude=%f&deviceType=EnvMonitor&sketchType=EnvMonitor"
	waitGroup.Add(100)
	for i := 0; i != 10; i++ {
		go create(34.121569, 108.793953, "0", strurl, client)
		go create(33.982307, 109.794306, "1", strurl, client)
		go create(33.45188, 109.966781, "2", strurl, client)
		go create(39.826775, 116.282962, "3", strurl, client)
		go create(31.14635, 121.306038, "4", strurl, client)
		go create(39.053262, 117.075408, "5", strurl, client)
		go create(29.500607, 106.443907, "6", strurl, client)
		go create(31.784154, 117.140586, "7", strurl, client)
		go create(30.508776, 117.008899, "8", strurl, client)
		go create(32.293027, 118.23299, "9", strurl, client)
	}
	waitGroup.Wait()
	fmt.Printf("注册总时间:%v\n", time.Now().Sub(now))
}

func create(latitude, longitude float64, prefix string, strurl string, client *http.Client) {
	for i := 0; i != 100; i++ {
		name := prefix
		b := i / 10
		y := i % 10
		switch {
		case b > 0:
			name = name + strconv.Itoa(i)
		default:
			name = name + "0" + strconv.Itoa(i)
		}
		lat := latitude + 0.04*float64(b)
		lon := longitude + 0.07*float64(y)
		createURL := fmt.Sprintf(strurl, name, lat, lon)
		fmt.Println(createURL)
		resp, err := client.Get(createURL)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		resp.Body.Close()
		status := resp.StatusCode
		fmt.Println(i, status)
	}
	waitGroup.Done()
}
