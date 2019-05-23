package common

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func InitLogger() error {
	Logger = logs.NewLogger(10000)
	config := make(map[string]interface{})
	config["filename"] = AppConf.LogPath
	config["level"] = convertLogLevel(AppConf.LogLevel)
	configstr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("init logger failed, marshal err:", err)
		return err
	}
	Logger.SetLogger("file", string(configstr))
	return nil
}
