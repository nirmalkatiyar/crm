package controllers

import (
	"context"

	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nirmal/crm/database"
	helper "github.com/nirmal/crm/helpers"
	"github.com/nirmal/crm/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/nirmal/crm/utils"
)

const (
	customerdatabaseName = "Cluster0"
	collectionName       = "customers"
)

var customerValidate = validator.New()
var CustomerCollection *mongo.Collection = database.OpenCollection(customerdatabaseName, collectionName)

// CustomerSignUp : Create a new customer
func CustomerSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Bind incoming json to customer model
		var customer models.Customer
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate incoming data
		validationErr := customerValidate.Struct(customer)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// Check if email already exists
		count, err := CustomerCollection.CountDocuments(ctx, bson.M{"email": customer.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for email"})
			log.Panic(err)
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		// Hash password
		password := HashPassword(*customer.Password)
		customer.Password = &password

		// Create new customer
		customer.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		customer.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		customer.ID = primitive.NewObjectID()
		customer.CustomerId = customer.ID.Hex()
		// Generate customer token
		token, _ := helper.GenerateCustomerToken(*customer.Email, *customer.Name, customer.CustomerId)
		customer.Token = &token

		resultInsertionNumber, insertErr := CustomerCollection.InsertOne(ctx, customer)
		if insertErr != nil {
			msg := fmt.Sprintln("Customer item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"insertId": resultInsertionNumber, "message": "Customer created successfully"})
	}
}

// CustomerLogIn : Customer login
func CustomerLogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Bind incoming json to customer model
		var customer models.Customer
		var foundCustomer models.Customer

		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//check the record with the email in DB
		err := CustomerCollection.FindOne(ctx, bson.M{"email": customer.Email}).Decode(&foundCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}
		//check if the password is correct
		passwordIsValid, msg := VerifyPassword(*customer.Password, *foundCustomer.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundCustomer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}
		// Generate customer token
		token, err := helper.GenerateCustomerToken(*foundCustomer.Email, *foundCustomer.Name, foundCustomer.CustomerId)
		if err != nil || token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Update customer token
		helper.UpdateCustomerToken(token, foundCustomer.CustomerId)
		// Find customer by utils.CUSTOMER_ID
		err = CustomerCollection.FindOne(ctx, bson.M{utils.CUSTOMER_ID: foundCustomer.CustomerId}).Decode(&foundCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"customer": foundCustomer, "message": "Customer logged in successfully"})
	}
}

// GetAllCustomers : Get all customers
func GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var customers []models.Customer
		// Find all customers
		cursor, err := CustomerCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing customers"})
			return
		}
		// Decode all customers
		if err = cursor.All(ctx, &customers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding customer data"})
			return
		}

		if len(customers) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no customers available"})
			return
		}

		c.JSON(http.StatusOK, customers)
	}
}

// GetCustomer : Get a customer by ID
func GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param(utils.CUSTOMER_ID)
		// Match customer type to utils.CUSTOMER_ID
		if err := helper.MatchCustomerTypeToCid(c, customerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Find customer by utils.CUSTOMER_ID
		var customer models.Customer
		err := CustomerCollection.FindOne(ctx, bson.M{utils.CUSTOMER_ID: customerId}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		// Check if customer email is nil
		if customer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}

		c.JSON(http.StatusOK, customer)
	}
}

// UpdateCustomer : Update a customer by ID
func UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param(utils.CUSTOMER_ID)
		// Match customer type to utils.CUSTOMER_ID
		if err := helper.MatchCustomerTypeToCid(c, customerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var customer models.Customer
		// Bind incoming json to customer model
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateObj := bson.M{}
		// Update customer fields
		if customer.Name != nil {
			updateObj["name"] = customer.Name
		}

		if customer.Email != nil {
			updateObj["email"] = customer.Email
		}

		if customer.Company != nil {
			updateObj["company"] = customer.Company
		}

		if customer.Phone != nil {
			updateObj["phone"] = customer.Phone
		}

		if customer.Password != nil {
			password := HashPassword(*customer.Password)
			updateObj["password"] = password
		}

		updateObj["updated_at"] = time.Now()

		filter := bson.M{utils.CUSTOMER_ID: bson.M{"$eq": customerId}}
		update := bson.M{"$set": updateObj}
		// Update customer
		_, err := CustomerCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "customer updated successfully"})
	}
}

// DeleteCustomer : Delete a customer by ID
func DeleteCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param(utils.CUSTOMER_ID)
         // Match customer type to utils.CUSTOMER_ID
		if err := helper.MatchCustomerTypeToCid(c, customerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.M{utils.CUSTOMER_ID: bson.M{"$eq": customerId}}
		// Delete customer
		_, err := CustomerCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "customer deleted successfully"})
	}
}
