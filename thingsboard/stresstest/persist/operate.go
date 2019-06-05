package persist

import (
	"database/sql"
	"os"
	"fmt"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
)

func checkDB() {
	if tbdb == nil {
		fmt.Println("error: not connect db!")
		os.Exit(0)
	}
}

func Insert(dv *common.Device) error {
	checkDB()
	sqlstr := "insert into device(name, deviceid, accesstoken) values($1, $2, $3)"
	_, err := tbdb.Query(sqlstr, dv.Name, dv.DeviceId, dv.AccessToken)
	return err
}

func Delete(id int) {
	sqlstr := "delete from device where id = $1"
	tbdb.Query(sqlstr, id)
}

func GetDevices(index, num int) (sdv []common.Device, err error) {
	checkDB()
	var rows *sql.Rows
	if num > 0 {
		sqlstr := "select * from device where id > $1 limit $2"
		rows, err = tbdb.Query(sqlstr, index, num)
	} else {
		sqlstr := "select * from device"
		rows, err = tbdb.Query(sqlstr)
	}
	for rows.Next() {
		var tmp common.Device
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.DeviceId, &tmp.AccessToken, &tmp.Created)
		if err != nil {
			return
		}
		sdv = append(sdv, tmp)
	}
	return
}