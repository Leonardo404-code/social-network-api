package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotaLogin = Routes{
	URI:                    "/login",
	Metodo:                 http.MethodPost,
	Funcao:                 controllers.Login,
	RequiresAuthentication: false,
}
