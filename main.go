package main

import (
	"fmt"

	"go-test/config"
	"go-test/db"
	"go-test/routes/setting"
	"go-test/routes/user"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	UserName    string  `bson:"userName"`
	TgId        string  `bson:"tgId"`
	Email       string  `bson:"email"`
	TotalPoints float32 `bson:"totalPoints"`
}

func main() {
	// load configuration data from config package
	cfg := config.LoadConfig()

	//================================================================================== setting router
	router := mux.NewRouter()
	user_router.RegisterUserRoute(router)
	setting_router.RegisterUserRoute(router)

	//================================================================================== Connect to DB
	db.Connect(cfg.DbUrl)

	fmt.Println("Server is running on port: ", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, router)
}
