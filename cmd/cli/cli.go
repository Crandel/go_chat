package main

import (
	"log"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

func main() {
	SetupSqlite()
}

func SetupSqlite() {
	db, err := godb.Open(sqlite.Adapter, "storage.db")
	createTable := `
	create table users (
		id          integer not null primary key autoincrement,
		name        text not null,
		second_name text null,
		email       text not null,
		token       text not null,
		role        text check( role in ('Admin', 'Member') ) not null default 'Member',
		created     date not null
	);
	`
	_, err = db.CurrentDB().Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
