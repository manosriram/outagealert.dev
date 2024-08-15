package models

import (
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type DbConn struct {
	Query *db.Queries
}
