package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	add "github.com/Crandel/go_chat/internal/adding"
	ath "github.com/Crandel/go_chat/internal/auth"
	cht "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
	ntw "github.com/Crandel/go_chat/internal/network"
	rdn "github.com/Crandel/go_chat/internal/reading"
	mem "github.com/Crandel/go_chat/internal/storage/memory"
)

const port = 8080

func main() {
	debug := os.Getenv("DEBUG")
	log := lg.InitLogger()
	log.PrintDebug = debug == "1"
	log.Log(lg.Info, "Starting server on port", port)

	memory := mem.NewStorage()
	aths := ath.NewService(&memory)
	adds := add.NewService(&memory)
	rdns := rdn.NewService(&memory)
	chts := cht.NewService(&memory)
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
