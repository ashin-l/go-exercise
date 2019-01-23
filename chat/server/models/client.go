package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/ashin-l/go-exercise/chat/common"
	"github.com/ashin-l/go-exercise/chat/proto"
)

type Client struct {
	user common.User
	conn net.Conn
	buf  [8192]byte
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (p *Client) Process() (err error) {
	for {
		var msg proto.Message
		msg, err = proto.ReadMessage(p.conn)
		if err != nil {
			clientmgr.DelClient(p.user.Id)
			go notifyUserStatus(p.user, common.UserStatusOffline)
			return err
		}

		err = p.processMsg(msg)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
	}
}

func (p *Client) processMsg(msg proto.Message) (err error) {
	switch msg.Cmd {
	case proto.UserLoginReq:
		err = p.login(msg)
	case proto.UserRegisterReq:
		err = p.register(msg)
	default:
		err = errors.New("unsupport message")
		return
	}
	return
}

func (p *Client) loginResp(user *common.User, err error) {
	var respMsg proto.Message
	respMsg.Cmd = proto.UserLoginRes
	var data proto.LoginResData
	if err != nil {
		data.Error = err.Error()
	} else {
		p.user = *user
		clientmgr.AddClient(user.Id, p)
		setOnlineUsers(user.Id, &data.Users)
	}
	go notifyUserStatus(*user, common.UserStatusOnline)
	proto.WriteMessage(proto.UserLoginRes, data, p.conn)
}

func (p *Client) login(msg proto.Message) (err error) {
	var user *common.User
	defer func() {
		p.loginResp(user, err)
	}()

	fmt.Printf("recv user login request, data:%v\n", msg)
	var data proto.LoginReqData
	err = json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return
	}

	user, err = usermgr.Login(data.Id, data.Password)
	return
}

func (p *Client) registerResp(user *common.User, err error) {
	var data proto.RegisterResData
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Id = user.Id
		p.user = *user
		clientmgr.AddClient(user.Id, p)
		setOnlineUsers(user.Id, &data.Users)
	}
	go notifyUserStatus(*user, common.UserStatusOnline)
	proto.WriteMessage(proto.UserRegisterRes, data, p.conn)
}

func (p *Client) register(msg proto.Message) (err error) {
	var data proto.RegisterReqData
	defer func() {
		p.registerResp(&data.User, err)
	}()
	err = json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return
	}

	err = usermgr.Register(&data.User)
	if err != nil {
		return
	}

	return
}

func setOnlineUsers(id int, users *[]common.User) {
	onlineusers := clientmgr.GetAllUsers()
	for k, _ := range onlineusers {
		if k == id {
			continue
		}
		*users = append(*users, onlineusers[k].user)
	}
}

func notifyUserStatus(user common.User, status int) {
	var data proto.NotifyUserStatusData
	data.Status = status
	data.User = user
	onlineusers := clientmgr.GetAllUsers()
	for k, v := range onlineusers {
		if k == user.Id {
			continue
		}
		go proto.WriteMessage(proto.NotifyUserStatus, data, v.conn)
	}
}
