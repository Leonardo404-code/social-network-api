package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasPublicacoes = []Routes{
	{
		URI:                    "/publicacoes",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.CreatePost,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchPosts,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes/{publicacaoId}",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchPost,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes/{publicacaoId}",
		Metodo:                 http.MethodPut,
		Funcao:                 controllers.UpdatePost,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes/{publicacaoId}",
		Metodo:                 http.MethodDelete,
		Funcao:                 controllers.DeletePost,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/usuarios/{usuarioId}/publicacoes",
		Metodo:                 http.MethodGet,
		Funcao:                 controllers.SearchUserPublications,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes/{publicacaoId}/curtir",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.LikePost,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publicacoes/{publicacaoId}/descurtir",
		Metodo:                 http.MethodPost,
		Funcao:                 controllers.UnlikePost,
		RequiresAuthentication: true,
	},
}
