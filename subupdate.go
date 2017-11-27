package main;

import (
	"time"
	"strconv"
	"log"
)

func SubUpdate() {
	dbh := ctx.dbh
	rows, err := dbh.Query("select id,topic,lifetime,start from subscriptions")
	if (err != nil) {
		log.Fatal(err)
	}
	for rows.Next() {
		var id string
		var topic string
		var lifetime int
		var start int

		err := rows.Scan(&id, &topic, &lifetime, &start)

		if (err != nil) {
			log.Fatal(err)
		}

		// We're storing our time as unix timestamps
		// so lifetime + start is valid and correct
		expiry,_ := time.Parse(time.UnixDate, strconv.Itoa(lifetime + start))
		if (time.Until(expiry) < (time.Hour * 24)) {
//			httpSubscribe()
		}
	}
}
