package device

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/ashin-l/go-exercise/thermometer/common"
	"github.com/ashin-l/go-exercise/thermometer/persist"
)

var (
	other         []byte
	successMsgNum int
	failMsgNum    int
	muSuccessMsg  sync.Mutex
	muFailMsg     sync.Mutex
)

func Create(id, num int) (err error) {
	for i := 0; i < num; i++ {
		index := strconv.Itoa(id)
		deviceID := common.AppConf.IDprefix + index
		name := common.AppConf.Nameprefix + index
		err = persist.CreateTable(deviceID, common.AppConf.Tenant, name)
		if err != nil {
			return
		}
		id++
	}
	return
}

func saveData(ctx context.Context, dv common.Device) {
	ticker := time.NewTicker(time.Duration(common.AppConf.Pubinterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done ", dv.DeviceID)
			return
		case <-ticker.C:
			err := persist.InsertData(dv.DeviceID, rand.Float32()+30, string(other))
			if err != nil {
				fmt.Println(time.Now(), err)
				muFailMsg.Lock()
				failMsgNum++
				fmt.Printf("device:%s, fail msg num:%d\n", dv.DeviceID, failMsgNum)
				common.Logger.Info("device:%s, fail msg num:%d\n", dv.DeviceID, failMsgNum)
				muFailMsg.Unlock()
			} else {
				muSuccessMsg.Lock()
				successMsgNum++
				//fmt.Printf("device:%s, success msg num:%d\n", dv.DeviceID, successMsgNum)
				muSuccessMsg.Unlock()
				//common.Logger.Info("device:%s, success msg num:%d\n", dv.Name, successMsgNum)
			}
		}
	}
}

func Run() (err error) {
	dvs, err := persist.GetDevices()
	if err != nil {
		return
	}
	l := len(dvs)
	if l < common.AppConf.DeviceNum {
		err = Create(l, common.AppConf.DeviceNum-l)
		fmt.Println("create device finish!")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(50 * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	other = make([]byte, 30)
	for i := range other {
		other[i] = 'd'
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	exit := false
	start := time.Now()
Loop:
	for _, v := range dvs {
		select {
		case <-ticker.C:
			go saveData(ctx, v)
		case <-sigChan:
			exit = true
			break Loop
		}
	}
	totaltime := time.Now().Sub(start)
	fmt.Println("all device start senddata! use ", totaltime)
	//common.Logger.Info("all device start senddata! use :%v\n", totaltime)
	if !exit {
		<-sigChan
	}
	cancel()

	totaltime = time.Now().Sub(start)
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
