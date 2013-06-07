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
  userID := "12345"
  bugsnag.Notify(err, userID)
}
```
