package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/http/rest"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port ", port)
	router := rest.InitHandlers()
	log.Fatal(http.ListenAndServe(":8080", router))
}
