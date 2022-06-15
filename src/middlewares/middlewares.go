package middlewares

import (
	"api/src/auth"
	"api/src/responses"
	"log"
	"net/http"
)

// Logger retrurn a custom log
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

// Authentication validates if the token is valid
func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := auth.TokenValidation(r); erro != nil {
			responses.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		next(w, r)
	}
}
