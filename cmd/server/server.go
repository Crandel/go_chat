package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/http/rest"
	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	sa "github.com/samonzeweb/godb/adapters/sqlite"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	// memory := memory.NewStorage()
	sql_db, _ := godb.Open(sa.Adapter, "./storage.db")
	sqlite_storage := sqlite.NewStorage(sql_db)
	ls := login.NewService(&sqlite_storage)
	sis := signin.NewService(&sqlite_storage)
	router := rest.InitHandlers(ls, sis)
	log.Fatal(http.ListenAndServe(":8080", router))
}
