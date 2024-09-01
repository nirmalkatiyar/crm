package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nirmal/crm/database"
	"github.com/nirmal/crm/helpers"
	"github.com/nirmal/crm/models"
	"github.com/nirmal/crm/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	InteractionDatabaseName   = "Cluster0"
	InteractionCollectionName = "interactions"
)

var InteractionValidate = validator.New()
var InteractionCollection *mongo.Collection = database.OpenCollection(InteractionDatabaseName, InteractionCollectionName)
// CreateInteractionAndSendEmail 
func CreateInteractionAndSendEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var interaction models.Interaction
		if err := c.BindJSON(&interaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := userValidate.Struct(interaction)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		interaction.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		interaction.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		interaction.ID = primitive.NewObjectID()
		interaction.InteractionId = interaction.ID.Hex()

		userIDStr := c.GetString("uid")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		interaction.UserID = userID

		customerIDStr := c.Param(utils.CUSTOMER_ID)
		customerID, err := primitive.ObjectIDFromHex(customerIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}
		interaction.CustomerID = customerID

		var customer models.Customer
		err = CustomerCollection.FindOne(ctx, bson.M{utils.CUSTOMER_ID: customerIDStr}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if customer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}

		resultInsertionNumber, insertErr := InteractionCollection.InsertOne(ctx, interaction)
		if insertErr != nil {
			msg := fmt.Sprintln("fialed to create Interaction")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		userEmail := c.GetString("email")
		customerEmail := customer.Email

		errorChan := make(chan error, 2)

		// Send email to the user
		go func() {
			if err := utils.SendInteractionNotificationWithEmail(interaction, userEmail, interaction.StartTime.String()); err != nil {
				errorChan <- fmt.Errorf("failed to send email to user: %w", err)
			}
			close(errorChan) // Close channel when done
		}()

		// Send email to the customer
		go func() {
			if err := utils.SendInteractionNotificationWithEmail(interaction, *customerEmail, interaction.StartTime.String()); err != nil {
				errorChan <- fmt.Errorf("failed to send email to customer: %w", err)
			}
		}()

		go func() {
			for err := range errorChan {
				fmt.Println("Error:", err)
			}
		}()

		c.JSON(http.StatusCreated, resultInsertionNumber)
	}
}

// get all interactions (only admin can access)
func GetAllInteractions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := helpers.CheckUserType(c, utils.ROLE_ADMIN); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var interactions []models.Interaction

		cursor, err := InteractionCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing users"})
			return
		}

		if err = cursor.All(ctx, &interactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding user data"})
			return
		}

		if len(interactions) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no interactions available"})
			return
		}

		c.JSON(http.StatusOK, interactions)
	}
}

// get all interactions related to a user
func GetInteractionsByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.GetString("uid")

		if err := helpers.MatchUserTypeToUid(c, userIdStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var interactions []models.Interaction

		cursor, err := InteractionCollection.Find(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing interactions"})
			return
		}

		if err = cursor.All(ctx, &interactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding interaction data"})
			return
		}

		if len(interactions) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No interactions available related to this user"})
			return
		}

		c.JSON(http.StatusOK, interactions)
	}
}
// get all interactions related to a user
func DeleteInteraction() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		userIdStr := c.GetString("uid")

		if err := helpers.MatchUserTypeToUid(c, userIdStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		interactionIdStr := c.Param("interaction_id")
		interactionId, err := primitive.ObjectIDFromHex(interactionIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interaction ID"})
			return
		}

		// check if interaction exists and belongs to the user
		var interaction models.Interaction

		err = InteractionCollection.FindOne(ctx, bson.M{"_id": interactionId}).Decode(&interaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching interaction"})
			return
		}

		if interaction.UserID.Hex() != userIdStr {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UnAuthorized to delete this interaction"})
			return
		}

		result, err := InteractionCollection.DeleteOne(ctx, bson.M{"_id": interactionId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting interaction"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Interaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Interaction deleted successfully"})
	}
}
