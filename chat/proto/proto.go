package proto

import (
	"github.com/ashin-l/go-exercise/chat/common"
)

const (
	UserLoginReq     = "UserLoginReq"
	UserLoginRes     = "UserLoginRes"
	UserRegisterReq  = "UserRegisterReq"
	UserRegisterRes  = "UserRegisterRes"
	NotifyUserStatus = "NotifyUserStatus"
	SendMessage      = "SendMessage"
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

type RegisterReqData struct {
	User common.User `json:"user"`
}

type RegisterResData struct {
	Id    int           `json:"id"`
	Error string        `json:"error"`
	Users []common.User `json:"users"`
}

type NotifyUserStatusData struct {
	User   common.User `json:"user"`
	Status int         `json:"status"`
}

type SendMessageData struct {
	UserInfo common.UserInfo `json:"userinfo"`
	Content  string          `json:"content"`
}
