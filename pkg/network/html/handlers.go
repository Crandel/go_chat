package html

import "net/http"

const TEMPLATE_ROOT = "resources/templates"

func RootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
