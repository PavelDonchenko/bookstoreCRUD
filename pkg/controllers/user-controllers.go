package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/PavelDonchenko/bookstoreCRUD/api/responses"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/auth"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/models"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()

	err = user.Validate("")
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(s.DB)
	if err != nil {

		formattedError := utils.FormatError(err.Error())

		responses.JsonError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JsonFormat(w, http.StatusCreated, userCreated)
}

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(s.DB)
	if err != nil {
		responses.JsonError(w, http.StatusInternalServerError, err)
		return
	}
	responses.JsonFormat(w, http.StatusOK, users)
}

func (s *Server) GetUserById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}

	userGotten, err := user.FindUserByID(s.DB, uint32(uid))
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	responses.JsonFormat(w, http.StatusOK, userGotten)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.JsonError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()

	err = user.Validate("update")
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateAUser(s.DB, uint32(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.JsonError(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JsonFormat(w, http.StatusOK, updatedUser)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.JsonError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.DeleteAUser(s.DB, uint32(uid))
	if err != nil {
		responses.JsonError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JsonFormat(w, http.StatusNoContent, "")
}
