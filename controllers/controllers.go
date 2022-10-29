package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aflahahamed/ecommerce/database"
	"github.com/aflahahamed/ecommerce/models"
	token "github.com/aflahahamed/ecommerce/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "Users")
var productCollection *mongo.Collection = database.ProductData(database.Client, "Product")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or password incorrrect"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := token.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserter := userCollection.InsertOne(ctx, user)
		if inserter != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}
		c.JSON(http.StatusCreated, "successfully signed in!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the login pass or user incorrcet"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login or password incorrrect"})
			fmt.Println(msg)
			return
		}
		signedToken, refreshToken, err := token.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the token was not generated"})
			return
		}
		defer cancel()
		token.UpdateAllTokens(signedToken, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusFound, foundUser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := productCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not inserted"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Sucessfully added")

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		cursor, err := productCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try again later")
			return
		}

		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)
		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")
		fmt.Println("test match", queryParam)

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQueryDb, err := productCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(400, "something went wrong")
			return
		}
		err = searchQueryDb.All(ctx, &searchProducts)
		fmt.Println("result", searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchQueryDb.Close(ctx)

		if err := searchQueryDb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid requets")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}
