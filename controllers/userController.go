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
	"golang.org/x/crypto/bcrypt"
	"github.com/nirmal/crm/utils"
)

const (
	userdatabaseName   = "Cluster0"
	userCollectionName = "users"
)

var userValidate = validator.New()
var UserCollection *mongo.Collection = database.OpenCollection(userdatabaseName, userCollectionName)

//HashPassword: Hash password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword, foundUserPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(foundUserPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintln("email or password is incorrect")
		check = false
	}

	return check, msg
}

// UserSignUp godoc
func UserSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
        
		var user models.User
		// Bind incoming json to user model
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate user input
		validationErr := userValidate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
        // Check if email already exists
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for email"})
			log.Panic(err)
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}
        // Hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
        // Generate user token
		token, _ := helper.GenerateUserToken(*user.Email, *user.Name, user.UserId, *user.Role)
		user.Token = &token
        // Insert user into database
		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			msg := fmt.Sprintln("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// Return response
		c.JSON(http.StatusCreated, gin.H{"user":user,"message":"User created successfully"})
	}
}

// UserLogIn 
func UserLogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        // Find user by email
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}
        // Verify password user input with found user password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		// Generate user token
		token, err := helper.GenerateUserToken(*foundUser.Email, *foundUser.Name, foundUser.UserId, *foundUser.Role)
		if err != nil || token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update user token
		helper.UpdateUserToken(token, foundUser.UserId)

		err = UserCollection.FindOne(ctx, bson.M{utils.USER_ID: foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "user": foundUser,"msg":"User logged in successfully"})
	}
}

// GetUsers all users
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := helper.CheckUserType(c, utils.ROLE_ADMIN); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var users []models.User

		cursor, err := UserCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing users"})
			return
		}

		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding user data"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no users available"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// GetUser on
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param(utils.USER_ID)

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{utils.USER_ID: userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if user.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// UpdateUser
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param(utils.USER_ID)

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateObj := bson.M{}

		if user.Name != nil {
			updateObj["name"] = user.Name
		}

		if user.Email != nil {
			updateObj["email"] = user.Email
		}

		if user.Password != nil {
			password := HashPassword(*user.Password)
			updateObj["password"] = password
		}

		updateObj["updated_at"] = time.Now()

		filter := bson.M{utils.USER_ID: bson.M{"$eq": userId}}
		update := bson.M{"$set": updateObj}

		_, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
	}
}

// DeleteUser: Delete a user by utils.USER_ID
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param(utils.USER_ID)

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := UserCollection.DeleteOne(ctx, bson.M{utils.USER_ID: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting user"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found !!!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
