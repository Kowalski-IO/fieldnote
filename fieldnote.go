package main

import (
	"kowalski.io/fieldnote/db"
	"kowalski.io/fieldnote/server"
)

func main() {
	db.InitDB("/Users/brandon/Desktop/filenote.sqlite")
	server.Init()
}
