package user_router

import (
	user_controller "go-test/controllers/user"

	"github.com/gorilla/mux"
)

func RegisterUserRoute(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{id}", user_controller.GetUser).Methods("POST")
	userRouter.HandleFunc("/top/{id}", user_controller.GetTopUsers).Methods("GET")
	userRouter.HandleFunc("/", user_controller.CreateUser).Methods("POST")
	userRouter.HandleFunc("/friend/{id}", user_controller.GetFriendById).Methods("GET")
}
