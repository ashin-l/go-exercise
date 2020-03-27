package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/astaxie/beego/config"
)

var sData SensorData
var AppConf *AppConfig

type AppConfig struct {
	Addr       string
	URL        string
	KafkaAddrs []string
	ESAddr     string
	EtcdAddrs  []string
}

func main() {
	fmt.Println("server start...")
	//resp, _ := http.Post("http://192.168.152.48/api/v1/ugnzrLLxYS6Ed4sTiKA9/telemetry", "application/json", strings.NewReader("{'haha':'hehe'}"))
	//fmt.Println(resp.Status)
	err := InitConfig("ini", "agent1.conf")
	if err != nil {
		panic(err)
		return
	}
	sData = SensorData{}
	l, err := net.Listen("tcp", AppConf.Addr)
	if err != nil {
		fmt.Println("server error: ", err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("connect error: ", err.Error())
			continue
		}
		fmt.Println("connect!!!")
		go process(conn)
	}
}

func process(conn net.Conn) {
	fmt.Println("1111111111111")
	defer conn.Close()
	for {
		fmt.Println("xxx")
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err.Error())
			return
		}
		sData.Temperature = ByteToFloat32(buf, 7)
		sData.Humidity = ByteToFloat32(buf, 12)
		sData.Airpressure = ByteToFloat32(buf, 17)
		sData.Windspeed = ByteToFloat32(buf, 22)
		sData.Winddirection = ByteToFloat32(buf, 27)
		sData.Timestamp = time.Now().UnixNano() / 1e6
		jData, err := json.Marshal(sData)
		if err != nil {
			fmt.Println("marshal json failed:", err)
		} else {
			fmt.Println("jData:", jData)
			resp, err := http.Post(AppConf.URL, "application/json", bytes.NewReader(jData))
			if err != nil {
				fmt.Println("post failed:", err)
			} else {
				fmt.Println(resp.Status)
			}
		}
	}

}

func ByteToFloat32(buf []byte, index int) float32 {
	bits := binary.LittleEndian.Uint32(buf[index : index+4])
	return math.Float32frombits(bits)
}

func InitConfig(confType, filename string) error {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		return err
	}

	AppConf = &AppConfig{}
	AppConf.Addr = conf.String("addr")
	if len(AppConf.Addr) == 0 {
		AppConf.Addr = "0.0.0.0:7108"
	}

	AppConf.URL = conf.String("url")
	if len(AppConf.URL) == 0 {
		return fmt.Errorf("ERROR: must set url in conf file!")
	}
	return nil
}
