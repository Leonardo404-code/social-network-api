package routes

import (
	"api/src/routes/rotas"

	"github.com/gorilla/mux"
)

// GenerateRoutes Returns a configured route
func GenerateRoutes() *mux.Router {
	r := mux.NewRouter()

	return rotas.Configure(r)
}
