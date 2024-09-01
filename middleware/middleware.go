package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/nirmal/crm/helpers"
)

// AuthenticateUser - middleware to authenticate user
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header found"})
			c.Abort()
			return
		}
		//	validate token
		claims, err := helper.ValidateUserToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
        // set claims in context
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("role", claims.Role)
		c.Set("uid", claims.Uid)
        // continue
		c.Next()
	}
}

// AuthenticateCustomer - middleware to authenticate customer
func AuthenticateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header found"})
			c.Abort()
			return
		}
		//	validate token
		claims, err := helper.ValidateCustomerToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		// set claims in context
		c.Set("cid", claims.Cid)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
        // continue
		c.Next()
	}
}
