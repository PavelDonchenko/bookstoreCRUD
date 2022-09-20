package routes

import (
	"github.com/PavelDonchenko/40projects/go-bookstore/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/books/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/books/", controllers.GetBook).Methods("GET")
	router.HandleFunc("/books/{id}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/books/{id}", controllers.DeleteBook).Methods("DELETE")
	//router.HandleFunc("/books/{id}", controllers.UpdateBook).Methods("PUT")
}
