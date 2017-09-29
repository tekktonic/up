package main;

import (
	"net/http"
)

var getCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	"/post/:id" : PermalinkCB,
//	"/timeline/:max" : TimelineCB,
}

var postCallbacks = map[string]http.Handler {}
