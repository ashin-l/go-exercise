package common

import (
	"time"
)

type Device struct {
	Id int
	Name string
	DeviceId string
	AccessToken string
	Created time.Time
}
