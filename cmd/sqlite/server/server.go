package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/http/rest"
	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	as "github.com/samonzeweb/godb/adapters/sqlite"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	db, _ := godb.Open(as.Adapter, "storage.db")
	sqlite_service := sqlite.NewStorage(db)
	ls := login.NewService(&sqlite_service)
	sis := signin.NewService(&sqlite_service)
	as := adding.NewService(&sqlite_service)
	rs := reading.NewService(&sqlite_service)
	router := rest.InitHandlers(ls, sis, as, rs)
	log.Fatal(http.ListenAndServe(":8080", router))
}
