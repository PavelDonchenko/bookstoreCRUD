package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string    `gorm:"size:255;not null;unique" json:"name"`
	BookAuthor string    `gorm:"size:255;not null;" json:"book_author"`
	User       User      `json:"user"`
	UserID     uint32    `gorm:"not null" json:"user_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (b *Book) Prepare() {
	b.ID = 0
	b.Name = html.EscapeString(strings.TrimSpace(b.Name))
	b.BookAuthor = html.EscapeString(strings.TrimSpace(b.BookAuthor))
	b.User = User{}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
}

func (b *Book) Validate() error {

	if b.Name == "" {
		return errors.New("Required Name")
	}
	if b.BookAuthor == "" {
		return errors.New("Required Author")
	}
	if b.UserID < 1 {
		return errors.New("Required User")
	}
	return nil
}

func (b *Book) CreateBook(db *gorm.DB) (*Book, error) {
	var err error
	err = db.Debug().Model(&Book{}).Create(&b).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.UserID).Take(&b.User).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) GetAllBooks(db *gorm.DB) (*[]Book, error) {
	var err error
	books := []Book{}
	err = db.Debug().Model(&Book{}).Limit(100).Find(&books).Error
	if err != nil {
		return &[]Book{}, err
	}
	if len(books) > 0 {
		for i, _ := range books {
			err := db.Debug().Model(&User{}).Where("id = ?", books[i].UserID).Take(&books[i].User).Error
			if err != nil {
				return &[]Book{}, err
			}
		}
	}
	return &books, nil
}

func (b *Book) GetBookById(db *gorm.DB, bid uint64) (*Book, error) {
	var err error
	err = db.Debug().Model(&Book{}).Where("id = ?", bid).Take(&b).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.UserID).Take(&b.User).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) UpdateBook(db *gorm.DB) (*Book, error) {

	var err error

	err = db.Debug().Model(&Book{}).Where("id = ?", b.ID).Updates(Book{Name: b.Name, BookAuthor: b.BookAuthor, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.UserID).Take(&b.User).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) DeleteBook(db *gorm.DB, bid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Book{}).Where("id = ? and user_id = ?", bid, uid).Take(&Book{}).Delete(&Book{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Book not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
