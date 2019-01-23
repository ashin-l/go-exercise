package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ashin-l/go-exercise/chat/common"
	"github.com/gomodule/redigo/redis"
)

const (
	UsreTable = "user:"
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

func (m *UserMgr) GetUser(conn redis.Conn, id int) (user *common.User, err error) {
	key := UsreTable + strconv.Itoa(id)
	result, err := redis.StringMap(conn.Do("hgetall", key))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrUserNotExist
		}
		return
	}

	user = &common.User{
		Id:            id,
		NickName:      result["nickname"],
		Password:      result["password"],
		Sex:           result["sex"],
		ImgUri:        result["imguri"],
		Lastlogintime: result["lastlogintime"],
		Createtime:    result["createtime"],
	}
	return
}

func (m *UserMgr) Login(id int, password string) (user *common.User, err error) {
	conn := m.pool.Get()
	defer conn.Close()
	user, err = m.GetUser(conn, id)
	if err != nil {
		return
	}

	if user.Password != password {
		err = ErrInvalidPasswd
		return
	}
	user.Lastlogintime = fmt.Sprintf("%v", time.Now())
	return
}

func (m *UserMgr) Register(user *common.User) (err error) {
	conn := m.pool.Get()
	defer conn.Close()
	if user == nil {
		fmt.Println("invalid user!")
		err = ErrInvalidParams
		return
	}

	result, err := redis.Int(conn.Do("sismember", "nickname", user.NickName))
	if err != nil {
		return
	}
	if result == 1 {
		err = ErrNicknameExist
		return
	}

	conn.Do("sadd", "nickname", user.NickName)
	id, err := redis.Int(conn.Do("incr", "userid"))
	key := UsreTable + strconv.Itoa(id)
	user.Id = id
	_, err = conn.Do("hmset", key,
		"id", user.Id,
		"nickname", user.NickName,
		"password", user.Password,
		"sex", user.Sex,
		"imguri", user.ImgUri,
		"createtime", fmt.Sprintf("%v", time.Now()),
		"status", common.UserStatusOffline)
	if err != nil {
		conn.Do("srem", "nickname", user.NickName)
		conn.Do("decr", "userid")
	}
	return
}

func (m *UserMgr) Update(user *common.User) error {
	data, _ := json.Marshal(user)
	_, err := m.pool.Get().Do("hset", UsreTable, user.NickName, data)
	return err
}
