package bugsnag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

var (
	APIKey              string
	AppVersion          string
	OSVersion           string
	ReleaseStage        = "production"
	NotifyReleaseStages = []string{"production"}
	AutoNotify          = true
	UseSSL              = false
	Verbose             = false
	Notifier            = &bugsnagNotifier{
		Name:    "Bugsnag Go client",
		Version: "0.0.1",
		URL:     "https://github.com/toggl/bugsnag_client",
	}
)

type (
	bugsnagNotifier struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		URL     string `json:"url"`
	}
	bugsnagPayload struct {
		APIKey   string           `json:"apiKey"`
		Notifier *bugsnagNotifier `json:"notifier"`
		Events   []*bugsnagEvent  `json:"events"`
	}
	bugsnagException struct {
		ErrorClass string              `json:"errorClass"`
		Message    string              `json:"message,omitempty"`
		Stacktrace []bugsnagStacktrace `json:"stacktrace,omitempty"`
	}
	bugsnagStacktrace struct {
		File       string `json:"file"`
		LineNumber string `json:"lineNumber"`
		Method     string `json:"method"`
		InProject  bool   `json:"inProject,omitempty"`
	}
	bugsnagEvent struct {
		UserID       string                            `json:"userId,omitempty"`
		AppVersion   string                            `json:"appVersion,omitempty"`
		OSVersion    string                            `json:"osVersion,omitempty"`
		ReleaseStage string                            `json:"releaseStage"`
		Context      string                            `json:"context,omitempty"`
		Exceptions   []bugsnagException                `json:"exceptions"`
		MetaData     map[string]map[string]interface{} `json:"metaData,omitempty"`
	}
)

func send(events []*bugsnagEvent) error {
	if APIKey == "" {
		return errors.New("Missing APIKey")
	}
	payload := &bugsnagPayload{
		Notifier: Notifier,
		APIKey:   APIKey,
		Events:   events,
	}
	protocol := "http"
	if UseSSL {
		protocol = "https"
	}
	if b, err := json.MarshalIndent(payload, "", "\t"); err != nil {
		return err
	} else if resp, err := http.Post(protocol+"://notify.bugsnag.com", "application/json", bytes.NewBuffer(b)); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	} else if Verbose {
		println(string(b))
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		println(resp.StatusCode)
		println(resp.Status)
		println(string(b))
	}
	return nil
}

func getStacktrace() []bugsnagStacktrace {
	var stacktrace []bugsnagStacktrace
	i := 3 // First 3 lines are our own functions, not interesting
	for {
		if pc, file, line, ok := runtime.Caller(i); !ok {
			break
		} else {
			methodName := "unnamed"
			if f := runtime.FuncForPC(pc); f != nil {
				methodName = f.Name()
			}
			traceLine := bugsnagStacktrace{
				File:       file,
				LineNumber: strconv.Itoa(line),
				Method:     methodName,
			}
			stacktrace = append(stacktrace, traceLine)
		}
		i += 1
	}
	return stacktrace
}

// Notify sends an error to bugsnag
func Notify(err error) error {
	return New(err).Notify()
}

// NotifyRequest sends an error to bugsnag, and sets request
// URL as the event context.
func NotifyRequest(err error, r *http.Request) error {
	return New(err).WithContext(r.URL.String()).Notify()
}

// New returns returns a bugsnag event instance, that can be further configured
// before sending it.
func New(err error) *bugsnagEvent {
	return &bugsnagEvent{
		AppVersion:   AppVersion,
		OSVersion:    OSVersion,
		ReleaseStage: ReleaseStage,
		Exceptions: []bugsnagException{
			bugsnagException{
				ErrorClass: reflect.TypeOf(err).String(),
				Message:    err.Error(),
				Stacktrace: getStacktrace(),
			},
		},
	}
}

// WithUserID sets the user_id property on the bugsnag event.
func (event *bugsnagEvent) WithUserID(userID string) *bugsnagEvent {
	event.UserID = userID
	return event
}

func (event *bugsnagEvent) WithContext(context string) *bugsnagEvent {
	event.Context = context
	return event
}

// WithMetaDataValues sets bunch of key-value pairs under a tab in bugsnag
func (event *bugsnagEvent) WithMetaDataValues(tab string, values map[string]interface{}) *bugsnagEvent {
	if event.MetaData == nil {
		event.MetaData = make(map[string]map[string]interface{})
	}
	event.MetaData[tab] = values
	return event
}

// WithMetaData adds a key-value pair under a tab in bugsnag
func (event *bugsnagEvent) WithMetaData(tab string, name string, value interface{}) *bugsnagEvent {
	if event.MetaData == nil {
		event.MetaData = make(map[string]map[string]interface{})
	}
	if event.MetaData[tab] == nil {
		event.MetaData[tab] = make(map[string]interface{})
	}
	event.MetaData[tab][name] = value
	return event
}

// Notify sends the configured event off to bugsnag.
func (event *bugsnagEvent) Notify() error {
	return send([]*bugsnagEvent{event})
}
