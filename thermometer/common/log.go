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
	Logger = logs.GetBeeLogger()
	config := make(map[string]interface{})
	config["filename"] = AppConf.LogPath
	config["level"] = convertLogLevel(AppConf.LogLevel)
	config["maxdays"] = 30
	configstr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("init logger failed, marshal err:", err)
		return err
	}
	Logger.SetLogger("file", string(configstr))
	Logger.EnableFuncCallDepth(true)
	//Logger.Async()
	return nil
}
