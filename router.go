package main

import (
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ConfiguredRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = HandlerWithMiddleware(httpNotFoundHandler)

	router.
		Methods("GET", "HEAD").
		Path("/status").
		Handler(HandlerWithMiddleware(httpStatusHandler))

	router.
		Methods("POST").
		Path("/notifier_api/v2/notices").
		Handler(HandlerWithMiddleware(httpAirbrake23NotificationHandler))

	return router
}

func HandlerWithMiddleware(h http.HandlerFunc) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(raven.RecoveryHandler(h)))
}
