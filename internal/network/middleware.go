package network

import (
	"log"
	"net/http"
)

// Define our struct
type AuthenticationMiddleware struct {
	tokenUsers map[string]string
}

func NewAuthMiddleware() *AuthenticationMiddleware {
	tokenUsers := make(map[string]string)
	return &AuthenticationMiddleware{
		tokenUsers: tokenUsers,
	}
}

func (amw *AuthenticationMiddleware) Populate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		token := r.Context().Value("token")
		nick := r.Context().Value("nick")
		log.Println("Populate token and nick", token, nick)
		if token != nil && nick != nil {
			amw.tokenUsers[token.(string)] = nick.(string)
		}
	})
}

// Middleware function, which will be called for each request
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
