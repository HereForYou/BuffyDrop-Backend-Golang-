package setting_controller

import (
	"context"
	"fmt"
	"go-test/db"
	"go-test/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SetSetting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is setting controller set setting")
}

func GetAllSetting(w http.ResponseWriter, r *http.Request) {
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("setting")

	var setting models.Setting
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&setting)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, setting)
}
