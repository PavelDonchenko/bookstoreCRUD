package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PavelDonchenko/bookstoreCRUD/api/responses"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/auth"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/models"
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
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

	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		responses.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := s.SignIn(user.Email, user.Password)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.JsonError(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JsonFormat(w, http.StatusOK, token)
}

func (s *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = s.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
