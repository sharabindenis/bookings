package render

import (
	"bytes"
	"fmt"
	"github.com/sharabindenis/bookings/pkg/config"
	"github.com/sharabindenis/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

// NewTemplates установка кнофигурации для пакета шаблонов
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData создает дефолтные данные для передачи во все шаблоны
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate отрисовка кеш шаблона
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// берем кеш из конфигурации приложения
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

// CreateTemplateCache сборка из разных tmpl файлов, страниц и слоёв в кеш
func CreateTemplateCache() (map[string]*template.Template, error) {
	// создается мап из tmpl файлов
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// выбираются все страницы
	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("Page is currently", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// попадается слой
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		// если слой есть добавить его в мап
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		// добавить в кеш все шаблоны
		myCache[name] = ts

	}
	//fmt.Println(myCache["about.page.tmpl"])
	return myCache, nil
}
