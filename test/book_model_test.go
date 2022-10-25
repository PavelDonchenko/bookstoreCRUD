package test

import (
	"log"
	"testing"

	"github.com/PavelDonchenko/bookstoreCRUD/pkg/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllBooks(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatalf("Error refreshing user and book table %v\n", err)
	}
	_, _, err = seedUsersAndBooks()
	if err != nil {
		log.Fatalf("Error seeding user and book  table %v\n", err)
	}
	books, err := bookInstance.GetAllBooks(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the books: %v\n", err)
		return
	}
	assert.Equal(t, len(*books), 2)
}

func TestSaveBook(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatalf("Error user and book refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newBook := models.Book{
		ID:         1,
		Name:       "This is the title",
		BookAuthor: "This is the content",
		UserID:     user.ID,
	}
	savedBook, err := newBook.CreateBook(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the book: %v\n", err)
		return
	}
	assert.Equal(t, newBook.ID, savedBook.ID)
	assert.Equal(t, newBook.Name, savedBook.Name)
	assert.Equal(t, newBook.BookAuthor, savedBook.BookAuthor)
	assert.Equal(t, newBook.UserID, savedBook.UserID)

}

func TestGetBookByID(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatalf("Error refreshing user and book table: %v\n", err)
	}
	book, err := seedOneUserAndOneBook()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundBook, err := bookInstance.GetBookById(server.DB, book.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundBook.ID, book.ID)
	assert.Equal(t, foundBook.Name, book.Name)
	assert.Equal(t, foundBook.BookAuthor, book.BookAuthor)
}

func TestUpdateABook(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatalf("Error refreshing user and book table: %v\n", err)
	}
	book, err := seedOneUserAndOneBook()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	bookUpdate := models.Book{
		ID:         1,
		Name:       "modiUpdate",
		BookAuthor: "modiupdate@gmail.com",
		UserID:     book.UserID,
	}
	updatedBook, err := bookUpdate.UpdateBook(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedBook.ID, bookUpdate.ID)
	assert.Equal(t, updatedBook.Name, bookUpdate.Name)
	assert.Equal(t, updatedBook.BookAuthor, bookUpdate.BookAuthor)
	assert.Equal(t, updatedBook.UserID, bookUpdate.UserID)
}

func TestDeleteABook(t *testing.T) {

	err := refreshUserAndBookTable()
	if err != nil {
		log.Fatalf("Error refreshing user and book table: %v\n", err)
	}
	book, err := seedOneUserAndOneBook()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := bookInstance.DeleteBook(server.DB, book.ID, book.UserID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
