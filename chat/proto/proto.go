package proto

import (
	"github.com/ashin-l/go-exercise/chat/models"
)

const (
	UserLogin    = "UserLogin"
	UserLoginRes = "UserLoginRes"
	UserRegister = "UserRegister"
)

type Message struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type LoginData struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type RegisterCmd struct {
	User models.User `json:"user"`
}

type LoginCmdRes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
