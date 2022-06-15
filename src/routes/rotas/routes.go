package rotas

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Rotas cria um struct de rotas padrão
type Routes struct {
	URI                    string
	Metodo                 string
	Funcao                 func(http.ResponseWriter, *http.Request)
	RequiresAuthentication bool
}

// Coloca todas as rotas dentro do router
func Configure(r *mux.Router) *mux.Router {
	Routes := rotasUsuarios
	Routes = append(Routes, rotaLogin)
	Routes = append(Routes, rotasPublicacoes...)

	for _, route := range Routes {
		if route.RequiresAuthentication {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(middlewares.Authentication(route.Funcao)),
			).Methods(route.Metodo)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Funcao)).Methods(route.Metodo)
		}
	}

	return r
}
