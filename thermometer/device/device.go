package device

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/persist"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	devicenum = 10
	userstr   = `{"username":"%s", "password":"%s"}`
	savestr   = `{"name":"test%d", "type":"stress"}`
	djson     = `{"deviceid":"%s","value":"%d","clienttime":%d,"other":"%s"}`
)

var (
	//clients       map[int]MQTT.Client
	clients       sync.Map
	cli           = &http.Client{}
	other         []byte
	successMsgNum int
	failMsgNum    int
	muSuccessMsg  sync.Mutex
	muFailMsg     sync.Mutex
	successCli    int
	muSuccessCli  sync.Mutex
	wg            sync.WaitGroup
)

func save(id int, expired chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	dv.AccessToken = mdevice["credentialsId"].(string)
	err = persist.Insert(&dv)
	if err != nil {
		common.Logger.Error("insert to db error:", err, ", devicename:", dv.Name)
	}
	return
}

func Create(id, num int) (err error) {
	wg := &sync.WaitGroup{}
	expired := make(chan bool)
	//ticker := time.NewTicker(time.Duration(common.AppConf.Createinterval) * time.Millisecond)
	//defer func() {
	//	ticker.Stop()
	//}()
	//for i := 0; i != num; i++ {
	//	select {
	//	case <-ticker.C:
	//		wg.Add(1)
	//		go save(id, token, expired, wg)
	//		id++
	//	case <-expired:
	//		wg.Wait()
	//		os.Exit(0)
	//	}
	//}

	for i := 0; i != num; i++ {
		wg.Add(1)
		go save(id, expired, wg)
		id++
	}
	return
}

func Run() (err error) {
	common.Logger.Info("in run ...")
	sdv, err := persist.GetDevices(0, common.AppConf.DeviceNum)
	if err != nil {
		return
	}
	l := len(sdv)
	if l == 0 {
		err = Create(1, common.AppConf.DeviceNum)
	} else if l < common.AppConf.DeviceNum {
		err = Create(sdv[l-1].Id+1, common.AppConf.DeviceNum-l)
		var adddv []common.Device
		adddv, err = persist.GetDevices(sdv[l-1].Id, common.AppConf.DeviceNum-l)
		sdv = append(sdv, adddv...)
	}
	if err != nil {
		return
	}
	other = make([]byte, common.AppConf.MsgSize)
	for i := range other {
		other[i] = 'd'
	}
	var send func(context.Context, common.Device)
	switch common.AppConf.Transport {
	case "http":
		send = sendHttp
	case "mqtt":
		//clients = make(map[int]MQTT.Client)
		for _, v := range sdv {
			wg.Add(1)
			go func(v common.Device) {
				opts := MQTT.NewClientOptions().AddBroker(common.AppConf.Mqttbroker)
				opts.SetCleanSession(true)
				opts.SetClientID(v.Name)
				opts.SetUsername(v.AccessToken)
				c := MQTT.NewClient(opts)
				for {
					if token := c.Connect(); token.Wait() && token.Error() != nil {
						fmt.Println(token.Error())
						common.Logger.Error(token.Error().Error(), v.Name)
						fmt.Println("???????????????????????????", v.Name)
						continue
					}
					muSuccessCli.Lock()
					successCli++
					muSuccessCli.Unlock()
					clients.Store(v.Id, c)
					//clients[v.Id] = c
					wg.Done()
					return
				}
			}(v)
			//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++", v.Name, "connected!, num:", successCli)
			//common.Logger.Info("+++++++++++++++++++++++++++++++++++++++++++++++  %s, %s, %d", v.Name, "connected!, num:", successCli)
		}
		send = sendMQTT
	}
	wg.Wait()
	fmt.Println("init...")
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(50 * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	exit := false
Loop:
	for _, v := range sdv {
		select {
		case <-ticker.C:
			go send(ctx, v)
		case <-sigChan:
			exit = true
			break Loop
		}
	}
	fmt.Println("down...")
	common.Logger.Info("============================ down...")
	if !exit {
		<-sigChan
	}
	cancel()
	totaltime := time.Now().Sub(start)
	fmt.Println("\tstop...")
	fmt.Println("======================================================")
	fmt.Println("消息协议: ", common.AppConf.Transport)
	fmt.Println("设备数: ", common.AppConf.DeviceNum)
	fmt.Println("成功连接数: ", successCli)
	fmt.Println("成功消息: ", successMsgNum)
	fmt.Println("失败消息: ", failMsgNum)
	fmt.Println("发送消息间隔（毫秒）: ", common.AppConf.Pubinterval)
	fmt.Println("消息体大小（字节）: ", common.AppConf.MsgSize)
	fmt.Printf("发送总时间:%v\n", totaltime)
	fmt.Println("======================================================")

	common.Logger.Info("======================================================")
	common.Logger.Info("消息协议: %s", common.AppConf.Transport)
	common.Logger.Info("设备数: %d", common.AppConf.DeviceNum)
	common.Logger.Info("成功连接数: %d", successCli)
	common.Logger.Info("成功消息: %d", successMsgNum)
	common.Logger.Info("失败消息: %d", failMsgNum)
	common.Logger.Info("发送消息间隔（毫秒）: %d", common.AppConf.Pubinterval)
	common.Logger.Info("消息体大小（字节）: %d", common.AppConf.MsgSize)
	common.Logger.Info("发送总时间:%v\n", totaltime)
	common.Logger.Info("======================================================")
	common.Logger.Close()

	//for _, c := range clients {
	//	c.Disconnect(5)
	//}
	time.Sleep(10 * time.Second)
	fmt.Println("disconnect!")
	return
}
