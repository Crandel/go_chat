package network

import (
	"context"
	"net/http"

	lg "github.com/Crandel/go_chat/internal/logging"

	"github.com/Crandel/go_chat/internal/auth"
)

var log = lg.Logger

// Define authenticationMiddleware struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func NewAuthMiddleware(aths auth.Service) *authenticationMiddleware {
	tokenUsers := make(map[string]string)

	users := aths.ReadAuthUsers()
	for _, user := range users {
		tokenUsers[user.Token] = user.Nick
	}
	return &authenticationMiddleware{
		tokenUsers: tokenUsers,
	}
}

func (amw *authenticationMiddleware) Populate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxUser := &auth.AuthUser{}
		ctx = context.WithValue(ctx, auth.AuthKey, ctxUser)
		next.ServeHTTP(w, r.WithContext(ctx))
		authUserCtx := ctx.Value(auth.AuthKey)
		log.Debugln("authUserCtx", authUserCtx)
		if authUserCtx != nil {
			authUser := authUserCtx.(*auth.AuthUser)
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
			log.Debugf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
