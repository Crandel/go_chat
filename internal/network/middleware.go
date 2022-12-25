package network

import (
	"context"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/internal/auth"
)

// Define authenticationMiddleware struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func NewAuthMiddleware() *authenticationMiddleware {
	tokenUsers := make(map[string]string)
	return &authenticationMiddleware{
		tokenUsers: tokenUsers,
	}
}

func (amw *authenticationMiddleware) Populate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxUser := &auth.CtxUser{}
		ctx = context.WithValue(ctx, auth.AuthKey, ctxUser)
		next.ServeHTTP(w, r.WithContext(ctx))
		authUserCtx := ctx.Value(auth.AuthKey)
		log.Println("authUserCtx", authUserCtx)
		if authUserCtx != nil {
			authUser := authUserCtx.(*auth.CtxUser)
			amw.tokenUsers[authUser.Token] = authUser.Nick
		}
	})
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
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
