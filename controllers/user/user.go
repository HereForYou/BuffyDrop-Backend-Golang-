package user_controller

import (
	"context"
	"time"

	"encoding/json"
	"fmt"
	"go-test/db"
	"go-test/models"
	"go-test/utils"

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

type GetUserResponse struct {
	User       models.User `json:"user"`
	SignIn     bool        `json:"signIn"`
	RemainTime float32     `json:"remainTime"`
	CycleTime  float32     `json:"cycleTime"`
}

var cycleTime = 10

func GetUser(w http.ResponseWriter, r *http.Request) {
	//======================================================================== Get the tgId from params
	vars := mux.Vars(r)
	tgId := vars["id"]
	inviteLink := tgId
	fmt.Println("This is telegram Id", inviteLink)

	//======================================================================== Get data from POST request (Content-Type == x-www-form-urlencoded)
	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// req := GetUserRequest{
	// 	UserName:   r.FormValue("userName"),
	// 	FirstName:  r.FormValue("firstName"),
	// 	LastName:   r.FormValue("lastName"),
	// 	StartParam: r.FormValue("start_param"),
	// 	Style:      r.FormValue("style"),
	// }

	//======================================================================== Get data from POST request (Content-Type == application/json)
	var req GetUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//======================================================================== Connecting to user and setting collection of BuffyDrop database
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
			log.Printf("🔴 Error finding setting" + err.Error())
		}
	}

	//======================================================================== Finding user by telegram Id
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No document found in user collection by telegram Id from request: %s", tgId)
		} else {
			log.Printf("🔴 Error while finding document: %v" + err.Error())
		}
	}

	//======================================================================== If user exists in database
	if user.TgId != "" {
		//======================================================================================================================== Calculate elapsed time since start farming
		nomarlizedDateStr, err := utils.NormalizeDateString(user.StartFarming.String())
		if err != nil {
			log.Printf("🔴 Error normalizing date" + err.Error())
			return
		}

		start, err := time.Parse("2006-01-02 15:04:05.000 -0700 MST", nomarlizedDateStr)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		now := time.Now()
		countTime := now.Sub(start).Seconds()

		//======================================================================================================================== If farming is ended
		if countTime > float64(cycleTime) {
			user.Cliamed = false
			_, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{
				"$set": user,
			})
			if err != nil {
				http.Error(w, "Internal Server Error: "+err.Error(), http.StatusBadRequest)
				return
			}
			response := GetUserResponse{
				User:       user,
				SignIn:     true,
				RemainTime: 0,
				CycleTime:  float32(cycleTime),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Internal Server Error: "+err.Error(), http.StatusBadRequest)
			}
			//======================================================================================================================== while farming
		} else {
			response := GetUserResponse{
				User:       user,
				SignIn:     true,
				RemainTime: float32(countTime),
				CycleTime:  float32(cycleTime),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Internal Server Error: "+err.Error(), http.StatusBadRequest)
			}
		}
		//============================================================================================================================ If user does not exist in database (when user is new)
	} else {
		inviteRevenue := setting.InviteRevenue
		rankCount, err := userCollection.CountDocuments(context.TODO(), bson.D{})
		if err != nil {
			http.Error(w, "Internal Server Error -> counting documents in user collection", http.StatusBadRequest)
		}
		totalPoints := rankCount + 1

		//======================================================================================================================== When user is invited
		if req.StartParam != "" {
			var inviter models.User
			if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": req.StartParam}).Decode(&inviter); err != nil {
				if err == mongo.ErrNoDocuments {
					log.Println("No document found in user collection by start_param")
					http.Error(w, "Unauthorized invitation link", http.StatusBadRequest)
				} else {
					log.Printf("Error: finding iniviter by start_param: %v", err.Error())
				}
			}

			newUser := models.User{
				TgId:         tgId,
				UserName:     req.UserName,
				FirstName:    req.FirstName,
				LastName:     req.LastName,
				IsInvited:    true,
				InviteLink:   tgId,
				TotalPoints:  float64(totalPoints),
				JoinRank:     int(totalPoints),
				Style:        req.Style,
				StartFarming: time.Now(),
				LastLogin:    time.Now(),
				Friends:      []models.Friend{},
				Task:         []string{},
			}
			_, err := userCollection.InsertOne(context.TODO(), newUser)
			if err != nil {
				log.Printf("🔴 " + err.Error())
				http.Error(w, "Internal Server Error while inserting new user", http.StatusBadRequest)
			}

			if !utils.HasFriendWithId(inviter.Friends, newUser.TgId) {
				inviter.Friends = append(inviter.Friends, models.Friend{Id: newUser.TgId, Revenue: inviteRevenue * newUser.TotalPoints})
				inviter.TotalPoints += inviteRevenue * newUser.TotalPoints
				if len(inviter.Friends) != 0 && len(inviter.Friends)%3 == 0 {
					inviter.TotalPoints += 200000
				}
				if _, err := userCollection.UpdateByID(context.TODO(), inviter.Id, bson.M{
					"$set": inviter,
				}); err != nil {
					log.Printf("🔴 " + err.Error())
					http.Error(w, "Internal Server Error while updating the inviter", http.StatusBadRequest)
				}
			}

			//======================================================================================================================== Sending response to client
			if err := json.NewEncoder(w).Encode(GetUserResponse{User: newUser, SignIn: false, RemainTime: 0, CycleTime: float32(cycleTime)}); err != nil {
				log.Printf("🔴 " + err.Error())
				http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
			}
			//======================================================================================================================== When user is not invited and is new user
		} else {
			newUser := models.User{
				TgId:         tgId,
				UserName:     req.UserName,
				FirstName:    req.FirstName,
				LastName:     req.LastName,
				InviteLink:   tgId,
				TotalPoints:  float64(totalPoints),
				JoinRank:     int(totalPoints),
				Style:        req.Style,
				StartFarming: time.Now(),
				LastLogin:    time.Now(),
				Friends:      []models.Friend{},
				Task:         []string{},
			}

			if _, err := userCollection.InsertOne(context.TODO(), newUser); err != nil {
				log.Printf("🔴 Error inserting new user: %v" + err.Error())
				http.Error(w, "Internal Server Error while saving new user into mongoDB", http.StatusInternalServerError)
				return
			}

			var reward float64
			if totalPoints < 11 {
				reward = 0.1001
			} else if totalPoints < 101 {
				reward = 0.1
			} else if totalPoints < 1001 {
				reward = 0.096
			} else if totalPoints < 10001 {
				reward = 0.0949
			} else if totalPoints < 100001 {
				reward = 0.065
			} else if totalPoints < 1000001 {
				reward = 0.019
			} else {
				reward = float64(totalPoints) * 0.01
			}
			log.Println("This is reward -> ", reward)
			setting.InviteRevenue = reward

			if _, err := settingCollection.UpdateOne(context.TODO(), bson.M{"taskList": setting.TaskList}, bson.M{
				"$set": setting,
			}); err != nil {
				log.Printf("🔴 " + err.Error())
				http.Error(w, "Internal Server Error while saving setting", http.StatusBadRequest)
			}

			if err := json.NewEncoder(w).Encode(GetUserResponse{User: newUser, SignIn: false, RemainTime: 0, CycleTime: float32(cycleTime)}); err != nil {
				log.Printf("🔴 " + err.Error())
				http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
			}
		}
	}
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
	findOptions.SetSort(bson.M{"totalPoints": -1})
	findOptions.SetLimit(int64(numUsers))

	// Return only bottom fields
	projection := bson.M{
		"totalPoints": 1,
		"userName":    1,
		"tgId":        1,
		"style":       1,
	}

	// Get all users according to filter and sort
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions.SetProjection(projection))
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}
	defer cursor.Close(context.TODO())

	var users []models.User // all users

	for cursor.Next(context.TODO()) {
		var tempUser models.User
		if err := cursor.Decode(&tempUser); err != nil {
			log.Printf("🔴 " + err.Error())
		}
		users = append(users, tempUser)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("🔴 " + err.Error())
	}

	var curUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"tgId": tgId}, options.FindOne().SetProjection(projection)).Decode(&curUser)
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}

	totalMembers, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}

	ranking, err := collection.CountDocuments(context.TODO(), bson.M{"totalPoints": bson.M{"$gt": curUser.TotalPoints}})
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}

	if err := json.NewEncoder(w).Encode(struct {
		CurUser      models.User   `json:"curUser"`
		TopUsers     []models.User `json:"topUsers"`
		TotalMembers int           `json:"totalMembers"`
		Ranking      int           `json:"ranking"`
	}{
		TopUsers:     users,
		CurUser:      curUser,
		TotalMembers: int(totalMembers),
		Ranking:      int(ranking + 1),
	}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func GetFriendById(w http.ResponseWriter, r *http.Request) {
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("user")

	vars := mux.Vars(r)
	tgId := vars["id"]

	var curUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&curUser)
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}

	var friendIds []string
	for _, friend := range curUser.Friends {
		friendIds = append(friendIds, friend.Id)
	}
	fmt.Println("FriendIds", friendIds)

	projection := bson.M{
		"totalPoints": 1,
		"userName":    1,
		"tgId":        1,
		"style":       1,
		"revenue":     1,
	}
	friends, err := collection.Find(context.TODO(), bson.M{"tgId": bson.M{"$in": friendIds}}, options.Find().SetProjection(projection))
	if err != nil {
		log.Printf("🔴 " + err.Error())
	}
	defer friends.Close(context.TODO())

	var users []struct {
		Info    models.User
		Revenue int `json:"revenue"`
	}
	for friends.Next(context.TODO()) {
		var tempUser models.User
		if err := friends.Decode(&tempUser); err != nil {
			log.Printf("🔴 " + err.Error())
		}
		for _, friend := range curUser.Friends {
			if friend.Id == tempUser.TgId {
				users = append(users, struct {
					Info    models.User
					Revenue int "json:\"revenue\""
				}{Info: tempUser, Revenue: int(friend.Revenue)})
			}
		}
	}

	if err := json.NewEncoder(w).Encode(struct {
		InviteLink  string `json:"inviteLink"`
		FriendsInfo []struct {
			Info    models.User
			Revenue int `json:"revenue"`
		} `json:"friendsInfo"`
	}{
		InviteLink:  curUser.InviteLink,
		FriendsInfo: users,
	}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is user controller create user")
}

