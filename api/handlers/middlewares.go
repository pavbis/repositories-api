package handlers

import (
	"net/http"
)

func BasicAuthMiddleware(user, pass string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkBasicAuth(r, user, pass) {
			next(w, r)
			return
		}

		respond(w, http.StatusUnauthorized, []byte("Unauthorized"))
	}
}

func checkBasicAuth(r *http.Request, user, pass string) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return u == user && p == pass
}
