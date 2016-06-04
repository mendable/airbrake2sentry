package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
)

// Not found fallback handler, so we may perform custom notification logic
// when unexpected requests are receive, which may indicate a problem.
func httpNotFoundHandler(response http.ResponseWriter, request *http.Request) {
	msg := fmt.Sprintf("Unexpected HTTP request: %s %s : This could indicate an unsupported"+
		" client is attempting to relay notifications", request.Method, request.URL.Path)

	log.Println(msg)
	raven.CaptureMessage(msg, nil)

	response.WriteHeader(http.StatusNotFound)
}

// GET /status
func httpStatusHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}

// POST /notifier_api/v2/notices
func httpAirbrake23NotificationHandler(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	airbrakeNotice, err := NewAirbrake23Notice(body)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", airbrakeNotice)

	sentryNotification, err := NewSentryNotificationFromAirbrake23Notice(airbrakeNotice)
	if err != nil {
		panic(err)
	}

	go sentryNotification.Send()

	// Return an Airbrake 2.3 compliant API response.
	responseNotice := &Airbrake23ResponseNotice{
		Id:  fmt.Sprintf("%s", sentryNotification.EventID),
		Url: sentryNotification.SearchURL()}

	bodyXML, err := xml.MarshalIndent(responseNotice, "", "  ")
	if err != nil {
		panic(err)
	}

	response.Header().Set("Content-Type", "text/xml")
	response.Write(bodyXML)
}
