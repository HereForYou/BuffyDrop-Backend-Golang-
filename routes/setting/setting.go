package setting_router

import (
	setting_controller "go-test/controllers/setting"

	"github.com/gorilla/mux"
)

func RegisterUserRoute(r *mux.Router) {
	settingRouter := r.PathPrefix("/api/setting").Subrouter()
	settingRouter.HandleFunc("", setting_controller.SetSetting).Methods("GET")
	settingRouter.HandleFunc("/all", setting_controller.GetAllSetting).Methods("GET")
}
