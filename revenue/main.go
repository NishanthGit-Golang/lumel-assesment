package main

import (
	"revenue/database"
	"revenue/routes"
)

func main() {

	database.InitDB()

	r := routes.SetupRouter()

	r.Run(":8080")
}
