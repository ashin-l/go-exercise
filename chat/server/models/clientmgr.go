package models

import "fmt"

type ClientMgr struct {
	onlineUsers map[int]*Client
}

func NewClientMgr() *ClientMgr {
	return &ClientMgr{make(map[int]*Client)}
}

func (m *ClientMgr) AddClient(id int, client *Client) {
	m.onlineUsers[id] = client
}

func (m *ClientMgr) GetClient(id int) (client *Client, err error) {
	client, ok := m.onlineUsers[id]
	if !ok {
		err = fmt.Errorf("user %d not online!", id)
		return
	}
	return
}

func (m *ClientMgr) DelClient(id int) {
	delete(m.onlineUsers, id)
}

func (m *ClientMgr) GetAllUsers() map[int]*Client {
	return m.onlineUsers
}
