package server

import (
	"github.com/PavelDonchenko/bookstoreCRUD/pkg/controllers"
)

var server = controllers.Server{}

func Run() {

	server.Initialize("mysql")

	server.Run(":8081")

}
