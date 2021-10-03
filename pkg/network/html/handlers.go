package html

import (
	"fmt"
	"html/template"
	"net/http"
)

const TEMPLATE_ROOT = "resources/templates"

func RootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(TEMPLATE_ROOT + "/base.gtpl"))
		tmpl.Execute(w, nil)
	}
}

func StaticHandler() http.Handler {
	fmt.Println("Static handler")
	return http.FileServer(http.Dir("./resources/public"))
}
