package user_controller

import (
	"context"
	"fmt"
	"go-test/db"
	"go-test/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("user")

	var user models.User

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

func GetTopUsers(w http.ResponseWriter, r *http.Request) {
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("user")

	// Get telegramId from params of request
	vars := mux.Vars(r)
	tgId := vars["id"]
	fmt.Println("tgId: ", tgId)

	// Get specific number of users from request query or 100 as default
	numUsersStr := r.URL.Query().Get("num")
	numUsers, err := strconv.Atoi(numUsersStr)
	if err != nil || numUsers < -1 {
		numUsers = 100
	}

	// Set the filter to find the users and sort by totalPoints
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"totalPoints", -1}})
	findOptions.SetLimit(int64(numUsers))

	// Return only bottom fields
	projection := bson.D{
		{"totalPoints", 1},
		{"userName", 1},
		{"tgId", 1},
		{"style", 1},
	}

	// Get all users according to filter and sort
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions.SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var users []models.User // all users

	for cursor.Next(context.TODO()) {
		var tempUser models.User
		if err := cursor.Decode(&tempUser); err != nil {
			log.Fatal(err)
		}
		users = append(users, tempUser)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	var curUser models.User
	err = collection.FindOne(context.TODO(), bson.D{{"tgId", tgId}}, options.FindOne().SetProjection(projection)).Decode(&curUser)
	if err != nil {
		log.Fatal(err)
	}

	totalMembers, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	ranking, err := collection.CountDocuments(context.TODO(), bson.D{{"totalPoints", bson.D{{"$gt", curUser.TotalPoints}}}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "Top Users: ", users, "\n", curUser, totalMembers, ranking)
}

func GetFriendById(w http.ResponseWriter, r *http.Request) {
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("user")

	vars := mux.Vars(r)
	tgId := vars["id"]

	var curUser models.User
	err := collection.FindOne(context.TODO(), bson.D{{"tgId", tgId}}).Decode(&curUser)
	if err != nil {
		log.Fatal(err)
	}

	var friendIds []string
	for _, friend := range curUser.Friends {
		friendIds = append(friendIds, friend.Id)
	}

	projection := bson.D{
		{"totalPoints", 1},
		{"userName", 1},
		{"tgId", 1},
		{"style", 1},
	}
	cursor, err := collection.Find(context.TODO(), bson.D{{"_id", bson.D{{"$in", curUser.Friends}}}}, options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.TODO()) {
		var tempUser models.User
		if err := cursor.Decode(&tempUser); err != nil {
			log.Fatal(err)
		}
		users = append(users, tempUser)
	}

	fmt.Fprint(w, "Friends invited by me: ", users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is user controller create user")
}
