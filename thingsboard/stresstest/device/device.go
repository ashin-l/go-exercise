package device

import (
	"strconv"
	"math/rand"
	"os/signal"
	"os"
	"time"
	"sync"
	"context"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/persist"
	"net/url"
	"errors"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
	"net/http"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	userstr = `{"username":"%s", "password":"%s"}`
	savestr = `{"name":"test%d", "type":"stress"}`
	djson = `{"deviceid":"%s","value":"%d","clienttime":%d,"other":"%s"}`
)

var (
	clients map[int]MQTT.Client
	cli = &http.Client{}
	other []byte
	successMsgNum    int
	failMsgNum       int
	muSuccessMsg     sync.Mutex
	muFailMsg        sync.Mutex
	successCli int
	//muSuccessCli sync.Mutex
)

func gettoken() (token string, err error) {
	str := fmt.Sprintf(userstr, common.AppConf.Username, common.AppConf.Password)
	resp, err := http.Post(common.AppConf.Gettoken, "application/json", bytes.NewReader([]byte(str)))
	if err != nil {
		fmt.Println("post failed:", err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("resp error:", err)
		return
	}
	mtoken := make(map[string]string)
	json.Unmarshal(data, &mtoken)
	token = mtoken["token"]
	return
}

func save(id int, token string, expired chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	str := fmt.Sprintf(savestr, id)
	req, err := http.NewRequest("POST", common.AppConf.Savedevice, bytes.NewReader([]byte(str)))
	if err != nil {
		fmt.Println("new request error:", err)
		common.Logger.Error("new request error:", err)
	}
	req.Header.Set("X-Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	if resp.StatusCode == 401 {
		fmt.Println("token expired!")
		expired <- true
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Createdevice error, http error code:" + resp.Status + ", devicename:test" + strconv.Itoa(id))
		common.Logger.Error(err.Error())
		return 
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("resp error:", err)
		common.Logger.Error("resp error:", err)
		return
	}
	mdevice := make(map[string]interface{})
	json.Unmarshal(data, &mdevice)
	dv := common.Device{}
	dv.Id = id
	dv.Name = mdevice["name"].(string)
	dv.DeviceId = mdevice["id"].(map[string]interface{})["id"].(string)
	rurl := fmt.Sprintf(common.AppConf.Getdevicecredentials, dv.DeviceId)
	req.URL, _ = url.Parse(rurl)
	req.Method = "GET"
	resp, err = cli.Do(req)
	if err != nil {
		common.Logger.Error("resp error:", err)
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("getDeviCecredentials error, http error code:" + resp.Status + ", devicename:" + dv.Name)
		common.Logger.Error(err.Error())
		return 
	}
	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("resp error:", err)
		common.Logger.Error("resp error:", err, ", devicename:", dv.Name)
		return
	}
	json.Unmarshal(data, &mdevice)
	dv.AccessToken = mdevice["credentialsId"].(string)
	err = persist.Insert(&dv)
	if err != nil {
		common.Logger.Error("insert to db error:", err, ", devicename:", dv.Name)
	}
	return
}

