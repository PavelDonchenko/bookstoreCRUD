package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/PavelDonchenko/40projects/go-bookstore/pkg/config"
	"github.com/PavelDonchenko/40projects/go-bookstore/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter, r *http.Request) {
	db := config.GetMySQLBase()
	bookModel := models.BookModel{DB: db}
	books, _ := bookModel.GetAllBooks()
	res, _ := json.Marshal(books)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	db := config.GetMySQLBase()
	vars := mux.Vars(r)
	bookId := vars["id"]
	id, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error while paring")
	}
	bookModel := models.BookModel{DB: db}
	book, _ := bookModel.GetBookById(id)
	res, _ := json.Marshal(book)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	db := config.GetMySQLBase()
	bookModel := models.BookModel{DB: db}
	book := models.Book{
		Id:          22,
		Name:        "Sleep",
		Author:      "John",
		Publication: "sdsd",
	}
	err := bookModel.CreateBook(&book)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Lastest id:", book.Id)
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	db := config.GetMySQLBase()
	bookModel := models.BookModel{DB: db}
	rows, err := bookModel.DeleteBook(6)
	if err != nil {
		fmt.Println(err)
	} else {
		if rows > 0 {
			fmt.Println("Done")
		}
	}
}
