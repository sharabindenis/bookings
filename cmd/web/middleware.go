package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf добавляет CSRF токен для всех POST запросов
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad загружает и сохраняет сессию на каждый запрос
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
