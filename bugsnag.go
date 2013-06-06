package bugsnag

import (
	"bytes"
	"encoding/json"
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
	ReleaseStage        = "development"
	NotifyReleaseStages = []string{"production"}
	AutoNotify          = true
	UseSSL              = false
	Verbose             = false
	notifier            = &bugsnagNotifier{
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
		Events   []bugsnagEvent   `json:"events"`
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
		UserID       string                     `json:"userId,omitempty"`
		AppVersion   string                     `json:"appVersion,omitempty"`
		OSVersion    string                     `json:"osVersion,omitempty"`
		ReleaseStage string                     `json:"releaseStage"`
		Context      string                     `json:"context,omitempty"`
		Exceptions   []bugsnagException         `json:"exceptions"`
		MetaData     map[string]bugsnagMetaData `json:"metaData,omitempty"`
	}
	bugsnagMetaData struct {
		Key       string            `json:"key"`
		SetOfKeys map[string]string `json:"setOfKeys"`
	}
)

func newBugsnagEvent(err error, userId string) bugsnagEvent {
	exception := bugsnagException{
		ErrorClass: reflect.TypeOf(err).String(),
		Message:    err.Error(),
		Stacktrace: getStacktrace(),
	}
	exceptions := []bugsnagException{exception}
	return bugsnagEvent{
		UserID:       userId,
		AppVersion:   AppVersion,
		OSVersion:    OSVersion,
		ReleaseStage: ReleaseStage,
		Exceptions:   exceptions,
	}
}

// Notify sends an error to bugsnag.com
func Notify(err error, userId string) error {
	event := newBugsnagEvent(err, userId)
	return send([]bugsnagEvent{event})
}

func send(events []bugsnagEvent) error {
	payload := &bugsnagPayload{
		Notifier: notifier,
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
