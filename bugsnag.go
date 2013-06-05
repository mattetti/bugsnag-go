// Package bugsnag is a client for bugsnag.com, written in Go.
// API type descriptions are taken from https://bugsnag.com/docs/notifier-api (© Bugsnag Inc. 2013)
package bugsnag

// BugsnagClient describes the notifier itself. These properties are used 
// within Bugsnag to track error rates from a notifier.
type BugsnagClient struct {
        // The notifier name
    Name string `json:"name"`
        // The notifier's current version
    Version string `json:"version"`
        // The URL associated with the notifier
    URL string `json:"url"`
}


// BugsnagPayload is the actual data that is serialized to JSON
// and sent to bugsnag.com
type BugsnagPayload struct {
    // The API Key associated with the project. Informs Bugsnag which project 
    // has generated this error.
    APIKey string `json:"apiKey"` // c9d60ae4c7e70c4b6c4ebd3e8056d2b8

    Notifier *BugsnagClient `json:"notifier"`

    // An array of error events that Bugsnag should be notified of. A notifier
    // can choose to group notices into an array to minimize network traffic, or
    // can notify Bugsnag each time an event occurs. 
    Events []BugsnagEvent `json:"events"`
}

// BugsnagException is data about the error that happened.
// One payload can contain many Exceptions.
type BugsnagException struct {
    // The class of error that occurred. This field is used to group the
    // errors together so should not contain any contextual information
    // that would prevent correct grouping. This would ordinarily be the
    // Exception name when dealing with an exception.
    ErrorClass string `json:"errorClass"` // "NoMethodError"

    // The error message associated with the error. Usually this will 
    // contain some information about this specific instance of the error
    // and is not used to group the errors (optional, default none).
    Message string `json:"message,omitempty"` // "Unable to connect to database."

    // An array of stacktrace objects. Each object represents one line in
    // the exception's stacktrace. Bugsnag uses this information to help
    // with error grouping, as well as displaying it to the user.
    Stacktrace []BugsnagStacktrace `json:"stacktrace"`
}

// BugsnagStacktrace represents one line in
// the exception's stacktrace. Bugsnag uses this information to help
// with error grouping, as well as displaying it to the user.
type BugsnagStacktrace struct {
    // The file that this stack frame was executing.
    // It is recommended that you strip any unnecessary or common
    // information from the beginning of the path.
    File string `json:"file"` // "controllers/auth/session_controller.rb"

    // The line of the file that this frame of the stack was in.
    LineNumber string `json:"lineNumber"` // 1234

    // The method that this particular stack frame is within.
    Method string `json:"method"` // "create"

    // Is this stacktrace line is in the user's project code, set 
    // this to true. It is useful for developers to be able to see 
    // which lines of a stacktrace are within their own application, 
    // and which are within third party libraries. This boolean field
    // allows Bugsnag to display this information in the stacktrace
    // as well as use the information to help group errors better.
    // (Optional, defaults to false).
    InProject bool `json:"inProject,omitempty"` // true
}

// BugsnagEvent that Bugsnag should be notified of.
type BugsnagEvent struct {
    // A unique identifier for a user affected by this event. This could be 
    // any distinct identifier that makes sense for your application/platform.
    // This field is optional but highly recommended.
    UserID string `json:"userId,omitempty"`

    // The version number of the application which generated the error.
    // (optional, default none)
    AppVersion string `json:"appVersion,omitempty"`

    // The operating system version of the client that the error was 
    // generated on. (optional, default none)
    OSVersion string `json:"osVersion,omitempty"`

    // The release stage that this error occurred in, for example 
    // "development" or "production". This can be any string, but "production"
    // will be highlighted differently in bugsnag in the future, so please use
    // "production" appropriately.
    ReleaseStage string `json:"releaseStage"`

    // A string representing what was happening in the application at the 
    // time of the error. This string could be used for grouping purposes, 
    // depending on the event.
    // Usually this would represent the controller and action in a server 
    // based project. It could represent the screen that the user was 
    // interacting with in a client side project.
    // For example,
    //   * On Ruby on Rails the context could be controller#action
    //   * In Android, the context could be the top most Activity.
    //   * In iOS, the context could be the name of the top most UIViewController
    Context string `json:"context"`

    // An array of exceptions that occurred during this event. Most of the
    // time there will only be one exception, but some languages support 
    // "nested" or "caused by" exceptions. In this case, exceptions should 
    // be unwrapped and added to the array one at a time. The first exception
    // raised should be first in this array.
    Exceptions []BugsnagException `json:"exceptions"`
    // An object containing any further data you wish to attach to this error
    // event. This should contain one or more objects, with each object being
    // displayed in its own tab on the event details on the Bugsnag website.
    // (Optional).
    MetaData map[string]BugsnagMetaData `json:"metaData,omitempty"`
}

// BugsnagMetaData contains any further data you wish to attach to this error
// event. This should contain one or more objects, with each object being
// displayed in its own tab on the event details on the Bugsnag website.
type BugsnagMetaData struct {
    // A key value pair that will be displayed in the first tab
    Key string  `json:"key"`
    // This is shown as a section within the first tab
    SetOfKeys map[string]string `json:"setOfKeys"`
}
