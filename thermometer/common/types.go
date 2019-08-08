package common

import "database/sql"

type DBConn struct {
	DBS        *sql.DB
	StmtInsert *sql.Stmt
}

type Device struct {
	DeviceID string
	TenantID string
	Name     string
}
