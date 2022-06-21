package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	ath "github.com/Crandel/go_chat/pkg/auth"
	cht "github.com/Crandel/go_chat/pkg/chatting"
	ntw "github.com/Crandel/go_chat/pkg/network"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	sql "github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

const port = 8080

func main() {
	fmt.Println("Starting server on port", port)
	sqlDB, _ := godb.Open(sqlite.Adapter, "./storage.db")
	sqliteStorage := sql.NewStorage(sqlDB)
	aths := ath.NewService(&sqliteStorage)
	adds := add.NewService(&sqliteStorage)
	rdns := rdn.NewService(&sqliteStorage)
	chts := cht.NewService(&sqliteStorage)
	go chts.Run()
	router := ntw.NewRouter(aths, adds, rdns, chts)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
