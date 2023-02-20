package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/marcelofranco/webapp-go-demo/internal/config"
)

var app *config.AppConfig

// NewHelpers set up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Printf("Client error with a status of %d\n", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
