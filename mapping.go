package main;

import (
	"net/http"

)

var apiheader string = "/api/v1"

var getCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	"/post/:id" : PermalinkCB,
	"/push/callback/:id" : HubResponseCB,
	apiheader + "/timeline/:max" : TimelineCB,
	apiheader + "/timeline/" : TimelineCB,
	apiheader + "/debug/webfinger/:remote" : webfingerCB,
	"/.well-known/host-meta" : HostMetaCB,
	"/.well-known/webfinger" : WebfingerCB,
	"/atom/feed.atom" : GenTimeline,

//	"/push/callback/:id" : WebfingerCB,
}

var postCallbacks = map[string]func(http.ResponseWriter, *http.Request) {
	apiheader + "/post/" : PostCB,
	apiheader + "/subscribe/" : SubscribeCB,
	"/push/hub" : HubCB,
	"/push/callback/:id" : RemotePostCB,
}
