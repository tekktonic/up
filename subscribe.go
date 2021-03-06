package main;

import (
	"math/rand"
	"strconv"
	"fmt"
	"log"
	"time"
	"net/http"
	"net/url"
//	"io/ioutil"
//	"database/sql"
	"github.com/antchfx/xquery/xml"
	"github.com/husobee/vestigo"
	"github.com/kr/pretty"
)


type Subscription struct {
	id int
	topic string
	name string
	challenge string
	hub string
	pending int
	start int64
	time int64
}

func genKey() string {
	// Entropy of 26**n...
	const n int = 10
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var ret [n]byte

	for i := 0; i < n; i++ {
		// Generate a random letter in a-z
		ret[i] = (byte)(r.Intn(26) + 97)
	}
	realret := string(ret[:])
	
	return realret
}

func saveSubscription(hub string, finger FingerInfo) Subscription {
	log.Println("Made it into saveSubscription")
	dbh := ctx.dbh;

	id := rand.Int()
	now := time.Now().Unix()
	ret := Subscription{
		id: id,
		topic: finger.atomfeed,
		name: finger.id,
		challenge: genKey(),
		pending: 0,
		start: now,
		time: 0,
	}

	pretty.Println(ret)
	tx, err := dbh.Begin();
	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Println("Preparing to insert...")
	insert, err := tx.Prepare("insert into subscriptions(topic, id, name, challenge, pending, lifetime, start, hub) values(?,?,?,?,?,?,?,?);")

	if (err != nil) {
		log.Fatal(err)
	}

	defer insert.Close()
	fmt.Println("Executing insert...")
	_, err = insert.Exec(ret.topic, ret.id, ret.name, ret.challenge, ret.pending, ret.time, ret.start, "")

	if (err != nil) {
		log.Fatal(err)
	}

	err = tx.Commit()

	if (err != nil) {
		log.Fatal(err)
	}

	fmt.Println("Got to end");
	return ret;
}

/*
  Subscription failed somewhere early in the process, just roll it back to keep the db clean
*/
func rollbackSubscribe(id int) {
	dbh := ctx.dbh
	tx, err := dbh.Begin()
	if (err != nil) {
		log.Fatal(err);
	}

	delete, err := tx.Prepare("delete from subscriptions where id=?;");

	defer delete.Close()

	_, err = delete.Exec(id);

	if (err != nil) {
		log.Fatal(err)
	}

	err = tx.Commit()

	if (err != nil) {
		log.Fatal(err)
	}

	fmt.Println("Made it out of saveSubscription")
}

func httpSubscribe(hub string, sub Subscription) {
	fmt.Println("httpSubscribe");
	_, err := http.PostForm(hub, url.Values{
		"hub.callback": {"https://" + config.Server + "/push/callback/" + strconv.Itoa(sub.id)},
		"hub.topic": {sub.topic},
		"hub.mode": {"subscribe"},
		
	})

	if (err != nil) {
		rollbackSubscribe(sub.id)
		log.Fatal(err);
	}
}

func SubscribeCB(w http.ResponseWriter, r *http.Request) {
	pledge()
	fmt.Println("SubscribeCB");
	str := auth(r.Header.Get("X-Up-Auth"))
	if (str == "") {
		r.ParseForm()
		user := (r.PostForm["user"])[0]
		fmt.Println("Trying to subscribe to " + user)
		finger := getFingerInfo(user)

		feed := getAtomFeed(finger.feed)
//		fmt.Println("User's feed is " + feed.String())
		huburl := getHub(feed)
		subscription := saveSubscription(huburl, finger)
		httpSubscribe(huburl, subscription)

		// Subscription is only half done here; we need to have the hub verify it (below).
	}
}


