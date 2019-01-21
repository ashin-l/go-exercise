package proto

import (
	"github.com/ashin-l/go-exercise/chat/common"
)

const (
	UserLoginReq = "UserLoginReq"
	UserLoginRes = "UserLoginRes"
	UserRegister = "UserRegister"
)

type Message struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type LoginReqData struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type LoginResData struct {
	Code  int           `json:"code"`
	Error string        `json:"error"`
	Users []common.User `json:"users"`
}

type RegisterCmd struct {
	common.User `json:"user"`
}
