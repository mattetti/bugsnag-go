/*

Package bugsnag provides a go interface to the Bugsnag API.

For more information about the API see https://bugsnag.com/docs/notifier-api

Usage:

	package main

	import (
		"errors"
		"github.com/mattetti/bugsnag"
	)

	func main() {
		// First, configure bugsnag. Only APIKey is mandatory, other settings are optional.
		bugsnag.APIKey = "c9d60ae4c7e70c4b6c4ebd3e8056d2b8"
		bugsnag.AppVersion = "1.0.2"
		bugsnag.OSVersion = "Windows XP"
		bugsnag.ReleaseStage = "production"
		bugsnag.NotifyReleaseStages = []string{"production"}
		bugsnag.UseSSL = true
		bugsnag.Verbose = true

		// Then, out of the blue, an error happens:
		err := errors.New("Something bad just happened")
		bugsnag.Notify(err)

		// To notify about a HTTP handler error (assuming r is a *http.Request):
		// bugsnag.NotifyRequest(err, r)

		// In case you need to supply the user ID:
		bugsnag.New(err).WithUserID("12345").Notify()

		// To report what your app was doing while error happened:
		bugsnag.New(err).WithContext("something").Notify()

		// Metadata can be set all at once:
		values := map[string]interface{}{
			"user_agent": "ie4",
			"account_id": 5555,
		}
		bugsnag.New(err).WithMetaDataValues("tab_name_in_bugsnag", values).Notify()

		// Or one key-value pair at a time:
		bugsnag.New(err).WithMetaData("tab_name_in_bugsnag", "user_agent", "ie4").WithMetaData("tab_name_in_bugsnag", "account_id", 5555).Notify()

		// To capture a HTTP handler panic, add this to your handler (assuming r is a *http.Request):
		// defer bugsnag.CapturePanic(r)
	}

The bugsnag event instance sends a stacktrace to the API in the metadata. This stacktrace can be filtered to remove application-specific noise:

	bugsnag.TraceFilterFunc = func(traces []bugsnagStacktrace) []bugsnagStacktrace {
		// modify the stacktrace here
		return traces
	}

	bugsnag.New(err).Notify()
*/

package bugsnag
