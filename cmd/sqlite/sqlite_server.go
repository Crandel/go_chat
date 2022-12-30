package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Crandel/go_chat/cmd/sqlite/migrations"
	add "github.com/Crandel/go_chat/internal/adding"
	ath "github.com/Crandel/go_chat/internal/auth"
	cht "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
	ntw "github.com/Crandel/go_chat/internal/network"
	rdn "github.com/Crandel/go_chat/internal/reading"
	sql "github.com/Crandel/go_chat/internal/storage/sqlite"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

const port = 8080

func main() {
	debug := os.Getenv("DEBUG")
	migrationFolder := os.Getenv("MIGRATIONS")
	if migrationFolder == "" {
		migrationFolder = "./migrations"
	}
	log := lg.InitLogger()
	log.PrintDebug = debug == "1"
	log.Println("Starting server on port", port)
	sqlDB, _ := godb.Open(sqlite.Adapter, "./storage.db")
	sqlDB.SetLogger(log)
	err := migrations.RunMigrations(sqlDB, migrationFolder)
	if err != nil {
		log.Fatal(err)
		return
	}
	sqliteStorage := sql.NewStorage(sqlDB)
	aths := ath.NewService(&sqliteStorage)
	adds := add.NewService(&sqliteStorage)
	rdns := rdn.NewService(&sqliteStorage)
	chts := cht.NewService(&sqliteStorage)
	go chts.Run()
	router := ntw.NewRouter(aths, adds, rdns, chts)
	srv := &http.Server{
		Addr:         ":" + fmt.Sprint(port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
