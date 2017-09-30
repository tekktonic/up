package main;

import (
	"net/http"

)

var getCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	"/post/:id" : PermalinkCB,
	"/timeline/:max" : TimelineCB,
	"/timeline/x" : TimelineCB,
}

var postCallbacks = map[string]http.Handler {}
