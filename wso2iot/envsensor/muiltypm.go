package main

import (
	"database/sql"
	"fmt"

	"math/rand"
	"time"

	"github.com/ashin-l/go-exercise/conf"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

type DbWorker struct {
	//mysql data source name
	Dsn string
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	dbw := DbWorker{
		Dsn: "lqc:111@tcp(192.168.152.48:3306)/EnvMonitorDM_DB",
	}
	fmt.Println(dbw.Dsn)
	db, err := sql.Open("mysql", dbw.Dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT  EnvMonitor_DEVICE_ID From EnvMonitor_DEVICE`)
	defer rows.Close()
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
	}
	var dids []string
	did := ""
	for rows.Next() {
		err := rows.Scan(&did)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		if did == "" {
			fmt.Println("DeviceId is nil")
		} else {
			dids = append(dids, did)
		}
	}

	//	for _, id := range dids {
	//		fmt.Println(id)
	//	}
	//	fmt.Println(len(dids))
	//	panic("Finish")

	myConfig := new(config.Config)
	myConfig.InitConfig("./deviceConfig.properties")
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	fmt.Println("Begin ...")
	opts := MQTT.NewClientOptions().AddBroker(myConfig.Read("Device-Configurations", "mqtt-ep"))
	opts.SetClientID("admin:EnvMonitor")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	deviceOwner := myConfig.Read("Device-Configurations", "owner")
	deviceType := "pmsensor"
	deviceId := myConfig.Read("Device-Configurations", "deviceId")
	topic := "carbon.super/envmonitor/" + deviceId

	//	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
	//		fmt.Println(token.Error())
	//		os.Exit(1)
	//	}

	djson := `{
                "event": {
                    "metaData": {
                        "owner": "%s",
                        "deviceId": "%s",
                        "type": "%s",
                        "timestamp": %d 
                    },
                    "payloadData": {
                        "pmsensor": %d,
                        "humiditysensor": %d
                    }
                }
            }`
	ticker := time.NewTicker(5 * time.Second)
	rand.Seed(37)
	mtime := time.Now().UnixNano() / 1e6
	for _ = range ticker.C {
		mtime += 5000
		payload := fmt.Sprintf(djson, deviceOwner, deviceId, deviceType, mtime, rand.Intn(40)+10, 55)
		//fmt.Println(payload)
		token := c.Publish(topic, 0, true, payload)
		token.Wait()
	}

	c.Disconnect(250)
}
