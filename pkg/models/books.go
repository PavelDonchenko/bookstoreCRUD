package models

import "database/sql"

type Book struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publicate"`
}
type BookModel struct {
	DB *sql.DB
}

func (bm BookModel) GetAllBooks() ([]Book, error) {
	rows, err := bm.DB.Query("select * from testdb2")
	if err != nil {
		return nil, err
	} else {
		books := []Book{}
		for rows.Next() {
			var Id int64
			var Name string
			var Author string
			var Publication string
			err2 := rows.Scan(&Id, &Name, &Author, &Publication)
			if err2 != nil {
				return nil, err2
			} else {
				book := Book{Id, Name, Author, Publication}
				books = append(books, book)
			}
		}
		return books, nil
	}

}

func (bm BookModel) GetBookById(id int64) (Book, error) {
	rows, err := bm.DB.Query("select * from testdb2 where id = ?", id)
	if err != nil {
		return Book{}, err
	} else {
		book := Book{}
		for rows.Next() {
			var id int64
			var name string
			var author string
			var publication string
			err2 := rows.Scan(&id, &name, &author, &publication)
			if err2 != nil {
				return Book{}, err2
			} else {
				book = Book{id, name, author, publication}
			}
		}
		return book, nil
	}

}

func (bm BookModel) CreateBook(book *Book) error {
	result, err := bm.DB.Exec("insert into testdb2(Name, Author, Publication) values(?,?,?)", book.Name, book.Author, book.Publication)
	if err != nil {
		return err
	} else {
		book.Id, _ = result.LastInsertId()
		return nil
	}
}

func (bm BookModel) DeleteBook(id int64) (int64, error) {
	result, err := bm.DB.Exec("delete from testdb2 where id = ?", id)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
