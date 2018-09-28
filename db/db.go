package db

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

var (
	conn *sqlx.DB
)

func InitDB(ds string) {
	conn = sqlx.MustConnect("postgres", ds)
}
