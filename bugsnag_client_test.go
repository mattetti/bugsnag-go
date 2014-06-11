package bugsnag

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNotify(t *testing.T) {
	// Configure bugsnag
	Verbose = true
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	AppVersion = "1.2.3"
	OSVersion = "3.2.1"

	// Notify about an error
	e := errors.New("This is a test")
	if err := Notify(e); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyRequest(t *testing.T) {
	// Configure bugsnag
	Verbose = true
	APIKey = os.Getenv("BUGSNAG_APIKEY")

	// Notify about an error
	e := errors.New("This is a test")
	if r, err := http.NewRequest("GET", "some URL", nil); err != nil {
		t.Fatal(err)
	} else if err := NotifyRequest(e, r); err != nil {
		t.Fatal(err)
	}
}

func TestSetMetaDataBeforeNotify(t *testing.T) {
	// Configure bugsnag
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("This is another test")
	values := map[string]interface{}{
		"account_id": 5555,
		"user_agent": "ie4",
	}
	if err := New(e).WithUserID("12345").WithMetaDataValues("user_info", values).Notify(); err != nil {
		t.Fatal(err)
	}
}

func TestWithMetaDataBeforeNotify(t *testing.T) {
	// Configure bugsnag
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("Another error")
	if err := New(e).WithUserID("12345").WithMetaData("user_info", "name", "mr. Fu").Notify(); err != nil {
		t.Fatal(err)
	}
}

func TestNewNotify(t *testing.T) {
	// Configure bugsnag
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("One more error")
	if err := New(e).Notify(); err != nil {
		t.Fatal(err)
	}
}

func TestCapturePanic(t *testing.T) {
	// Configure bugsnag
	Verbose = true
	APIKey = os.Getenv("BUGSNAG_APIKEY")

	// Notify about an error
	r, err := http.NewRequest("GET", "some URL", nil)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			if rec.(string) != "This should be reported!" {
				t.Fatal("Unexpected capture: " + rec.(string))
			}
		} else {
			t.Fatal("Panic not captured")
		}
	}()

	defer CapturePanic(r)

	panic("This should be reported!")
}

func TestStacktraceFunc(t *testing.T) {
	TraceFilterFunc = func(traces []bugsnagStacktrace) []bugsnagStacktrace {
		for i := 0; i < len(traces); i++ {
			traces[i].File = "ceci n'est pas un string" + traces[i].File
		}
		return traces
	}
	defer func() {
		TraceFilterFunc = nil
	}()

	traces := getStacktrace()
	for _, trace := range traces {
		if !strings.HasPrefix(trace.File, "ceci n'est pas un string") {
			t.Fatal("TraceFilterFunc is not correctly filtering the stacktrace")
		}
	}

}
