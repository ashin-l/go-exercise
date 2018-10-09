package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ashin-l/go-exercise/wso2iot/login"
)

func main() {
	fmt.Println("Begin ---")
	host := "https://192.168.160.19:9443"
	client := login.Login(host)
	// step 5 循环创建设备
	//https://192.168.160.19:9443/devicemgt/api/devices/sketch/download?deviceName=env113&latitude=35.1559455&longitude=109.90908&deviceType=EnvMonitor&sketchType=EnvMonitor
	strurl := host + "/devicemgt/api/devices/sketch/download?deviceName=env%s&latitude=%f&longitude=%f&deviceType=EnvMonitor&sketchType=EnvMonitor"
	latitude, longitude := 34.111781, 108.785222
	for i := 1; i != 100; i++ {
		s := ""
		b := i / 10
		y := i % 10
		switch {
		case b > 0:
			s = "0" + strconv.Itoa(i)
		default:
			s = "00" + strconv.Itoa(i)
		}
		lat := latitude + 0.01*float64(b)
		lon := longitude + 0.03*float64(y)
		createURL := fmt.Sprintf(strurl, s, lat, lon)
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
}
