# Airbrake2Sentry

A simple HTTP proxy that captures API requests intended for Airbrake, translates
them into Sentry API format, and forwards them onto a Sentry DSN.

This is really useful if you have many old applications that cannot be easily
switched to using Sentry directly, or, you just want to try Sentry out to see if
it's for you, and don't want to have to do a load of client updates to try it.

Compatible with the [Airbrake 2.3 API](https://help.airbrake.io/kb/api-2/notifier-api-v23)
format, other versions not yet supported.

## Setup & Installation

```
  cp config/airbrake2sentry.toml.example config/airbrake2sentry.toml
  go get
```

## Run locally in development

```
  go run *.go
```

## Manual Testing

To experiment with some test notifications manually, use the curl examples
provided in the examples directory.

```
  ./examples/notification.sh
```
