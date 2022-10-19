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

func (server *Server) CreateBook(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	book := models.Book{}

	err = json.Unmarshal(body, &book)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}
	book.Prepare()
	err = book.Validate()
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != book.UserID {
		responses.JsonError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := book.CreateBook(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.JsonError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JsonFormat(w, http.StatusCreated, postCreated)
}

func (server *Server) GetAllBooks(w http.ResponseWriter, r *http.Request) {

	book := models.Book{}

	books, err := book.GetAllBooks(server.DB)
	if err != nil {
		responses.JsonError(w, http.StatusInternalServerError, err)
		return
	}
	responses.JsonFormat(w, http.StatusOK, books)
}

func (server *Server) GetBookById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}
	book := models.Book{}

	bookReceived, err := book.GetBookById(server.DB, bid)
	if err != nil {
		responses.JsonError(w, http.StatusInternalServerError, err)
		return
	}
	responses.JsonFormat(w, http.StatusOK, bookReceived)
}

func (server *Server) UpdateBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the book id is valid
	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	//Check if the user token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Check if the book exist
	book := models.Book{}

	err = server.DB.Debug().Model(models.Book{}).Where("id = ?", bid).Take(&book).Error
	if err != nil {
		responses.JsonError(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	// If a user attempt to update a book not belonging to him
	if uid != book.UserID {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// Read the data posted
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	bookUpdate := models.Book{}
	err = json.Unmarshal(body, &bookUpdate)
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != bookUpdate.UserID {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	bookUpdate.Prepare()
	err = bookUpdate.Validate()
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookUpdate.ID = book.ID //this is important to tell the model the book id to update, the other update field are set above

	postUpdated, err := bookUpdate.UpdateBook(server.DB)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.JsonError(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JsonFormat(w, http.StatusOK, postUpdated)
}

func (server *Server) DeleteBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid book id given to us?
	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the book exist
	book := models.Book{}
	err = server.DB.Debug().Model(models.Book{}).Where("id = ?", bid).Take(&book).Error
	if err != nil {
		responses.JsonError(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this book?
	if uid != book.UserID {
		responses.JsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = book.DeleteBook(server.DB, bid, uid)
	if err != nil {
		responses.JsonError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", bid))
	responses.JsonFormat(w, http.StatusNoContent, "")
}
