package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nirmal/crm/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
// SignedUserDetails : Signed user details
type SignedUserDetails struct {
	Email string
	Name  string
	Uid   string
	Role  string
	jwt.StandardClaims
}

const (
	userDatabaseName = "Cluster0"
	collectionName   = "users"
)

var UserCollection *mongo.Collection = database.OpenCollection(userDatabaseName, collectionName)
var USER_SECRET_KEY string = os.Getenv("USER_SECRET_KEY") // secret key

// GenerateUserToken : Generate a new user token
func GenerateUserToken(email, name, uId, role string) (signedToken string, err error) {
	// Create the Claims
	claims := &SignedUserDetails{
		Email: email,
		Name:  name,
		Uid:   uId,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(), // token expires in 24 hours
		},
	}
    // Create token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(USER_SECRET_KEY))
	if err != nil {
		msg := fmt.Sprintf("Error signing Token: %v", err)
		return "", errors.New(msg)
	}

	return token, nil
}

// UpdateUserToken : Update user token
func UpdateUserToken(signedToken, userId string) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
    // Update user token
	_, err := UserCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}

}

// ValidateUserToken : Validate user token
func ValidateUserToken(signedToken string) (claims *SignedUserDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedUserDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(USER_SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	// Validate token claims
	claims, ok := token.Claims.(*SignedUserDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}
	// Check if token has expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "the token has expired"
		return
	}
	return claims, msg
}
