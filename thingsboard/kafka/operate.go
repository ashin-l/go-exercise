package main

import (
	"time"
	"database/sql"
	"fmt"
)

type DVdata struct {
	Id int `json:"-"`
	Deviceid string `json:"deviceid"`
	Value int `json:"value"`
	Other string `json:"other"`
	Clienttime int64 `json:"clienttime"`
	Servertime int64 `json:"servertime"`
	Createdtime int64 `json:"-"`
	Created time.Time `json:"-"`
}

func Insert(data *DVdata) error {
	fmt.Println(data)
	sqlstr := "insert into dvdata(deviceid, value, other, clienttime, servertime) values($1, $2, $3, $4, $5)"
	_, err := tbdb.Query(sqlstr, data.Deviceid, data.Value, data.Other, data.Clienttime, data.Servertime)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func Delete(id int) {
	sqlstr := "delete from device where id = $1"
	tbdb.Query(sqlstr, id)
}

func Getdatas(index, num int) (data []DVdata, err error) {
	var rows *sql.Rows
	if num > 0 {
		sqlstr := "select * from dvdata where id > $1 limit $2"
		rows, err = tbdb.Query(sqlstr, index, num)
	} else {
		sqlstr := "select * from dvdata"
		rows, err = tbdb.Query(sqlstr)
	}
	for rows.Next() {
		var tmp DVdata
		err = rows.Scan(&tmp.Id, &tmp.Deviceid, &tmp.Value, &tmp.Other, &tmp.Clienttime, &tmp.Servertime, &tmp.Createdtime, &tmp.Created)
		if err != nil {
			return
		}
		if tmp.Createdtime == 0 {
			tmp.Createdtime = tmp.Created.UnixNano()/1e6
		}
		data = append(data, tmp)
	}
	return
}