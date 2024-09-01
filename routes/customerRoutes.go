package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nirmal/crm/controllers"
	"github.com/nirmal/crm/middleware"
)

// CustomerRoutes - routes for customer
func CustomerRoutes(customerRoutes *gin.Engine) {
    // middleware to authenticate customer
	customerRoutes.Use(middleware.AuthenticateCustomer())

	// customer operations
	customerRoutes.GET("/customers", controller.GetAllCustomers())
	customerRoutes.GET("/customers/:cust_id", controller.GetCustomer())
	customerRoutes.PUT("/customers/:cust_id", controller.UpdateCustomer())
	customerRoutes.DELETE("/customers/:cust_id", controller.DeleteCustomer())

	// get all tickets
	customerRoutes.GET("/customers/tickets/", controller.GetAllTickets())
	customerRoutes.POST("/customers/ticket/:interaction_id", controller.CreateTicket())
	customerRoutes.GET("/customers/ticket/:user_id", controller.GetTicketsByUserID())
	customerRoutes.PUT("/customers/ticket/:ticket_id", controller.UpdateTicket())
	customerRoutes.DELETE("/customers/ticket/:ticket_id", controller.DeleteTicket())
}
