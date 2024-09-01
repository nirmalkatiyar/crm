package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nirmal/crm/controllers"
)

func DataExpImportRoutes(dataExpImportRoutes *gin.Engine) {
	// export data in csv or json format
	dataExpImportRoutes.GET("/export/customer_data", controller.ExportCustomerData())
	dataExpImportRoutes.POST("/import/customer_data", controller.ImportCustomerData())
}