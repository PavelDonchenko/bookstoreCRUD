package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/PavelDonchenko/bookstoreCRUD/pkg/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateBook(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		bookAuthor   string
		user_id      uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"The name", "book_author": "the bookAuthor", "user_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			name:         "The name",
			bookAuthor:   "the bookAuthor",
			user_id:      user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"The name", "book_author": "the bookAuthor", "user_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "name Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"name":"When no token is passed", "book_author": "the bookAuthor", "user_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"name":"When incorrect token is passed", "book_author": "the bookAuthor", "user_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"name": "", "book_author": "The bookAuthor", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name": "This is a name", "book_author": "", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			inputJSON:    `{"name": "This is an awesome name", "book_author": "the bookAuthor"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required User",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"name": "This is an awesome name", "book_author": "the bookAuthor", "user_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/books", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateBook)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["book_author"], v.bookAuthor)
			assert.Equal(t, responseMap["user_id"], float64(v.user_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetBooks(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndBooks()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAllBooks)
	handler.ServeHTTP(rr, req)

	var books []models.Book
	err = json.Unmarshal([]byte(rr.Body.String()), &books)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(books), 2)
}
func TestCGetBookByID(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatal(err)
	}
	book, err := seedOneUserAndOneBook()
	if err != nil {
		log.Fatal(err)
	}
	bookSample := []struct {
		id           string
		statusCode   int
		name         string
		bookAuthor   string
		user_id      uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(book.ID)),
			statusCode: 200,
			name:       book.Name,
			bookAuthor: book.BookAuthor,
			user_id:    book.UserID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range bookSample {

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetBookById)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, book.Name, responseMap["name"])
			assert.Equal(t, book.BookAuthor, responseMap["book_author"])
			assert.Equal(t, float64(book.UserID), responseMap["user_id"]) //the response author id is float64
		}
	}
}

func TestUpdateBook(t *testing.T) {

	var BookUserEmail, BookUserPassword string
	var AuthBookAuthorID uint32
	var AuthBookID uint64

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatal(err)
	}
	users, books, err := seedUsersAndBooks()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		BookUserEmail = user.Email
		BookUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(BookUserEmail, BookUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first book
	for _, book := range books {
		if book.ID == 2 {
			continue
		}
		AuthBookID = book.ID
		AuthBookAuthorID = book.UserID
	}
	// fmt.Printf("this is the auth book: %v\n", AuthBookID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		name         string
		bookAuthor   string
		user_id      uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"The updated book", "book_author": "This is the updated bookAuthor", "user_id": 1}`,
			statusCode:   200,
			name:         "The updated book",
			bookAuthor:   "This is the updated bookAuthor",
			user_id:      AuthBookAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"This is still another name", "book_author": "This is the updated bookAuthor", "user_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"This is still another name", "book_author": "This is the updated bookAuthor", "user_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			//Note: "Book 2" belongs to book 2, and name must be unique
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"Name 2", "book_author": "This is the updated bookAuthor", "user_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "name Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"", "book_author": "This is the updated bookAuthor", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"Awesome name", "bookAuthor": "", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"This is another name", "bookAuthor": "This is the updated bookAuthor"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthBookID)),
			updateJSON:   `{"name":"This is still another name", "bookAuthor": "This is the updated bookAuthor", "user_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/books", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateBook)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["book_author"], v.bookAuthor)
			assert.Equal(t, responseMap["user_id"], float64(v.user_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteBook(t *testing.T) {

	var BookUserEmail, BookUserPassword string
	var BookUserID uint32
	var AuthBookID uint64

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatal(err)
	}
	users, books, err := seedUsersAndBooks()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		BookUserEmail = user.Email
		BookUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(BookUserEmail, BookUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second book
	for _, book := range books {
		if book.ID == 1 {
			continue
		}
		AuthBookID = book.ID
		BookUserID = book.UserID
	}
	bookSample := []struct {
		id           string
		user_id      uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthBookID)),
			user_id:      BookUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthBookID)),
			user_id:      BookUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthBookID)),
			user_id:      BookUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			user_id:      1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range bookSample {

		req, _ := http.NewRequest("GET", "/books", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteBook)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
