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

// CreatePost Validates user inputs and creates a new post
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

// SearchPosts Get a param in URL and search by Post
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

// Find a single post with the ID
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

// Get Post ID, validate inputs and Update in database
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
		N??o ?? possivel atualizar uma publica????o que n??o seja sua!
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

// DeletePost get the post ID in URL and Delete
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
		N??o ?? possivel atualizar uma publica????o que n??o seja sua!
		`))
		return
	}

	if erro = repo.DeletePost(publicationID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// SearchUserPublications search the user's own posts
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

// LikePost add a like in a Post
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

// UnlikePost remove a like in a Post
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
