package main

import (
	"database/sql"
	"log"
	"fmt"
	"time"
//	"json"
)

type post struct {
	id string
	author string
	post string
	favorites int
	replyto string
	datetime time.Time
}

func Put(dbh *sql.DB, p *post) (string, error) {
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

func Get(dbh *sql.DB, id string) (*post, error) {
	fmt.Println("Grabbing post with id " + id)
	sel, err := dbh.Query("select * from posts where id = ?", id)

	if (err != nil) {
		return nil, NewUpError("Unable to retrieve post")
	}

	defer sel.Close();
	
	var ret post

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

func String(p *post) string {
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
