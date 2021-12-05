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
	mem "github.com/Crandel/go_chat/pkg/storage/memory"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	memory := mem.NewStorage()
	aths := ath.NewService(&memory)
	adds := add.NewService(&memory)
	rdns := rdn.NewService(&memory)
	chts := cht.NewService(&memory)
	router := ntw.InitHandlers(aths, adds, rdns, chts)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
