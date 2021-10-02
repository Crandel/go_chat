package main

import (
	"fmt"
	"log"
	"net/http"

	add "github.com/Crandel/go_chat/pkg/adding"
	ath "github.com/Crandel/go_chat/pkg/auth"
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
	router := ntw.InitHandlers(aths, adds, rdns)
	log.Fatal(http.ListenAndServe(":8080", router))
}
