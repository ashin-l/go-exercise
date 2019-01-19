package models

const (
	UserStatusOnline  = 1
	UserStatusOffline = iota
)

type User struct {
	NickName      string `json:"nickname"`
	Password      string `json:"password"`
	Repassword    string `json:"repassword"`
	Sex           string `json:"sex"`
	ImgUri        string `json:"imguri"`
	Lastlogintime string `json:"lastlogintime`
	Createtime    string `json:"createtime`
	Status        int    `json:"status"`
}
