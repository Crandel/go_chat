package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	"github.com/Crandel/go_chat/pkg/http/rest"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/storage/memory"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	memory := memory.NewStorage()
	aths := auth.NewService(&memory)
	adds := adding.NewService(&memory)
	rnms := reading.NewService(&memory)
	router := rest.InitHandlers(aths, adds, rnms)
	log.Fatal(http.ListenAndServe(":8080", router))
}
