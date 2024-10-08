package user_controller

import (
	"context"
	"fmt"
	"go-test/db"
	"go-test/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	client := db.Client

	collection := client.Database("BuffyDrop").Collection("user")

	err := collection.FindOne(context.TODO(), bson.D{{"tgId", "7202566339"}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
		} else {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(w, user.UserName)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is user controller create user")
}
