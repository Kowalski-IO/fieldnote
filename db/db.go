package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	conn *sqlx.DB
)

func InitDB(ds string) {
	conn = sqlx.MustConnect("sqlite3", ds)
}
