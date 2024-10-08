package main

import (
	"fmt"
	// "go-test/utils"
	"github.com/gorilla/mux"
	"go-test/config"
	setting_router "go-test/routes/setting"
	user_router "go-test/routes/user"
	"net/http"
)

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, world! This is Go test project!")
// }

func main() {
	cfg := config.LoadConfig()

	userRouter := mux.NewRouter()
	user_router.RegisterUserRoute(userRouter)
	setting_router.RegisterUserRoute(userRouter)

	fmt.Println("Server is running on port: ", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, userRouter)
	// var limit int
	// utils.SayHello("SmartFox")
	// fmt.Print("Enter a specific number: ")
	// fmt.Scan(&limit)
	// utils.FindEvens(limit)
	// http.HandleFunc("/", helloHandler)
	// fmt.Println("Server is running on port 8080!")
	// http.ListenAndServe(":8080", nil)
}
