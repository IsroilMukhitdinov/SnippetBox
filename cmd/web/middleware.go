package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		response.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("X-Frame-Options", "deny")
		response.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(response, request)
	}

	return http.HandlerFunc(fn)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.RequestURI)

		next.ServeHTTP(response, request)
	}

	return http.HandlerFunc(fn)
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				response.Header().Set("Connection", "Close")
				app.serverError(response, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(response, request)
	}

	return http.HandlerFunc(fn)
}
