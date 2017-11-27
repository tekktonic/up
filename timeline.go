package main;

import (
//	"fmt"

	"net/http"
//	"log"
	"strconv"
	"github.com/husobee/vestigo"
)

func TimelineCB (w http.ResponseWriter, r *http.Request) {
	pledge()
	max, err := strconv.Atoi(vestigo.Param(r, "max"))

	// If we need that big a max, something's probably misbehaving. I may change later.
	if (err != nil || max > 1024) {
		max = config.TimelineSize
	}

	str := auth(r.Header.Get("X-Up-Auth"))
	// Make sure we're allowed to look at our timeline
	if (str == "") {
		post := Timeline(ctx.dbh, max)
		for e := post.Front(); e != nil; e = e.Next(){
			w.Write(([]byte)((e.Value.(Post)).String()))
		}

		return
	}

	w.WriteHeader(401)

	w.Write(([]byte)(str))
}
