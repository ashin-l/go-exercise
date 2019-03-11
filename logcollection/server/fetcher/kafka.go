package fetcher

import (
	"fmt"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	"github.com/Shopify/sarama"
)

var consumer sarama.Consumer

func InitKafka() (err error) {
	common.Logger.Info("init kafka...")
	consumer, err = sarama.NewConsumer(common.AppConf.KafkaAddrs, nil)
	if err != nil {
		fmt.Println("init kafka error:", err)
		return
	}
	return
}
