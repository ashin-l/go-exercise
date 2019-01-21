package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/ashin-l/go-exercise/chat/common"

	"github.com/ashin-l/go-exercise/chat/proto"
)

type Client struct {
	conn net.Conn
	buf  [8192]byte
}

func (p *Client) readPackage() (msg proto.Message, err error) {
	n, err := p.conn.Read(p.buf[0:4])
	if n != 4 {
		err = errors.New("read header failed!")
		return
	}

	var packlen uint32
	packlen = binary.BigEndian.Uint32(p.buf[0:4])
	fmt.Printf("receive len:%d\n", packlen)
	n, err = p.conn.Read(p.buf[:packlen])
	if n != int(packlen) {
		err = errors.New("read body failed")
		return
	}

	fmt.Printf("receive data:%s\n", string(p.buf[0:packlen]))
	err = json.Unmarshal(p.buf[0:packlen], &msg)
	if err != nil {
		fmt.Println("unmarshal failed, err:", err)
	}
	return
}

func (p *Client) writePackage(data []byte) (err error) {
	packlen := uint32(len(data))

	binary.BigEndian.PutUint32(p.buf[0:4], packlen)
	n, err := p.conn.Write(p.buf[0:4])
	if err != nil {
		fmt.Println("write data  failed")
		return
	}

	n, err = p.conn.Write(data)
	if err != nil {
		fmt.Println("write data  failed")
		return
	}

	if n != int(packlen) {
		fmt.Println("write data  not finished")
		err = errors.New("write data not fninshed")
		return
	}

	return
}

func (p *Client) Process() (err error) {
	for {
		var msg proto.Message
		msg, err = proto.ReadMessage(p.conn)
		if err != nil {
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
	case proto.UserRegister:
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

	var loginRes proto.LoginResData
	loginRes.Code = 200

	if err != nil {
		loginRes.Code = 500
		loginRes.Error = fmt.Sprintf("%v", err)
		return
	}

	data, err := json.Marshal(loginRes)
	if err != nil {
		fmt.Println("marshal failed, ", err)
		return
	}

	respMsg.Data = string(data)
	data, err = json.Marshal(respMsg)
	if err != nil {
		fmt.Println("marshal failed, ", err)
		return
	}
	err = p.writePackage(data)
	if err != nil {
		fmt.Println("send failed, ", err)
		return
	}
	usermgr.Update(user)
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

func (p *Client) register(msg proto.Message) (err error) {
	var cmd proto.RegisterCmd
	err = json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		return
	}

	err = usermgr.Register(&cmd.User)
	if err != nil {
		return
	}

	return
}
