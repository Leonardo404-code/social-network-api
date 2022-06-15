package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasUsuarios = []Routes{
	{
		URI:                    "/usuarios",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.CreateUser,
		RequiresAuthentication: false,
	},
	{
		URI:                    "/usuarios",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchUsers,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}",
		Metodo:                 http.MethodPut,
		Funcao:                 controllers.UpdateUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}",
		Metodo:                 http.MethodDelete,
		Funcao:                 controllers.DeleteUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}/seguir",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.FollowUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}/parar-de-seguir",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.StopFollowingUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}/seguidores",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchFollowers,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}/atualizar-senha",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.UpdatePassword,
		RequiresAuthentication: true,
	},
}
