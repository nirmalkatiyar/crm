package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nirmal/crm/controllers"
)

// AuthenticateUserCustomer - user and customer authentication routes
func AuthenticateUserCustomer(request *gin.Engine) {

	// customer authentication
	request.POST("customer/signup", controller.CustomerSignUp())
	request.POST("customer/signin", controller.CustomerLogIn())

	// user authentication
	request.POST("user/signup", controller.UserSignUp())
	request.POST("user/signin", controller.UserLogIn())

}
