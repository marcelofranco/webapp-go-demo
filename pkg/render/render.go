package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/marcelofranco/webapp-go-demo/pkg/config"
	"github.com/marcelofranco/webapp-go-demo/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate renders template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatalf("Could not find template %s\n", tmpl)
		}
	}

	t, tmplFound := tc[tmpl]
	if !tmplFound {
		log.Fatalf("Could not find template %s\n", tmpl)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	err = t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	matches, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return myCache, err
	}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		pageName := filepath.Base(page)
		ts, err := template.New(pageName).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[pageName] = ts
	}

	return myCache, nil
}
