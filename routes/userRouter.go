package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nirmal/crm/controllers"
	"github.com/nirmal/crm/middleware"
)

// UserRoutes - user related all routes
func UserRoutes(userRoutes *gin.Engine) {
	// middleware to authenticate user
	userRoutes.Use(middleware.AuthenticateUser())

	// user operations
	userRoutes.GET("/users", controller.GetUsers())
	userRoutes.GET("/users/:user_id", controller.GetUser())
	userRoutes.PUT("/users/:user_id", controller.UpdateUser())
	userRoutes.DELETE("/users/:user_id", controller.DeleteUser())

	// get all interactions, 
	userRoutes.GET("/users/meetings/", controller.GetAllInteractions())

	// get all interactions by user id
	userRoutes.GET("/user/meetings/", controller.GetInteractionsByUserID())

	// get all interactions by customer id
	userRoutes.POST("/users/meetings/:cust_id", controller.CreateInteractionAndSendEmail())

	// delete interaction by meet id
	userRoutes.DELETE("/users/meetings/:interaction_id", controller.DeleteInteraction())
}
