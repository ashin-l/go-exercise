package proto

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func ReadMessage(conn net.Conn) (msg Message, err error) {
	var buf [8192]byte
	n, err := conn.Read(buf[0:4])
	if n != 4 {
		err = errors.New("read header failed!")
		return
	}

	var packlen uint32
	packlen = binary.BigEndian.Uint32(buf[0:4])
	n, err = conn.Read(buf[:packlen])
	if n != int(packlen) {
		err = errors.New("read body failed")
		return
	}

	err = json.Unmarshal(buf[0:packlen], &msg)
	if err != nil {
		fmt.Println("unmarshal failed, err:", err)
	}
	return
}

func WriteMessage(cmd string, data interface{}, conn net.Conn) (err error) {
	jdata, err := json.Marshal(data)
	msg := Message{
		Cmd:  cmd,
		Data: string(jdata),
	}
	jmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("marshal failed, ", err)
		return
	}

	var buf [4]byte
	packlen := uint32(len(jmsg))
	binary.BigEndian.PutUint32(buf[:4], packlen)
	n, err := conn.Write(buf[:4])
	if err != nil {
		fmt.Println("write header  failed")
		return
	}

	n, err = conn.Write(jmsg)
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
