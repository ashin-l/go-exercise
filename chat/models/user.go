package models

import "time"

const (
	UserStatusOnline  = 1
	UserStatusOffline = iota
)

type User struct {
	NickName      string    `json:"nickname"`
	Password      string    `json:"password"`
	Repassword    string    `json:"repassword"`
	Sex           string    `json:"sex"`
	ImgUri        string    `json:"imguri"`
	Lastlogintime time.Time `json:"lastlogintime`
	Createtime    time.Time `json:"createtime`
	Status        int       `json:"status"`
}
