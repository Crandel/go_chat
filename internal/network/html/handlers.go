package html

import (
	"html/template"
	"net/http"
)

const TEMPLATE_ROOT = "resources/templates"

func RootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(TEMPLATE_ROOT+"/base.gtpl", TEMPLATE_ROOT+"/room.gtpl"))
		tmpl.Execute(w, nil)
	}
}

func StaticHandler() http.Handler {
	dir := http.Dir("./resources/public")
	return http.FileServer(dir)
}
