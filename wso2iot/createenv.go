package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ashin-l/go-exercise/wso2iot/login"
)

func main() {
	fmt.Println("Begin ---")
	host := "https://192.168.152.48:9443"
	client := login.Login(host)
	// step 5 循环创建设备
	//https://192.168.160.19:9443/devicemgt/api/devices/sketch/download?deviceName=env113&latitude=35.1559455&longitude=109.90908&deviceType=EnvMonitor&sketchType=EnvMonitor
	strurl := host + "/devicemgt/api/devices/sketch/download?deviceName=env%s&latitude=%f&longitude=%f&deviceType=EnvMonitor&sketchType=EnvMonitor"
	//create(34.111781, 108.785222, "0", strurl)
	//create(34.219237, 108.797547, "1", strurl, client)
	create(33.376421, 108.697655, "2", strurl, client)
	/*
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
			lat := latitude + 0.05*float64(b)
			lon := longitude + 0.08*float64(y)
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
	*/
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
}
