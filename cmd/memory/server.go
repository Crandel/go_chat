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
	"github.com/Crandel/go_chat/pkg/storage/memory"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	memory := memory.NewStorage()
	ls := login.NewService(&memory)
	sis := signin.NewService(&memory)
	as := adding.NewService(&memory)
	rs := reading.NewService(&memory)
	router := rest.InitHandlers(ls, sis, as, rs)
	log.Fatal(http.ListenAndServe(":8080", router))
}
