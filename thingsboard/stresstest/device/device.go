package device

import (
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
)

const (
	userstr = `{"username":"%s", "password":"%s"}`
	savestr = `{"name":"test%d", "type":"stress"}`
	djson = `{"deviceid":"%s","value":"%d","clienttime":%d,"other":"%s"}`
)

var (
	cli = &http.Client{}
	other []byte
	successMsgNum    int
	failMsgNum       int
	muSuccessMsg     sync.Mutex
	muFailMsg        sync.Mutex
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

func save(id int, token string) (dv common.Device, err error) {
	str := fmt.Sprintf(savestr, id)
	req, err := http.NewRequest("POST", common.AppConf.Savedevice, bytes.NewReader([]byte(str)))
	if err != nil {
		fmt.Println("new request error:", err)
	}
	req.Header.Set("X-Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Createdevice error, http error code:" + resp.Status)
		return 
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("resp error:", err)
		return
	}
	mdevice := make(map[string]interface{})
	json.Unmarshal(data, &mdevice)
	dv.Name = mdevice["name"].(string)
	dv.DeviceId = mdevice["id"].(map[string]interface{})["id"].(string)
	rurl := fmt.Sprintf(common.AppConf.Getdevicecredentials, dv.DeviceId)
	req.URL, _ = url.Parse(rurl)
	req.Method = "GET"
	resp, err = cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("getDeviCecredentials error, http error code:" + resp.Status)
		return 
	}
	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("resp error:", err)
		return
	}
	json.Unmarshal(data, &mdevice)
	dv.AccessToken = mdevice["credentialsId"].(string)
	err = persist.Insert(&dv)
	return

}

func Create(id, num int) (sdv []common.Device, err error) {
	dv := common.Device{}
	token, err := gettoken()
	if err != nil {
		fmt.Println("get token error:", err)
		return
	}
	for i := 0; i != num; i++ {
		dv, err = save(id, token)
		if err != nil {
			common.Logger.Error(err.Error())
			fmt.Println(i, err)
			i--
			time.Sleep(1 * time.Second)
			continue
		}
		sdv = append(sdv, dv)
		id++
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
		fmt.Println(v.DeviceId)
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
		sdv, err = Create(1, common.AppConf.DeviceNum)
	} else if (l < common.AppConf.DeviceNum) {
		var tmp []common.Device
		tmp, err = Create(sdv[l-1].Id + 1, common.AppConf.DeviceNum - l)
		if err == nil {
			sdv = append(sdv, tmp...)
		}
	}
	if err != nil {
		return
	}
	other = make([]byte, common.AppConf.MsgSize)
	for i := range other {
		other[i] = 'd'
	}
	fmt.Println("start...")
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	for _, v := range sdv {
		go sendData(ctx, v)
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	<-sigChan
	cancel()
	totaltime := time.Now().Sub(start)
	fmt.Println("\tstop...")
	fmt.Println("======================================================")
	fmt.Println("设备数: ", common.AppConf.DeviceNum)
	fmt.Println("成功消息: ", successMsgNum)
	fmt.Println("失败消息: ", failMsgNum)
	fmt.Println("发送消息间隔（毫秒）: ", common.AppConf.Pubinterval)
	fmt.Println("消息体大小（字节）: ", common.AppConf.MsgSize)
	fmt.Printf("发送总时间:%v\n", totaltime)
	fmt.Println("======================================================")

	common.Logger.Info("======================================================")
	common.Logger.Info("设备数: %d", common.AppConf.DeviceNum)
	common.Logger.Info("成功消息: %d", successMsgNum)
	common.Logger.Info("失败消息: %d", failMsgNum)
	common.Logger.Info("发送消息间隔（毫秒）: %d", common.AppConf.Pubinterval)
	common.Logger.Info("消息体大小（字节）: %d", common.AppConf.MsgSize)
	common.Logger.Info("发送总时间:%v\n", totaltime)
	common.Logger.Info("======================================================")
	common.Logger.Close()
	return
}

func sendData(ctx context.Context, dv common.Device) {
	rand.Seed(int64(dv.Id))
	rurl := fmt.Sprintf(common.AppConf.Telemetryup, dv.AccessToken)
	ticker := time.NewTicker(time.Duration(common.AppConf.Pubinterval) * time.Millisecond)
	for {
		select {
		case <- ctx.Done():
			return
		case <- ticker.C:
			data := fmt.Sprintf(djson, dv.DeviceId, rand.Intn(100), time.Now().UnixNano()/1e6, other)
			_, err := http.Post(rurl, "application/json", bytes.NewReader([]byte(data)))
			if err != nil {
				errmsg := fmt.Sprintf("device %s post telemetry failed: %s", dv.Name, err.Error())
				fmt.Println(errmsg)
				common.Logger.Error(errmsg)
				muFailMsg.Lock()
				failMsgNum++
				muFailMsg.Unlock()
				continue
			}
			muSuccessMsg.Lock()
			successMsgNum++
			muSuccessMsg.Unlock()
		}
	}
}