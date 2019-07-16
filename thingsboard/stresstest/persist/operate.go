package persist

import (
	"database/sql"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
)

func Insert(dv *common.Device) error {
	//sqlstr := "insert into device(name, deviceid, accesstoken) values($1, $2, $3)"
	_, err := st.Exec(dv.Id, dv.Name, dv.DeviceId, dv.AccessToken)
	//_, err := tbdb.Query(sqlstr, dv.Name, dv.DeviceId, dv.AccessToken)
	return err
}

func Delete(id int) {
	delst.Exec(id)
}

func GetDevices(index, num int) (sdv []common.Device, err error) {
	var rows *sql.Rows
	if num > 0 {
		sqlstr := "select * from device where id > $1 order by id limit $2"
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