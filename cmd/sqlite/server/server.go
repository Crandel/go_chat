package main

import (
	"fmt"
	"log"
	"net/http"

	a "github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/http/rest"
	l "github.com/Crandel/go_chat/pkg/login"
	r "github.com/Crandel/go_chat/pkg/reading"
	s "github.com/Crandel/go_chat/pkg/signin"
	sql "github.com/Crandel/go_chat/pkg/storage/sqlite"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

func main() {
	port := 8080
	fmt.Println("Starting server on port", port)
	sql_db, _ := godb.Open(sqlite.Adapter, "./storage.db")
	sqlite_storage := sql.NewStorage(sql_db)
	ls := l.NewService(&sqlite_storage)
	sis := s.NewService(&sqlite_storage)
	as := a.NewService(&sqlite_storage)
	rs := r.NewService(&sqlite_storage)
	router := rest.InitHandlers(ls, sis, as, rs)
	log.Fatal(http.ListenAndServe(":8080", router))
}
