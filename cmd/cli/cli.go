package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	as "github.com/samonzeweb/godb/adapters/sqlite"
)

func main() {
	var init bool
	flag.BoolVar(&init, "i", false, "Init database")
	flag.Parse()

	if init {
		fmt.Println("Setup database")
		db, _ := godb.Open(as.Adapter, "storage.db")
		fmt.Println("Create tables")
		SetupSqlite(db)
		fmt.Println("Fill tables")
		FillSqlite(db)
	}
}

func SetupSqlite(db *godb.DB) {
	createTable := `
	create table users (
		id          integer not null primary key autoincrement,
		name        text not null,
		second_name text null,
		email       text not null,
		password    text not null,
		token       text not null,
		role        text check( role in ('Admin', 'Member') ) not null default 'Member',
		created     date not null
	);
	`
	result, err := db.CurrentDB().Exec(createTable)
	if err != nil {
		fmt.Println("Error while creating DB:", err)
	}
	fmt.Println(result)
}

func FillSqlite(db *godb.DB) {
	su := sqlite.User{
		Name:       "test_name",
		SecondName: "test_second_name",
		Email:      "test@email.com",
		Password:   "pass",
		Token:      "token",
		Role:       sqlite.Member,
		Created:    time.Now(),
	}
	err := db.Insert(&su).Do()
	fmt.Println("After insert")
	if err != nil {
		fmt.Println("Error while inserting", err)
	}
	fmt.Println("User with id:", su.ID, su.Role, "was inserted")
}
