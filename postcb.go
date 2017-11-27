package main;

import (
	"fmt"
	"strings"
	"net/http"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
//	"hash"
	"log"

)

func DistributePost(p Post) {
	fmt.Println("Distributing post")
	dbh := ctx.dbh;

	rows, err := dbh.Query("select callback,challenge from subscribers");
	defer rows.Close()

	if (err != nil) {
		log.Fatal(err)
	}

	feed := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:thr="http://purl.org/syndication/thread/1.0"  xmlns:media="http://purl.org/syndication/atommedia" xmlns:poco="http://portablecontacts.net/spec/1.0" xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:statusnet="http://status.net/schema/api/1/" xml:lang="en-US">
  <generator uri="https://tekk.in/up" version="0.1.0">Up</generator>
  <id>https://` + config.Server + `/timeline.atom</id>
  <title>` + config.Owner + ` timeline</title>
  <subtitle>Updates from ` + config.Owner + ` on ` + config.Server + `</subtitle>
  <updated>` + p.datetime.Format("2006-01-02T15:04:05-07:00") + `</updated>
  <author>
    <uri>https://` + config.Server + `/</uri>
    <name>` + config.Owner + `</name>
    <link rel="alternate" type="text/html" href="https://` + config.Server + `/" />
    <poco:preferredUsername>` + config.Owner + `</poco:preferredUsername>
    <poco:displayName>` + config.Owner + `</poco:displayName>
    <statusnet:profile_info local_id="1"/>
  </author>

  <link href="https://` + config.Server + `/push/hub" rel="hub"/>` + MakeEntry(p) +`

</feed>
`
	for rows.Next() {
		var subscriber string
		var challenge string
		if err := rows.Scan(&subscriber, &challenge); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Distributing to " + subscriber)
		mac := hmac.New(sha1.New, ([]byte)(challenge))
		mac.Write((([]byte)(feed)))
		signature := mac.Sum(nil)
		client := http.Client{}
		req, err := http.NewRequest("POST", subscriber, strings.NewReader(feed))
		if (err != nil) {
			log.Println(err)
		}
		req.Header.Add("Content-Type", "application/xrd+xml")
		req.Header.Add("X-Hub-Signature", "sha1=" + hex.EncodeToString(signature))
		_, err = client.Do(req);

		if (err != nil) {
			log.Fatal(err)
		}
	}
}
func PostCB (w http.ResponseWriter, r *http.Request) {
	pledge()

	str := auth(r.Header.Get("X-Up-Auth"))
	// Make sure we're allowed to look at our timeline
	if (str == "") {
		r.ParseForm()
		
		post := r.PostForm["text"];

		fmt.Println(post);
		// Make sure the post is of an allowed length
		if (len(post)> config.Max || len(post) == 0) {
			w.WriteHeader(413);
			w.Write([]byte("Payload invalid"))
			return;
		}
		
		entry := NewPost(strings.Join(post, " "))

		id, err := Put(ctx.dbh, &entry)

		if (err != nil) {
			log.Fatal(err)
		}

		go DistributePost(entry);
		w.Write([]byte("http://" + config.Server + "/post/" + id));
		return
	}

	w.WriteHeader(401)

	w.Write(([]byte)(""))

}
