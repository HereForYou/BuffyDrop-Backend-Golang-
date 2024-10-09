package user_controller

import (
	"context"
	// "encoding/json"
	"fmt"
	"go-test/db"
	"go-test/models"
	// "io/ioutil"
	"log"
	"net/http"
	"strconv"

	// routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetUserRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	UserName   string `json:"userName"`
	StartParam string `json:"start_param"`
	Style      string `json:"style"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	//======================================================================== Get the tgId from params
	vars := mux.Vars(r)
	tgId := vars["id"]
	inviteLink := tgId
	fmt.Println("This is telegram Id", inviteLink)

	//======================================================================== Get data from request
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println(string(body))
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("This is form", r.FormValue("userName"))
	req := GetUserRequest{
		UserName:   r.FormValue("userName"),
		FirstName:  r.FormValue("firstName"),
		LastName:   r.FormValue("lastName"),
		StartParam: r.FormValue("start_param"),
		Style:      r.FormValue("style"),
	}

	log.Println("This is request: ", req)

	// err = json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// err = json.Unmarshal(body, &req)
	// if err != nil {
	//     http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
	//     return
	// }

	// var c routing.Context
	// if err = c.Read(&req); err != nil {
	// 	http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
	// 	return
	// }

	// fmt.Println("This is req.data", req.UserName)
	//======================================================================== Connecting to user collection of BuffyDrop database
	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")
	settingCollection := client.Database("BuffyDrop").Collection("setting")

	var user models.User
	var setting models.Setting
	err := settingCollection.FindOne(context.TODO(), bson.D{}).Decode(&setting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
		} else {
			log.Fatal(err)
		}
	}

	err = userCollection.FindOne(context.TODO(), bson.D{{"tgId", tgId}}).Decode(&user)
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
