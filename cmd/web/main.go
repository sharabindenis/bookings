package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/sharabindenis/bookings/pkg/config"
	"github.com/sharabindenis/bookings/pkg/handlers"
	"github.com/sharabindenis/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8180"

// создаем переменную с конфигурацией
var app config.AppConfig
var session *scs.SessionManager

func main() {
	//изменить на тру в продакшене
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// создаем переменную с кешем
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create temple cache")
	}
	// создаем кеш для текущего приложения
	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting app on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
