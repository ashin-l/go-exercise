package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	UsreTable = "user"
)

type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (mgr *UserMgr) {
	mgr = &UserMgr{
		pool: pool,
	}
	return
}

func (m *UserMgr) GetUser(conn redis.Conn, nickname string) (user *User, err error) {
	result, err := redis.String(conn.Do("hget", UsreTable, nickname))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrUserNotExist
		}
		return
	}

	user = &User{}
	err = json.Unmarshal([]byte(result), user)
	if err != nil {
		return
	}
	return
}

func (m *UserMgr) Login(nickname, password string) (user *User, err error) {
	conn := m.pool.Get()
	defer conn.Close()
	user, err = m.GetUser(conn, nickname)
	if err != nil {
		return
	}

	if user.NickName != nickname || user.Password != password {
		err = ErrInvalidPasswd
		return
	}
	user.Status = UserStatusOnline
	user.Lastlogintime = fmt.Sprintf("%v", time.Now())
	return
}

func (m *UserMgr) Register(user *User) (err error) {
	conn := m.pool.Get()
	defer conn.Close()
	if user == nil {
		fmt.Println("invalid user!")
		err = ErrInvalidParams
		return
	}

	_, err = m.GetUser(conn, user.NickName)
	if err == nil {
		err = ErrUserExist
		return
	}

	if err != ErrUserNotExist {
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	_, err = conn.Do("hset", UsreTable, user.NickName, string(data))
	return
}

func (m *UserMgr) Update(user *User) error {
	data, _ := json.Marshal(user)
	_, err := m.pool.Get().Do("hset", UsreTable, user.NickName, data)
	return err
}