func Create(id, num int) (err error) {
	token, err := gettoken()
	if err != nil {
		fmt.Println("get token error:", err)
		return
	}
	wg := &sync.WaitGroup{}
	expired := make(chan bool)
	ticker := time.NewTicker(time.Duration(common.AppConf.Createinterval) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	for i := 0; i != num; i++ {
		select {
		case <- ticker.C:
			wg.Add(1)
			go save(id, token, expired, wg)
			id++
		case <- expired:
			wg.Wait()
			os.Exit(0)
		}
	}
	return
}

func Delall() {
	sdv, err := persist.GetDevices(0, -1)
	if err != nil {
		fmt.Println(" getdevice error ")
		return
	}
	token, err := gettoken()
	if err != nil {
		fmt.Println("get token error:", err)
		return
	}
	fmt.Println("len:", len(sdv))
	for _, v := range sdv {
		go func(common.Device) {
			fmt.Println(v.Name)
			req, err := http.NewRequest(http.MethodDelete, common.AppConf.Deldevice + v.DeviceId, nil)
			if err != nil {
				fmt.Println("new request error:", err)
				return
			}
			req.Header.Set("X-Authorization", "Bearer " + token)
			req.Header.Set("Content-Type", "application/json")
			resp, err := cli.Do(req)
			if err != nil {
				fmt.Println("del error:", err)
				return
			}
			if resp.StatusCode != 200 {
				fmt.Println("Deletedevice error, http error code:" + resp.Status)
				return
			}
			persist.Delete(v.Id)
		} (v)
		time.Sleep(100 * time.Millisecond)
	}
}

func Run() (err error) {
	common.Logger.Info("in run ...")
	sdv, err := persist.GetDevices(0, common.AppConf.DeviceNum)
	if err != nil {
		return
	}
	l := len(sdv)
	if (l == 0) {
		err = Create(1, common.AppConf.DeviceNum)
	} else if (l < common.AppConf.DeviceNum) {
		err = Create(sdv[l-1].Id + 1, common.AppConf.DeviceNum - l)
		var adddv []common.Device
		adddv, err = persist.GetDevices(sdv[l-1].Id, common.AppConf.DeviceNum - l)
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
		clients = make(map[int]MQTT.Client)
		opts := MQTT.NewClientOptions().AddBroker(common.AppConf.Mqttbroker)
		opts.SetCleanSession(true)
		for _, v := range sdv {
			opts.SetClientID(v.Name)
			opts.SetUsername(v.AccessToken)
			c := MQTT.NewClient(opts)
			if token := c.Connect(); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				common.Logger.Error(token.Error().Error(), v.Name)
				fmt.Println("???????????????????????????")
				continue
			}
			//muSuccessCli.Lock()
			successCli++
			fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++", v.Name, "connected!, num:", successCli)
			//common.Logger.Info("+++++++++++++++++++++++++++++++++++++++++++++++  %s, %s, %d", v.Name, "connected!, num:", successCli)
			//muSuccessCli.Unlock()
			clients[v.Id] = c
		}
		send = sendMQTT
	}
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
		case <- ticker.C:
			go send(ctx, v)
		case <-sigChan:
			exit = true
			break Loop
		}
	}
	fmt.Println("down...")
	common.Logger.Info("============================ down...")
	if !exit {
		<- sigChan
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
	common.Logger.Info("成功连接数: ", successCli)
	common.Logger.Info("成功消息: %d", successMsgNum)
	common.Logger.Info("失败消息: %d", failMsgNum)
	common.Logger.Info("发送消息间隔（毫秒）: %d", common.AppConf.Pubinterval)
	common.Logger.Info("消息体大小（字节）: %d", common.AppConf.MsgSize)
	common.Logger.Info("发送总时间:%v\n", totaltime)
	common.Logger.Info("======================================================")
	common.Logger.Close()

	for _, c := range clients {
		c.Disconnect(5)
	}
	time.Sleep(3 * time.Second)
	fmt.Println("disconnect!")
	return
}

func sendHttp(ctx context.Context, dv common.Device) {
	client := &http.Client{}
	rand.Seed(int64(dv.Id))
	rurl := fmt.Sprintf(common.AppConf.Telemetryup, dv.AccessToken)
	ticker := time.NewTicker(time.Duration(common.AppConf.Pubinterval) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <- ctx.Done():
			return
		case <- ticker.C:
			data := fmt.Sprintf(djson, dv.DeviceId, rand.Intn(100), time.Now().UnixNano()/1e6, other)
			req, err := http.NewRequest("POST", rurl, bytes.NewReader([]byte(data)))
			if err != nil {
				errmsg := fmt.Sprintf("device %s post telemetry failed: %s", dv.Name, err.Error())
				fmt.Println(errmsg)
				common.Logger.Error(errmsg)
				muFailMsg.Lock()
				failMsgNum++
				muFailMsg.Unlock()
				continue
			}
			//req.Close = true
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				common.Logger.Error(err.Error())
				continue
			}
			resp.Body.Close()
			muSuccessMsg.Lock()
			fmt.Printf("dvname:%s, success msg num:%d\n", dv.Name, successMsgNum)
			common.Logger.Info("dvname:%s, success msg num:%d\n", dv.Name, successMsgNum)
			successMsgNum++
			muSuccessMsg.Unlock()
		}
	}
}

func sendMQTT(ctx context.Context, dv common.Device) {
	ticker := time.NewTicker(time.Duration(common.AppConf.Pubinterval) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	rand.Seed(37)
	for {
		select {
		case <- ctx.Done():
			fmt.Println("done ", dv.Name)
			return
		case <- ticker.C:
			payload := fmt.Sprintf(djson, dv.DeviceId, rand.Intn(100), time.Now().UnixNano()/1e6, other)
			token := clients[dv.Id].Publish(common.AppConf.MqttTopic, 1, false, payload)
			if token.Error() != nil {
				fmt.Println(token.Error())
				common.Logger.Error(token.Error().Error(), dv.Name)
				muFailMsg.Lock()
				failMsgNum++
				muFailMsg.Unlock()
			} else if !token.WaitTimeout(time.Duration(common.AppConf.MsgTimeout) * time.Second) {
				muFailMsg.Lock()
				failMsgNum++
				fmt.Println("publish msg time out!", failMsgNum)
				common.Logger.Error("device: %s, publish msg time out! failmsgnum %d", dv.Name, failMsgNum)
				muFailMsg.Unlock()
			} else {
				muSuccessMsg.Lock()
				successMsgNum++
				//fmt.Printf("device:%s, success msg num:%d\n", dv.Name, successMsgNum)
				//common.Logger.Info("device:%s, success msg num:%d\n", dv.Name, successMsgNum)
				muSuccessMsg.Unlock()
			}
		}
	}
}