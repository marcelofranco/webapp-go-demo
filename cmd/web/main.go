package main

import (
	"log"
	"net/http"

	"github.com/marcelofranco/webapp-go-demo/pkg/config"
	"github.com/marcelofranco/webapp-go-demo/pkg/handlers"
	"github.com/marcelofranco/webapp-go-demo/pkg/render"
)

const portNumber = ":8085"

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	log.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
