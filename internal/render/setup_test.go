package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/marcelofranco/webapp-go-demo/internal/config"
	"github.com/marcelofranco/webapp-go-demo/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type responseWriter struct{}

func (tw *responseWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *responseWriter) Write(b []byte) (int, error) {
	lenght := len(b)
	return lenght, nil
}

func (tw *responseWriter) WriteHeader(statusCode int) {}
