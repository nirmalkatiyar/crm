package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/nirmal/crm/helpers"
	"github.com/nirmal/crm/models"
	"github.com/nirmal/crm/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// ExportData : Export data in csv or json format
func ExportCustomerData() gin.HandlerFunc {
	return func(c *gin.Context) {

		format := c.Query("format")
		// export data in csv or json format
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := helper.CheckUserType(c, utils.ROLE_ADMIN); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cursor, err := CustomerCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var customers []models.Customer
		if err = cursor.All(ctx, &customers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch format {
		case "json":
			c.JSON(http.StatusOK, customers)

		case "csv":
			c.Header("Content-Type", "text/csv")
			c.Header("Content-Disposition", "attachment;filename=customers.csv")

			writer := csv.NewWriter(c.Writer)
			defer writer.Flush()

			// Writing CSV headers
			err = writer.Write([]string{"ID", "CustomerId", "Name", "Email", "Password", "Company", "Phone", "CreatedAt", "UpdatedAt"})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Writing CSV records
			for _, customer := range customers {
				err = writer.Write([]string{
					customer.CustomerId,
					safeString(customer.Name),
					safeString(customer.Email),
					HashPassword(safeString(customer.Password)),
					safeString(customer.Company),
					safeString(customer.Phone),
					customer.CreatedAt.String(),
					customer.UpdatedAt.String(),
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					continue
				}
			}

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		}
	}
}

// only admin can import data
func ImportCustomerData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// import data in csv or json
		format := c.Query("format")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := helper.CheckUserType(c, utils.ROLE_ADMIN); err != nil {
			fmt.Println("Error:kdbv1", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		switch format {
		case "json":
			var newCustomersData []models.Customer
			if err := c.BindJSON(&newCustomersData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var docs []interface{}
			for _, customer := range newCustomersData {
				docs = append(docs, customer)
			}

			_, err := CustomerCollection.InsertMany(ctx, docs)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})

		case "csv":
			fmt.Println("Error:kdbv2")
			file, _, err := c.Request.FormFile("file")
			if err != nil {
				fmt.Println("Error:kdbv3", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("Error:kdbv4", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Skipping header row records[0]
			var docs []interface{}
			for _, record := range records[1:] {
				customer := models.Customer{
					Name:      &record[1],
					Email:     &record[2],
					Company:   &record[3],
					Phone:     &record[4],
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				docs = append(docs, customer)
			}

			_, err = CustomerCollection.InsertMany(ctx, docs)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		}

	}
}

func safeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
