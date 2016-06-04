package main

import (
	"net/url"
	"strings"

	"github.com/getsentry/raven-go"
)

// Constructs a fully populated raven Packet and SentryNotification for the
// passed in Airbrake23Notice.
func NewSentryNotificationFromAirbrake23Notice(airbrakeNotice *Airbrake23Notice) (*SentryNotification, error) {
	packet := BuildRavenPacketFromAirbrake23Notice(airbrakeNotice)
	return NewSentryNotification(airbrakeNotice.APIKey, packet)
}

// Constructs a fully populated raven Packet for the Airbrake23Notice,
// translating fields from
func BuildRavenPacketFromAirbrake23Notice(airbrakeNotice *Airbrake23Notice) *raven.Packet {
	eventID, _ := NewSentryEventID()
	return &raven.Packet{
		EventID:    eventID,
		Message:    airbrakeNotice.ErrorMessage,
		ServerName: airbrakeNotice.ServerEnvironment.Hostname,
		Interfaces: BuildRavenInterfacesFromAirbrake32Notice(airbrakeNotice),
		Tags:       BuildRavenTagsFromAirbrake23Notice(airbrakeNotice),
		Extra:      BuildRavenExtraFromAirbrake32Notice(airbrakeNotice)}
}

func BuildRavenInterfacesFromAirbrake32Notice(airbrakeNotice *Airbrake23Notice) []raven.Interface {
	interfaces := []raven.Interface{}

	// Append details of the exception to the exception interface
	// https://docs.getsentry.com/hosted/clientdev/interfaces/#failure-interfaces
	exception := &raven.Exception{
		Type:       airbrakeNotice.ErrorClass,
		Value:      airbrakeNotice.ErrorMessage,
		Stacktrace: BuildRavenStacktraceFromAirbrake32Notice(airbrakeNotice)}

	interfaces = append(interfaces, exception)

	// If exception generated as a result of a HTTP request, then append details of
	// the request to the request interface.
	// https://docs.getsentry.com/hosted/clientdev/interfaces/#context-interfaces
	if airbrakeNotice.RequestURL != "" {
		queryString := url.Values{}
		for _, v := range airbrakeNotice.RequestParams {
			queryString.Add(v.Name, v.Value)
		}

		env := make(map[string]string)
		headers := make(map[string]string)
		for _, v := range airbrakeNotice.RequestEnvironment {
			if strings.HasPrefix(v.Name, "HTTP_") {
				tidyName := strings.Replace(v.Name, "HTTP_", "", -1)
				tidyName = strings.Replace(tidyName, "_", "-", -1)
				headers[tidyName] = v.Value
			} else {
				env[v.Name] = v.Value
			}
		}

		request := &raven.Http{
			Method:  "UNKNOWN",
			Query:   queryString.Encode(),
			URL:     airbrakeNotice.RequestURL,
			Headers: headers,
			Cookies: headers["COOKIE"],
			Env:     env}
		interfaces = append(interfaces, request)
	}

	return interfaces
}

// Builds tags for the notification
func BuildRavenTagsFromAirbrake23Notice(airbrakeNotice *Airbrake23Notice) raven.Tags {
	tags := raven.Tags{}

	tags = append(tags, raven.Tag{Key: "environment", Value: airbrakeNotice.ServerEnvironment.EnvironmentName})

	return tags
}

// Builds extra information section for the notification.
func BuildRavenExtraFromAirbrake32Notice(airbrakeNotice *Airbrake23Notice) map[string]interface{} {
	extra := make(map[string]interface{})

	extra["Notified Via"] = APP_NAME
	extra["Origin Notifier Name"] = airbrakeNotice.NotifierName
	extra["Origin Notifier Version"] = airbrakeNotice.NotifierVersion

	return extra
}

// Builds a raven stacktrace for the notification based on the backtrace
// included in the Airbrake23Notice.
// https://docs.getsentry.com/hosted/clientdev/interfaces/#failure-interfaces
func BuildRavenStacktraceFromAirbrake32Notice(airbrakeNotice *Airbrake23Notice) *raven.Stacktrace {
	var frames []*raven.StacktraceFrame

	for _, line := range airbrakeNotice.Backtrace {
		filename := strings.Replace(line.File, airbrakeNotice.ServerEnvironment.ProjectRoot, "", -1)

		frame := &raven.StacktraceFrame{
			AbsolutePath: line.File,
			Filename:     filename,
			Lineno:       line.Number,
			Function:     line.Method,
			InApp:        false}
		frames = append(frames, frame)
	}

	return &raven.Stacktrace{frames}
}
