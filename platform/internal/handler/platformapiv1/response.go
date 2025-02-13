package platformapiv1

import (
	"encoding/json"
	"log"
	"net/http"

	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
)

func response(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	}
}

func errResponse(w http.ResponseWriter, _ *http.Request, code int, message string) {
	response(w, code, platformapiv1.Error{
		Message: message,
	})
}
