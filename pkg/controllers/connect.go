package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PavelDonchenko/bookstoreCRUD/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(Dbdriver string) {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	if Dbdriver == os.Getenv("MYSQL_DBDRIVER") {
		dns := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
		s.DB, err = gorm.Open(Dbdriver, dns)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.Book{}) //database migration

	s.Router = mux.NewRouter()

	s.RegisterBookStoreRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Println("Listening to port 8081")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
