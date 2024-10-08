package setting_router

import (
	"github.com/gorilla/mux"
	"go-test/controllers/setting"
)

func RegisterUserRoute(r *mux.Router) {
	settingRouter := r.PathPrefix("/setting").Subrouter()
	settingRouter.HandleFunc("", setting_controller.SetSetting).Methods("GET")
}
