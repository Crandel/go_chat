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

func printError(e error) {
	if e != nil {
		fmt.Println("Error while creating DB:", e)
	}
}

func SetupSqlite(db *godb.DB) {
	createUsers := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE users (
		email        VARCHAR(50) UNIQUE PRIMARY KEY,
		name         TEXT NOT NULL,
		second_name  TEXT,
		password     TEXT NOT NULL,
		token        TEXT NOT NULL,
		role         TEXT CHECK( role IN ('Admin', 'Member') ) NOT NULL DEFAULT 'Member',
		created      DATE NOT NULL
	);
	`

	createRooms := `
	DROP TABLE IF EXISTS rooms;
	CREATE TABLE rooms (
		name         TEXT UNIQUE PRIMARY KEY,
		created      DATE NOT NULL
	);
	`

	createMessages := `
	DROP TABLE IF EXISTS messages;
	CREATE TABLE messages (
		id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id      TEXT NOT NULL UNIQUE,
		room_name    TEXT NOT NULL,
		payload      TEXT,
		created      DATE NOT NULL,
		FOREIGN KEY(user_id)   REFERENCES users(email),
		FOREIGN KEY(room_name) REFERENCES rooms(name)
	);
	`
	_, erru := db.CurrentDB().Exec(createUsers)
	printError(erru)
	_, errm := db.CurrentDB().Exec(createMessages)
	printError(errm)
	_, errr := db.CurrentDB().Exec(createRooms)
	printError(errr)
	fmt.Println("DB was successfully setup")
}

func FillSqlite(db *godb.DB) {
	su := sqlite.User{
		Email:      "test@email.com",
		Name:       "test_name",
		SecondName: "test_second_name",
		Password:   "pass",
		Token:      "token",
		Role:       sqlite.Member,
		Created:    time.Now(),
	}
	err := db.Insert(&su).Do()
	printError(err)
	fmt.Println("User with id:", su.Email, " and role ", su.Role, "was inserted")
	sr := sqlite.Room{
		Name:    "room 1",
		Created: time.Now(),
	}
	err = db.Insert(&sr).Do()
	printError(err)
	sm := sqlite.Message{
		RoomName: sr.Name,
		UserID:   su.Email,
		Payload:  "Test message",
		Created:  time.Now(),
	}
	err = db.Insert(sm).Do()
	printError(err)

	fmt.Printf("Message %v was created", sm)
}
