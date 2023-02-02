package main

import (
	"github.com/gin-api-rest/database"
	"github.com/gin-api-rest/routes"
)

func main() {

	database.ConectaComBamcoDeDados()
	routes.HandleRequest()

}
