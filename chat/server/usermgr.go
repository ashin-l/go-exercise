package main

import (
	"github.com/ashin-l/go-exercise/chat/models"
)

var mgr *models.UserMgr

func initUserMgr() {
	mgr = models.NewUserMgr(pool)
}
