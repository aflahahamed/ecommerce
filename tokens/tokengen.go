package tokens

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aflahahamed/ecommerce/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, First_Name string, lastname string, uid string) (signedToken string, signedrefreshtoken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: First_Name,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		fmt.Println("error 1", err)
		log.Panic()
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic()
		return "", "", err
	}

	return token, refreshtoken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		fmt.Println("test 1")
		msg = "token invalid"
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		fmt.Println("test 2")

		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		fmt.Println("test 3")
		msg = "token is already expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken string, signedrefreshtoken string, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateobj primitive.D
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedToken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})
	upsert := true
	filter := bson.M{"userd_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	},
		&opt)

	if err != nil {
		log.Panic(err)
	}
}
