package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

const (
	CONFIG_FILE = "config/airbrake2sentry.toml"
	APP_NAME    = "Airbrake2Sentry"
)

var config Config

type Config struct {
	Airbrake2Sentry ConfigAirbrake2Sentry
	Projects        map[string]ConfigAirbrakeProject `toml:"projects"`
}

type ConfigAirbrake2Sentry struct {
	ListenHost   string `toml:"listen_host"`
	ListenPort   int    `toml:"listen_port"`
	OwnSentryDSN string `toml:"own_sentry_dsn"`
}

type ConfigAirbrakeProject struct {
	AirbrakeAPIKey         string `toml:"airbrake_api_key"`
	SentryOrganizationName string `toml:"sentry_organization_name"`
	SentryProjectName      string `toml:"sentry_project_name"`
	SentryDSN              string `toml:"sentry_dsn"`
}

// Loads the configuration from the the default config file into the global
// config variable.
func (config *Config) Load() {
	if _, err := toml.DecodeFile(CONFIG_FILE, &config); err != nil {
		log.Fatalf("Fatal error reading config file: %s \n", err)
	}
}

// Returns the listen host and port in a string format that can be
// passed directly to http.ListenAndServe.
func (config *Config) ListenHostAndPort() string {
	return fmt.Sprintf("%s:%d", config.Airbrake2Sentry.ListenHost, config.Airbrake2Sentry.ListenPort)
}

// Searches the configuration for a project definition for given Airbrake
// API key, and returns the ConfigAirbrakeProject struct if found.
func (config *Config) FindProjectForAirbrakeAPIKey(APIKey string) (*ConfigAirbrakeProject, error) {
	for _, value := range config.Projects {
		if value.AirbrakeAPIKey == APIKey {
			return &value, nil
		}
	}

	return nil, fmt.Errorf("Could not find config block with Airbrake API Key '%s'", APIKey)
}
