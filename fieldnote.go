package main

import (
	"kowalski.io/fieldnote/db"
	"kowalski.io/fieldnote/server"
)

const cockroach = "dbname=fieldnote user=root password=password host=localhost port=26257 sslmode=disable"

func main() {
	db.InitDB(cockroach)
	server.Init()
}
