package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	ath "github.com/Crandel/go_chat/pkg/auth"
	ntw "github.com/Crandel/go_chat/pkg/network"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	sql "github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

const port = 8080

func main() {
	fmt.Println("Starting server on port", port)
	sql_db, _ := godb.Open(sqlite.Adapter, "./storage.db")
	sqlite_storage := sql.NewStorage(sql_db)
	aths := ath.NewService(&sqlite_storage)
	adds := add.NewService(&sqlite_storage)
	rdns := rdn.NewService(&sqlite_storage)
	router := ntw.InitHandlers(aths, adds, rdns)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