func ClaimFarming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tgId := vars["id"]
	fmt.Println("This is telegram Id for ClaimFarming function", tgId)

	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")
	settingCollection := client.Database("BuffyDrop").Collection("setting")

	var user models.User
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while finding user by TG id", http.StatusBadRequest)
	}
	var setting models.Setting
	if err := settingCollection.FindOne(context.TODO(), bson.D{}).Decode(&setting); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while finding setting", http.StatusBadRequest)
	}

	user.TotalPoints += setting.DailyRevenue * float64(cycleTime)
	user.Cliamed = true
	user.IsStarted = false

	if _, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{"$set": user}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while saving user by TG id", http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(struct {
		Status     bool        `json:"status"`
		User       models.User `json:"user"`
		RemainTime float32     `json:"remainTime"`
	}{
		Status:     true,
		User:       user,
		RemainTime: float32(cycleTime),
	}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func StartFarming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tgId := vars["id"]
	fmt.Println("This is telegram Id for StartFarming functioin", tgId)

	//=========================================================================================================================== Connection DB
	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")

	var user models.User
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents in database by telegram Id")
			http.Error(w, "Invalid user", http.StatusBadRequest)
			return
		} else {
			log.Printf("🔴 " + err.Error())
		}
	}
	fmt.Println("This is user", user.UserName)

	user.StartFarming = time.Now()
	user.Cliamed = true
	user.IsStarted = true

	if _, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{"$set": user}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while updating user documnet", http.StatusBadRequest)
	}
	fmt.Println("Successfully saved!")

	if err := json.NewEncoder(w).Encode(struct {
		User      models.User `json:"user"`
		CycleTime int         `json:"cycleTime"`
	}{User: user, CycleTime: cycleTime}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func EndFarming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tgId := vars["id"]
	fmt.Println("This is TG id", tgId)

	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")

	var user models.User
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while finding user", http.StatusBadRequest)
	}

	user.Cliamed = false
	if _, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{"$set": user}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while saving user", http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(struct {
		User      models.User `json:"user"`
		CycleTime int         `json:"cycleTime"`
	}{
		User:      user,
		CycleTime: cycleTime,
	}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func Tap(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tgId := vars["id"]
	fmt.Println("This is TG id for tap game:", tgId)

	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")

	var user models.User
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		log.Printf("🔴 " + err.Error())
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
		}
	}

	user.TotalPoints += 1

	if _, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{"$set": user}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while saving user", http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(struct {
		Status bool        `json:"status"`
		User   models.User `json:"user"`
	}{
		Status: true,
		User:   user,
	}); err != nil {
		log.Printf("🔴 " + err.Error())
		http.Error(w, "Internal Server Error while sending response to client", http.StatusBadRequest)
	}
}

func HandleFollow(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tgId := vars["id"]
	fmt.Println("This is TG id from api/user/task request", tgId)

	var reqBody struct {
		Id     string  `json:"id"`
		Profit float32 `json:"profit"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	client := db.Client
	userCollection := client.Database("BuffyDrop").Collection("user")

	var user models.User
	if err := userCollection.FindOne(context.TODO(), bson.M{"tgId": tgId}).Decode(&user); err != nil {
		log.Printf("🔴 " + err.Error())
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
		}
		return
	}

	if !utils.ContainsValue(user.Task, "", reqBody.Id) {
		user.Task = append(user.Task, reqBody.Id)
		user.TotalPoints += float64(reqBody.Profit)

		if _, err := userCollection.UpdateOne(context.TODO(), bson.M{"tgId": tgId}, bson.M{"$set": user}); err != nil {
			log.Printf("🔴 " + err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Printf("🔴 " + err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		if err := json.NewEncoder(w).Encode(false); err != nil {
			log.Printf("🔴 " + err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
