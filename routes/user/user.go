package user_router

import (
	"github.com/gorilla/mux"
	"go-test/controllers/user"
)

func RegisterUserRoute(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", user_controller.GetUser).Methods("GET")
	userRouter.HandleFunc("/", user_controller.CreateUser).Methods("POST")
}
