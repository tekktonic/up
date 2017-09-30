package main

import (
	"container/list"
	"database/sql"
	"log"
	"fmt"
	"strconv"
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
}

type jsonpost struct {
	Post string `json:"post"`
	Replyto string `json:"replyto"`
}

func Put(dbh *sql.DB, p *Post) (string, error) {
	// 1 in 144 quadrillion is *probably* safe
	p.id = keyGen(10)

	fmt.Println("Current id is " + p.id)
	p.datetime = time.Now()
	fmt.Println(p.datetime.Format(time.UnixDate));
	post, err := Get(dbh, p.id);

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

	defer insert.Close()


	_, err = insert.Exec(p.id, p.author, p.post, p.favorites, p.replyto, p.datetime.Format(time.UnixDate))

	if (err != nil) {
		log.Fatal(err);
	}

	err = tx.Commit()
	
	if (err != nil) {
		log.Fatal(err);
	}

	return p.id, nil;
}

func Get(dbh *sql.DB, id string) (*Post, error) {
	fmt.Println("Grabbing post with id " + id)
	sel, err := dbh.Query("select * from posts where id = ?", id)

	if (err != nil) {
		return nil, NewUpError("Unable to retrieve post")
	}

	defer sel.Close();
	
	var ret Post

	// Safe, selected on primary key
	sel.Next()
	var stringtime string
	err = sel.Scan(&ret.id, &ret.author, &ret.post,
		&ret.favorites, &ret.replyto, &stringtime)

	if (err != nil) {
		return nil, NewUpError("Unable to retrieve post")
	}
	ret.datetime, _ = time.Parse(time.UnixDate, stringtime)

	return &ret, nil
}

func Timeline(dbh *sql.DB, max int) *list.List {
	ret := list.New()

	// Safe, max has been verified as sane at this point.
	sel, err := dbh.Query("select * from posts order by datetime desc limit " + strconv.Itoa(max))

	if (err != nil) {
		log.Fatal("Something went horribly wrong in retrieving the timeline")
	}

	defer sel.Close()

	for sel.Next() {
		var item Post
		var stringtime string
		err = sel.Scan(&item.id, &item.author, &item.post,
			&item.favorites, &item.replyto, &stringtime)

		if (err != nil) {
			log.Fatal("Database somehow corrupted")
		}

		item.datetime, _ = time.Parse(time.UnixDate, stringtime)

		// Have things in reverse order so that when we print them the most recent is at the bottom
		ret.PushFront(item)
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
		part2ret = partret
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
