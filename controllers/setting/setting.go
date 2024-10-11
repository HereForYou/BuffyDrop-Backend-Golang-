package setting_controller

import (
	"context"
	"encoding/json"
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
	log.Println("This is api/setting/all GET request")
	client := db.Client
	collection := client.Database("BuffyDrop").Collection("setting")

	var setting models.Setting
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&setting)
	if err != nil {
		log.Fatal("🔴 " + err.Error())
	}

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		w.WriteHeader(http.StatusNoContent) // No content for preflight
		return
	}

	// Handle actual request
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		http.Error(w, "Internal Server Error -> Can not encode the setting", http.StatusBadRequest)
	}
}
