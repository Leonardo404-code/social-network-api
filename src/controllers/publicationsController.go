package controllers

import (
	"api/src/auth"
	banco "api/src/database"
	"api/src/models"
	"api/src/repository"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, erro := auth.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publication models.Publication

	if erro = json.Unmarshal(body, &publication); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publication.AutorID = userID

	if erro = publication.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	publication.ID, erro = repo.CreatePost(publication)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, publication)
}

func SearchPosts(w http.ResponseWriter, r *http.Request) {
	userID, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)
	publicacoes, erro := repo.SearchPost(userID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacoes)
}

func SearchPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	publicationID, erro := strconv.ParseUint(params["publicacaoId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	publication, erro := repo.SearchPostByID(publicationID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publication)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)

	publicationID, erro := strconv.ParseUint(params["publicacaoId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	publicationInDatabase, erro := repo.SearchPostByID(publicationID)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if publicationInDatabase.AutorID != userID {
		responses.Erro(w, http.StatusForbidden, errors.New(`
		Não é possivel atualizar uma publicação que não seja sua!
		`))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var publication models.Publication

	if erro = json.Unmarshal(body, &publication); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = publication.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repo.UpdatePost(publicationID, publication); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)

	publicationID, erro := strconv.ParseUint(params["publicacaoId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	publicationInDatabase, erro := repo.SearchPostByID(publicationID)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if publicationInDatabase.AutorID != userID {
		responses.Erro(w, http.StatusForbidden, errors.New(`
		Não é possivel atualizar uma publicação que não seja sua!
		`))
		return
	}

	if erro = repo.DeletePost(publicationID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func SearchUserPublications(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId, erro := strconv.ParseUint(params["usuarioId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)
	publications, erro := repo.SearchUserPublications(userId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publications)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	publicationID, erro := strconv.ParseUint(params["publicacaoId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	if erro = repo.LikePost(publicationID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	publicationID, erro := strconv.ParseUint(params["publicacaoId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewPublicationRepository(db)

	if erro = repo.UnlikePost(publicationID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
