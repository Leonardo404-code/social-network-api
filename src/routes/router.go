package routes

import (
	"api/src/routes/rotas"

	"github.com/gorilla/mux"
)

// Gerar retornar um router com as rotas configuradas
func GenerateRoutes() *mux.Router {
	r := mux.NewRouter()

	return rotas.Configure(r)
}
