package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nirmal/crm/middleware"
	"github.com/nirmal/crm/routes"
)

// Starting point of the application
func main() {

	PORT := os.Getenv("PORT")

	// if PORT is not set, use default port 8080
	if PORT == "" {
		PORT = "8080"
	}
	router := gin.New()

	router.Use(middleware.RateLimiterMiddleware())
	// middleware to log all requests on console
	router.Use(gin.Logger()) 

	// Authetication routes for user and customer
	routes.AuthenticateUserCustomer(router)

	// Customer routes
	routes.CustomerRoutes(router)

    // User routes
	routes.UserRoutes(router)

	//export/import data in csv or json format
    routes.DataExpImportRoutes(router)

	// Run the server on PORT
	router.Run(":"+PORT)
}
