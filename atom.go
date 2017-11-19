package main;

import (
	"fmt"
	"net/http"
)
func GenTimeline(w http.ResponseWriter, r *http.Request) {
	dbh := ctx.dbh

	var dt string;
	// This isn't input, so it's safe; the only way to inject anything is if the owner
	// does it in their config
	err := dbh.QueryRow("select datetime from posts where id = " + config.Owner + "@" + config.Server +
		" order by datetime desc limit 1").Scan(&dt);

	if (err != nil) {
		dt = "1979-01-01T00:00:00+00:00"
	}

	fmt.Println("Got a request for the feed, sending it back.")
	feed := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:thr="http://purl.org/syndication/thread/1.0"  xmlns:media="http://purl.org/syndication/atommedia" xmlns:poco="http://portablecontacts.net/spec/1.0" xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:statusnet="http://status.net/schema/api/1/" xml:lang="en-US">
  <generator uri="https://tekk.in/up" version="0.1.0">Up</generator>
  <id>https://` + config.Server + `/timeline.atom</id>
  <title>` + config.Owner + ` timeline</title>
  <subtitle>Updates from ` + config.Owner + ` on ` + config.Server + `</subtitle>
  <updated>` + dt + `</updated>
  <author>
    <uri>https://` + config.Server + `/</uri>
    <name>` + config.Owner + `</name>
    <link rel="alternate" type="text/html" href="https://` + config.Server + `/" />
    <poco:preferredUsername>` + config.Owner + `</poco:preferredUsername>
    <poco:displayName>` + config.Owner + `</poco:displayName>
    <statusnet:profile_info local_id="1"/>
  </author>

  <link href="https://` + config.Server + `/push/hub" rel="hub"/>

</feed>
`

	w.Write([]byte(feed));
}


