package network

import (
	"context"
	"net/http"

	lg "github.com/Crandel/go_chat/internal/logging"

	"github.com/Crandel/go_chat/internal/auth"
)

var log = lg.InitLogger()

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
		ctx = context.WithValue(ctx, auth.AuthKey, ctxUser) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
		authUserCtx := ctx.Value(auth.AuthKey)
		log.Log(lg.Debug, "authUserCtx", authUserCtx)
		if authUserCtx != nil {
			authUser := authUserCtx.(*auth.AuthUser)
			amw.tokenUsers[authUser.Token] = authUser.Nick
		}
		log.Log(lg.Debug, amw.tokenUsers)
	})
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Log(lg.Debug, amw.tokenUsers)
		username, password, ok := r.BasicAuth()
		log.Log(lg.Debug, username, password, ok)
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Logf(lg.Debug, "User from Authorization header: %s\n", username)
		token := auth.MakeToken(username, password)
		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			ctx := r.Context()
			ctxUser := &auth.AuthUser{
				Nick:  user,
				Token: token,
			}
			ctx = context.WithValue(ctx, auth.AuthKey, ctxUser) //nolint:staticcheck

			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
