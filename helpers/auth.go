package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nirmal/crm/utils"
)

// user can access this resource only via his token...
func CheckUserType(c *gin.Context, role string) (err error) {
	// get user role from context
	userRole := c.GetString("role")

	if userRole != role {
		err = fmt.Errorf("UnAuthenticated to access this resource")
		return err
	}
	return nil
}

// user can access this resource only via his token or he is an ADMIN...
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	// get user id & role from context
	uid := c.GetString("uid")
	userType := c.GetString("role")
   // check if user is ADMIN
	if err := CheckUserType(c, utils.ROLE_ADMIN); err == nil {
		return nil
	}
   // check if user is USER
	if uid == userId && userType == utils.ROLE_USER {
		return nil
	}
    // if not ADMIN or USER
	err = fmt.Errorf("UnAuthenticated to access this resource")
	return err
}

// customer can access this resource only via cid...
func MatchCustomerTypeToCid(c *gin.Context, customerId string) (err error) {
	cid := c.GetString("cid")

	if cid == customerId {
		return nil
	}
	err = fmt.Errorf("UnAuthenticated to access this resource")
	return err
}
