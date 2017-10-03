package main;

import (
	"net/http"

)

var apiheader string = "/api/v1"

var getCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	"/post/:id" : PermalinkCB,
	apiheader + "/timeline/:max" : TimelineCB,
	apiheader + "/timeline/" : TimelineCB,
}

var postCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	apiheader + "/post/" : PostCB,
}
