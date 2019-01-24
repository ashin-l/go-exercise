package models

import (
	"fmt"
	"sync"
)

type ClientMgr struct {
	onlineUsers sync.Map
}

func NewClientMgr() *ClientMgr {
	return new(ClientMgr)
}

func (m *ClientMgr) AddClient(id int, client *Client) {
	m.onlineUsers.Store(id, client)
}

func (m *ClientMgr) GetClient(id int) (client *Client, err error) {
	v, ok := m.onlineUsers.Load(id)
	if !ok {
		err = fmt.Errorf("user %d not online!", id)
		return
	}
	client = v.(*Client)
	return
}

func (m *ClientMgr) DelClient(id int) {
	m.onlineUsers.Delete(id)
}

//func (m *ClientMgr) GetAllUsers() map[int]*Client {
//	return m.onlineUsers
//}
