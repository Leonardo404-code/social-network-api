package controllers

import (
	"api/src/auth"
	"api/src/models"
	"api/src/repository"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	database "api/src/database"

	"github.com/gorilla/mux"
)

// CreateUser Validate all the fields and create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)

		return
	}

	var user models.User

	if erro = json.Unmarshal(body, &user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)

		return
	}

	if erro = user.Prepare("cadastro"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)

		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)

	user.ID, erro = repo.CreateUser(user)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

// SearchUsers Search by multiple users
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)

	user, erro := repo.SearchUsers(nameOrNick)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

// SearchUser Search for one user
func SearchUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 32)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)

	user, erro := repo.SearchUserByID(userID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

// UpdateUser valid all fields and update a user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIDInToken, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if userIDInToken != userID {
		responses.Erro(w, http.StatusForbidden, errors.New(
			"n??o ?? possivel atualizar um usu??rio que n??o seja o seu",
		))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var user models.User

	if erro = json.Unmarshal(body, &user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = user.Prepare("edicao"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repository.NewRepositoryUser(db)
	if erro = repositorio.UpdateUser(userID, user); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// DeleteUser get parameter id and delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIDInToken, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if userID != userIDInToken {
		responses.Erro(w, http.StatusForbidden, errors.New(
			"n??o ?? possivel deletar um usuario que n??o seja seu",
		))
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)

	if erro = repo.DeleteUser(userID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// FollowUser Add a follower to a user
func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(w, http.StatusForbidden, errors.New("n??o ?? possivel seguir voc?? mesmo"))
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)
	if erro = repo.Follow(userID, followerID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, nil)
}

// StopFollowingUser Remove a follower to a user
func StopFollowingUser(w http.ResponseWriter, r *http.Request) {
	followerID, erro := auth.ExtractUserID(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(w, http.StatusForbidden, errors.New(
			"n??o ?? possivel parar de seguir voc?? mesmo",
		))
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)
	if erro = repo.StopFollowing(userID, followerID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// SearchFollowers Look for followers of the authenticated user
func SearchFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)

	followers, erro := repo.SearchFollowers(userID)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

// UpdatePassword update user password
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIDInToken, erro := auth.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if userIDInToken != userID {
		responses.Erro(w, http.StatusForbidden, errors.New(`
			N??o ?? possivel atualizar a senha de outro usu??rio
		`))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	var password models.Password

	if erro = json.Unmarshal(body, &password); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repo := repository.NewRepositoryUser(db)
	passwordInDB, erro := repo.SearchPassword(userID)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = security.VerifyPassword(passwordInDB, password.Old); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New(`
		a senha atual n??o condiz com a que est?? salva no banco
		`))
		return
	}

	passwordHash, erro := security.Hash(password.New)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repo.UpdatePassword(userID, string(passwordHash)); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
