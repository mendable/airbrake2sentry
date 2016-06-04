package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/getsentry/raven-go"
)

// SentryNotification is a wrapper around a RavenPacket provided by
// the raven-go package, with additional information about the event ID
// and what Sentry Project the notification will be forwarded to.
type SentryNotification struct {
	EventID          string
	OrganizationName string
	ProjectName      string
	SentryDSN        string
	RavenPacket      *raven.Packet
}

// Returns a URL that can be used to directly find the notification in the
// Sentry web interface once the notification has been sent.
func (sentryNotification *SentryNotification) SearchURL() string {
	return fmt.Sprintf("https://app.getsentry.com/%s/%s/?query=%s",
		sentryNotification.OrganizationName,
		sentryNotification.ProjectName,
		sentryNotification.EventID)
}

// Sends the raven Packet for the notification to the SentryDSN.
func (sentryNotification *SentryNotification) Send() {
	log.Printf("Notifying Sentry: %s", sentryNotification.EventID)

	client, err := raven.New(sentryNotification.SentryDSN)
	if err != nil {
		log.Println(err)
		return
	}

	eventID, ch := client.Capture(sentryNotification.RavenPacket, nil)
	if err = <-ch; err != nil {
		log.Println(err)
		return
	}

	log.Printf("Sent error with id %s to Sentry", eventID)
}

// Returns a populated SentryNotication struct for a given raven packet
// airbrake API key.
func NewSentryNotification(airbrakeAPIKey string, packet *raven.Packet) (*SentryNotification, error) {
	project, err := config.FindProjectForAirbrakeAPIKey(airbrakeAPIKey)
	if err != nil {
		return nil, err
	}

	sentryNotification := &SentryNotification{
		EventID:          packet.EventID,
		RavenPacket:      packet,
		OrganizationName: project.SentryOrganizationName,
		ProjectName:      project.SentryProjectName,
		SentryDSN:        project.SentryDSN}

	return sentryNotification, nil
}

// Sentry requires a 16byte UUID, longer formatted v4 UUIDs that Airbrake uses
// are are not accepted. This is a copy of the `uuid()` method from raven-go /
// raven.Client as of June 2016.
// https://github.com/getsentry/raven-go/blob/dffeb57df75d6a911f00232155194e43d79d38d7/client.go#L221-L232
func NewSentryEventID() (string, error) {
	id := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		return "", err
	}
	id[6] &= 0x0F // clear version
	id[6] |= 0x40 // set version to 4 (random uuid)
	id[8] &= 0x3F // clear variant
	id[8] |= 0x80 // set to IETF variant
	return hex.EncodeToString(id), nil
}
