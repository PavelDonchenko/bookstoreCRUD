package server

import (
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/controllers"
)

var server = controllers.Server{}

func Run() {
	// Init Database
	server.Initialize("mysql")

	//Run server
	server.Run(":8081")

}
