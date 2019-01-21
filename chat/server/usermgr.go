package main

import (
	"github.com/ashin-l/go-exercise/chat/server/models"
)

var usermgr *models.UserMgr

func initUserMgr() {
	usermgr = models.NewUserMgr(pool)
}
