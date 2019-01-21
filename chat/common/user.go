package common

const (
	UserStatusOnline  = 1
	UserStatusOffline = iota
)

type User struct {
	Id            int    `json:id`
	NickName      string `json:"nickname"`
	Password      string `json:"password"`
	Sex           string `json:"sex"`
	ImgUri        string `json:"imguri"`
	Lastlogintime string `json:"lastlogintime`
	Createtime    string `json:"createtime`
	Status        int    `json:"status"`
}