func HubResponseCB(w http.ResponseWriter, r *http.Request) {
	pledge()
	id,err := strconv.Atoi(vestigo.Param(r, "id"));

	fmt.Println("HubResponseCB");
	dbh := ctx.dbh
	// If the ID can't possibly be right, don't bother.
	if (err != nil) {
		w.WriteHeader(400);
		w.Write([]byte("Invalid callback URL"))
		return;
	}

	r.ParseForm()
	topic := r.Form["hub.topic"][0]
	mode := r.Form["hub.mode"][0]
	challenge := r.Form.Get("hub.secret")
	lease := r.Form.Get("hub.lease_seconds")

	
	if (mode == "subscribe") {
		sel, err := dbh.Query("select * from subscriptions where id = ?", id)

		if (err != nil) {
			w.WriteHeader(400);
			w.Write([]byte("Invalid callback URL"))
			return;			
		}

		var sub Subscription;
		// If there was no error, then there was an entry
		sel.Next()

		err = sel.Scan(&sub.topic, &sub.id, &sub.name, &sub.challenge, &sub.pending, &sub.time, &sub.start, &sub.hub)

		if (err != nil) {
			log.Fatal(err)
		}
		if (topic != sub.topic || id != sub.id/* || verify != sub.challenge*/) {
			fmt.Println(topic + " : " + sub.topic)
			fmt.Println(strconv.Itoa(id) + " : " + strconv.Itoa(sub.id))
			w.WriteHeader(404);
			w.Write([]byte("Subscription not requested."))
			return;
		}


		
		w.Write([]byte(challenge))

		sel.Close()
		
		tx, err := dbh.Begin()

		if (err != nil) {
			log.Fatal(err)
		}

		insert, err := tx.Prepare("update subscriptions set lifetime=?,start=? where id=?");
		if (err != nil) {
			log.Fatal(err)
		}

		leasetime, err := strconv.Atoi(lease)

		if (err != nil) {
			leasetime = (2 << 31) - 1
		}
		
		_,err = insert.Exec(strconv.Itoa(leasetime), time.Now().Format(time.UnixDate), sub.id);

		if (err != nil) {
			log.Fatal(err)
		}
		
		err = tx.Commit()

		if (err != nil) {
			log.Fatal(err)
		}

		insert.Close();
		return;
	} else if (mode == "unsubscribe") {
		w.WriteHeader(400);
		w.Write([]byte("Not yet implemented"))
		return
	} else
	{
		w.WriteHeader(400);
		w.Write([]byte("Invalid subscribe mode"))
		return;
	}
}

func RemotePostCB(w http.ResponseWriter, r *http.Request) {
	pledge()
	fmt.Println("Remote post")
	id,err := strconv.Atoi(vestigo.Param(r, "id"));

	dbh := ctx.dbh
	// If the ID can't possibly be right, don't bother.
	if (err != nil) {
		w.WriteHeader(400);
		w.Write([]byte("Invalid callback URL"))
		return;
	}
	fmt.Println("Callback at " + strconv.Itoa(id))
	sel, err := dbh.Query("select name from subscriptions where id = ?", id)
	fmt.Println("Did query")
	if (err != nil) {

		fmt.Println(r);
		w.WriteHeader(400);
		w.Write([]byte("Invalid callback URL"))
		return;			
	}


	if (!sel.Next()) {
		fmt.Println("Next failed")
		return;
	}
	
	fmt.Println("Got next")
	var name string;

	err = sel.Scan(&name);

	if (err != nil) {

		fmt.Println(id)
		//		log.Fatal(err)
		return
	}
	fmt.Println("Scanned")

	tree, err := xmlquery.Parse(r.Body)
	fmt.Println("Parsed body")
	text := xmlquery.FindOne(tree, "//feed/entry/content").InnerText()
	fmt.Println(tree);
/*	datetime, err := time.Parse("2006-01-02T15:04:05-07:00",
		xmlquery.FindOne(tree, "//feed/entry/published").InnerText())*/

	if (err != nil) {
		log.Fatal(err)
	}
	
	post := Post{post: text, replyto: "", foreign: true, favorites: 0,
		author: name,
		datetime: time.Now()}

	fmt.Println("Closing sel")
	sel.Close();
	_,err = Put(dbh, &post)
	if (err != nil) {
		log.Fatal(err)
	}
	
}
// Given a feed, retrieve our PuSH (now WebSub) hub.
func getHub(feed string) string {
	fmt.Println("Returning hub")
	resp, err := http.Get(feed)
	if (err != nil) {
		log.Fatal(err)
	}
	tree, _ := xmlquery.Parse(resp.Body)
	result := xmlquery.FindOne(tree, "//link[@rel=\"hub\"]")
	hub := getAttr(result, "", "href")

	
	return hub
}

