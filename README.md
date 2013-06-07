About
=====

[bugsnag](http://bugsnag.com) client in Go

Please see https://bugsnag.com/docs/notifier-api for more info about the bugsnag API.

How to use
==========

```go
package main

import (
    "errors"
	"github.com/toggl/bugsnag"
)

func main() {
	// First, configure bugsnag
	bugsnag.APIKey = "c9d60ae4c7e70c4b6c4ebd3e8056d2b8"
	bugsnag.Verbose = true

	// Then, out of the blue, an error happens:
	err := errors.New("Something bad just happened")
	bugsnag.Notify(err)
	
	// To notify about a HTTP handler error (assuming r is a *http.Request):
	// bugsnag.NotifyRequest(err, r)

	// In case you need to supply the user ID:
	bugsnag.New(err).SetUserID("12345").Notify()

	// To report what your app was doing while error happened:
	bugsnag.New(err).SetContext("something").Notify()

	// Metadata can be set all at once:
	values := map[string]interface{}{
		"user_agent": "ie4",
		"account_id": 5555,
	}
	bugsnag.New(err).SetMetaData("tab_name_in_bugsnag", values).Notify()

	// Or one key-value pair at a time:
	bugsnag.New(err).AddMetaData("tab_name_in_bugsnag", "user_agent", "ie4").AddMetaData("tab_name_in_bugsnag", "account_id", 5555).Notify()

	// AddMetaData(), SetMetaData() and SetUser() can be chained together (as above)
	// or left out alltogether:
	bugsnag.New(err).Notify()
}
```
