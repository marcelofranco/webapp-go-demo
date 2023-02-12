package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	tp, _ := template.ParseFiles("./templates/" + tmpl)
	err := tp.Execute(w, nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}
