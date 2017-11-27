package main;

import (
//	"fmt"
	"net/http"
	"log"
	"github.com/husobee/vestigo"
)

func PermalinkCB (w http.ResponseWriter, r *http.Request) {
	pledge()
	id := vestigo.Param(r, "id")

	post, err := Get(ctx.dbh, id)
	if (err != nil) {
		log.Fatal(err)
	}

	w.Write(([]byte)(post.String()))
}
