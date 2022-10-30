package main

import (
	"log"
	"os"

	"github.com/aflahahamed/ecommerce/controllers"
	"github.com/aflahahamed/ecommerce/database"
	"github.com/aflahahamed/ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router, app)
	log.Fatal(router.Run(":" + port))
}
