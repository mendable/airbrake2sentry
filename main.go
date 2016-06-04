package main

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
)

func main() {
	log.Println("Starting Airbrake2Sentry")
	config.Load()

	raven.SetDSN(config.Airbrake2Sentry.OwnSentryDSN)

	log.Printf("Listening on %s...", config.ListenHostAndPort())
	log.Fatal(http.ListenAndServe(config.ListenHostAndPort(), ConfiguredRouter()))
}
