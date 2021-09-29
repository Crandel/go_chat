package main

import (
	"fmt"
	"log"
	"net/http"

	add "github.com/Crandel/go_chat/pkg/adding"
	ath "github.com/Crandel/go_chat/pkg/auth"
	"github.com/Crandel/go_chat/pkg/http/rest"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/storage/memory"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	memory := memory.NewStorage()
	aths := ath.NewService(&memory)
	adds := add.NewService(&memory)
	rdns := rdn.NewService(&memory)
	router := rest.InitHandlers(aths, adds, rdns)
	log.Fatal(http.ListenAndServe(":8080", router))
}
