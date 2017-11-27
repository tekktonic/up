package main;

import (
	"net/http"
	"strings"
	"strconv"
	"regexp"
	"net/url"
	"sort"
	"log"
	"github.com/husobee/vestigo"
)

func buildTagSelect(tag string) string {
	re := regexp.MustCompile("'|\"|;")

	tag = string(re.ReplaceAll(([]byte)(tag), ([]byte)("")))
	return "select postid from tags where tagname = \"" + tag + "\" "
}
func TagCB(w http.ResponseWriter, r *http.Request) {
	tagstring := vestigo.Param(r, "tags");

	log.Println(vestigo.ParamNames(r));
	max, err := strconv.Atoi(vestigo.Param(r, "max"))

	if (err != nil || max > 1024) {
		max = config.TimelineSize
	}

	str := auth(r.Header.Get("X-Up-Auth"))

	if (str == "") {
		tags := strings.Split(tagstring, ",")
		len := len(tags)

		if (len < 1) {
			w.WriteHeader(400)
			return;
		}
		
		query := ""
		for i, encodedtag := range tags {
			tag,_ := url.QueryUnescape(encodedtag)

			query += buildTagSelect(tag)

			if (i+1 == len) {
				query += ";"
			} else {
				query += " intersect "
			}
		}

		rows, err := ctx.dbh.Query(query)

		if (err != nil){
			log.Println(query);
			log.Fatal(err)
		}

		var posts []Post
		for rows.Next() {
			var id string
			err := rows.Scan(&id)
			if (err != nil) {
				log.Fatal(err)
			}

			post, err := Get(ctx.dbh, id)

			if (err != nil) {
				log.Fatal(err)
			}

			posts = append(posts, *post)
		}

		sort.Sort(ByTime(posts))

		for _,post := range posts {
			w.Write((([]byte)(post.String())))
		}
	}

	w.WriteHeader(401)
	w.Write(([]byte)(str))
}
