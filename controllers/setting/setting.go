package setting_controller

import (
	"fmt"
	"net/http"
)

func SetSetting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is setting controller set setting")
}
