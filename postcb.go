package main;

import (
	"fmt"
	"strings"
	"net/http"
	"log"
)

func PostCB (w http.ResponseWriter, r *http.Request) {


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

		w.Write([]byte("http://" + config.Server + "/post/" + id));
		return
	}

	w.WriteHeader(401)

	w.Write(([]byte)(""))
}
