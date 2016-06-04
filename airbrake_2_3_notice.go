package main

import (
	"encoding/xml"
)

// https://help.airbrake.io/kb/api-2/notifier-api-v23
type Airbrake23Notice struct {
	XMLName xml.Name `xml:"notice"`

	// Required. The API key for the project that this error belongs to. The API key can
	// be found by viewing the edit project form on the Airbrake site.
	APIKey string `xml:"api-key"`

	// Required. The name of the notifier client submitting the request.
	NotifierName string `xml:"notifier>name"`

	// Required. The version number of the notifier client submitting the request.
	NotifierVersion string `xml:"notifier>version"`

	// Required. A URL at which more information can be obtained concerning the notifier client.
	NotifierURL string `xml:"notifier>url"`

	// Required. The class name or type of error that occurred.
	ErrorClass string `xml:"error>class"`

	// Optional. A short message describing the error that occurred.
	ErrorMessage string `xml:"error>message"`

	// Required. This element can occur more than once.
	Backtrace []Airbrake23NoticeBacktraceLine `xml:"error>backtrace>line"`

	// Required only if there is a request element. The URL at which the error occurred.
	RequestURL string `xml:"request>url"`

	// Required only if there is a request element. The component in which the error occurred.
	// In model-view-controller frameworks like Rails, this should be set to the controller.
	// Otherwise, this can be set to a route or other request category.
	RequestComponent string `xml:"request>component"`

	// Optional. The action in which the error occurred. If each request is routed to a
	// controller action, this should be set here. Otherwise, this can be set to a method
	// or other request subcategory.
	RequestAction string `xml:"request>action"`

	// Optional. A list of var elements describing request parameters from the query string,
	// POST body, routing, and other inputs. See the section on var elements below.
	RequestParams []Airbrake23NoticeVar `xml:"request>params>var"`

	// Optional. A list of var elements describing session variables from the request.
	// See the section on var elements below.
	RequestSession []Airbrake23NoticeVar `xml:"request>session>var"`

	// Optional. A list of var elements describing CGI variables from the request, such
	// as SERVER_NAME and REQUEST_URI. See the section on var elements below.
	RequestEnvironment []Airbrake23NoticeVar `xml:"request>cgi-data>var"`

	// Server environment details.
	ServerEnvironment Airbrake23NoticeServerEnvironment `xml:"server-environment"`
}

// Each line element describes one code location or frame in the backtrace when
// the error occurred, and requires @file and @number attributes. If the
// location includes a method or function, the @method attribute should be used.
type Airbrake23NoticeBacktraceLine struct {
	XMLName xml.Name `xml:"line"`
	Method  string   `xml:"method,attr"`
	File    string   `xml:"file,attr"`
	Number  int      `xml:"number,attr"`
}

type Airbrake23NoticeServerEnvironment struct {
	XMLName xml.Name `xml:"server-environment"`

	// Optional. The path to the project in which the error occurred, such as RAILS_ROOT.
	ProjectRoot string `xml:"project-root"`

	// Required. The name of the server environment in which the error occurred, such as "production."
	EnvironmentName string `xml:"environment-name"`

	// Optional. The version of the application that this error came from. If the App Version is
	// set on the project, then errors older than the project's app version will be ignored.
	// This version field uses Semantic Versioning style versioning.
	AppVersion string `xml:"app-version"`

	// Server hostname if available? (NB: Undocumented?)
	Hostname string `xml:"hostname"`
}

// The params, session, and cgi-data elements can contain one or more var
// elements for each parameter or variable that was set when the error occured.
// Each var element should have a @key attribute for the name of the variable,
// and element text content for the value of the variable.
type Airbrake23NoticeVar struct {
	XMLName xml.Name `xml:"var"`
	Name    string   `xml:"key,attr"`
	Value   string   `xml:",chardata"`
}

// Api response returned to Airbrake 2.3 clients when they POST a notification
// to us to forward to sentry.
type Airbrake23ResponseNotice struct {
	XMLName xml.Name `xml:"notice"`
	Id      string   `xml:"id"`
	Url     string   `xml:"url"`
}

// Turns a HTTP Body into a Airbrake23Notice struct.
func NewAirbrake23Notice(request_body []byte) (*Airbrake23Notice, error) {
	airbrakeNotice := &Airbrake23Notice{}
	err := xml.Unmarshal(request_body, &airbrakeNotice)
	return airbrakeNotice, err
}
