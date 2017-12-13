package main

import (
	"container/list"
	"database/sql"
	"html"
	"log"
	"fmt"
	"strconv"
	"strings"
	"time"
	"encoding/json"
)

type Post struct {
	id string
	author string
	post string
	favorites int
	replyto string
	datetime time.Time
	foreign bool
	tags []string
}

type jsonpost struct {
	Post string `json:"post"`
	Replyto string `json:"replyto"`
}


func NewPost(text string) Post {
	text = strings.Replace(text, "\n", " ", -1)
	
	ret := Post{author: config.Owner + "@" + config.Server,
		favorites: 0,
		post: text,
		datetime: time.Now(),
		replyto: "",
		foreign: false};


	for _,w := range (strings.Split(text, " ")) {
		// This is safe because the strings are UTF8
		if w[0] == '#' {
			ret.tags = append(ret.tags, w);
		}
	}
	return ret;
}
func Put(dbh *sql.DB, p *Post) (string, error) {
	// 1 in 144 quadrillion is *probably* safe
	p.id = keyGen(10)

	fmt.Println("Current id is " + p.id)
	post, err := Get(dbh, p.id);

	fmt.Println("Escaped our check");
	// Generate a key we know is collision free just in case
	for (post != nil) {
		p.id = keyGen(10)
		post, err = Get(dbh, p.id);
	}
	
	tx, err := dbh.Begin();
	if (err != nil) {
		log.Fatal(err);
	}
	
	insert, err := tx.Prepare("insert into posts(id, author, post, favorites, replyto, datetime) values(?,?,?,?,?,?)")

	if (err != nil) {
		log.Fatal(err);
	}


	log.Println("Inserting new post " + p.id + " into database at " + strconv.FormatInt(p.datetime.Unix(),10));
	_, err = insert.Exec(p.id, p.author, p.post, p.favorites, p.replyto, p.datetime.Unix())

	if (err != nil) {
		log.Fatal(err);
	}

	err = tx.Commit()

	
	if (err != nil) {
		log.Fatal(err);
	}


	insert.Close();

	tx, err = dbh.Begin();
	if (err != nil) {
		log.Fatal(err);
	}


	for _,tag := range p.tags {
		insert, err = tx.Prepare("insert into tags(postid, tagname) values(?,?)")

		if (err != nil) {
			log.Fatal(err);
		}

		defer insert.Close()


		_, err = insert.Exec(p.id, tag);

		if (err != nil) {
			log.Fatal(err);
		}

		err = tx.Commit()
		
		if (err != nil) {
			log.Fatal(err);
		}
	}
	
	return p.id, nil;
}

func Get(dbh *sql.DB, id string) (*Post, error) {
	fmt.Println("Grabbing post with id " + id)
	sel, err := dbh.Query("select * from posts where id = ?", id)

	if (err != nil) {
		log.Println("First")
		log.Fatal(err)
		return nil, NewUpError("Unable to retrieve post")
	}

	var ret Post

	// Safe, selected on primary key
	sel.Next()
	var stringtime int64
	err = sel.Scan(&ret.id, &ret.author, &ret.post,
		&ret.favorites, &ret.replyto, &stringtime)

	if (err != nil) {
		return nil, NewUpError("Unable to retrieve post")
	}
	fmt.Println("Stringtime is " + strconv.FormatInt(stringtime,10))
	ret.datetime = time.Unix(stringtime, 0)
	fmt.Println("Datetime is " + strconv.FormatInt(ret.datetime.Unix(),10));
	sel.Close()

	
	sel, err = dbh.Query("select tagname from tags where postid = ?", id)

	if (err != nil) {
		log.Println("Tag retrieval initial")
		log.Fatal(err)
		return nil, NewUpError("Unable to retrieve post")
	}

	defer sel.Close();
	

	// Safe, selected on primary key
	for sel.Next() {
		var tag string
		err = sel.Scan(&tag)

		if (err != nil) {
			log.Println("Something went wrong in tag application")
			log.Fatal(err);
			return nil, NewUpError("Unable to retrieve post")
		}

		ret.tags = append(ret.tags, tag);
	}

	return &ret, nil
}

func Timeline(dbh *sql.DB, max int) *list.List {
	ret := list.New()

	// Safe, max has been verified as sane at this point.
	sel, err := dbh.Query("select * from posts order by datetime desc limit " + strconv.Itoa(max))

	if (err != nil) {
		log.Fatal("Something went horribly wrong in retrieving the timeline")
	}


	for sel.Next() {
		var item Post
		var stringtime int64
		err = sel.Scan(&item.id, &item.author, &item.post,
			&item.favorites, &item.replyto, &stringtime)

		if (err != nil) {
			log.Fatal("Database somehow corrupted")
		}

		

		item.datetime = time.Unix(stringtime,0)

		// Have things in reverse order so that when we print them the most recent is at the bottom
		ret.PushFront(item)
	}

	sel.Close()

	// Build the tags for our selection
	for e := ret.Front(); e != nil; e = e.Next() {
		post := e.Value.(Post)
		
		sel, err = dbh.Query("select tagname from tags where postid = ?", post.id)

		if (err != nil) {
			return nil
		}

		defer sel.Close();
		

		// Safe, selected on primary key
		for sel.Next() {
			var tag string
			err = sel.Scan(&tag)

			if (err != nil) {
				return nil
			}

			post.tags = append(post.tags, tag);
		}
	}
	return ret
}

func (p Post) String() string {
	var ret string
	var part2ret string
	partret := p.author + "\r\n" + p.datetime.Format(time.UnixDate) + "\r\n";
	
	if (p.replyto == "") {
		part2ret = partret + p.replyto + "\r\n"
	} else {
		part2ret = partret + "[in reply to nobody]"
	}

	ret = part2ret + "\r\n" + p.post + "\r\n\r\n"

	return ret
}

func FromJSON(in []byte) Post {
	jsonresult := jsonpost{}

	err := json.Unmarshal(in, &jsonresult)
	if ((err != nil)) {
		log.Fatal(err)
	}

	return Post{}
	
}

func MakeEntry(p Post) string {
	return `<entry  xmlns="http://www.w3.org/2005/Atom" xmlns:thr="http://purl.org/syndication/thread/1.0"  xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:poco="http://portablecontacts.net/spec/1.0" xmlns:statusnet="http://status.net/schema/api/1/">
    <id>` + config.Server + "/" + p.id + `</id>
    <title/>
    <content type="html">` + html.EscapeString(p.post) + `</content>
    <published>` + p.datetime.Format("2006-01-02T15:04:05-07:00") + `</published>
    <updated>` + p.datetime.Format("2006-01-02T15:04:05-07:00") + `</updated>
    <statusnet:notice_info local_id="` + p.id + `" source="activity"/>
  </entry>`
}

type ByTime []Post

func (t ByTime) Len() int {return len(t)}
func (t ByTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTime) Less(i, j int) bool { return t[i].datetime.Before(t[j].datetime) }
