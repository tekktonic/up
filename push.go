package main;

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"strconv"
	"time"
	"github.com/kr/pretty"
)


// This is the minimal amount of information allowed for a PuSH subscription.
type Pushstruct struct {
	Callback string
	Action string
	Topic string
}

func HubCB(w http.ResponseWriter, r *http.Request) {
	pledge();
	log.Println("Got hit up on the hub")
	r.ParseForm();

	callback := r.Form["hub.callback"][0]
	topic := r.Form["hub.topic"][0]
	mode := r.Form["hub.mode"][0]
	secret := r.Form.Get("hub.secret")
	w.WriteHeader(202)
	w.Write([]byte(""))
	log.Println("callback: " + callback)

	log.Println("mode: " + mode)
	log.Println("topic: " + topic)
	go HubConfirm(callback, topic, secret);
	
}

func HubConfirm(c string, t string, s string) {
	dbh := ctx.dbh
	client := &http.Client{}
	req, _ := http.NewRequest("GET", c, nil);
	params := req.URL.Query()
	params.Add("hub.topic", t)
	params.Add("hub.mode", "subscribe")
	params.Add("hub.lease_seconds", strconv.Itoa((1 << 31) - 1))
	// q.Set("hub.topic", t)
	// q.Set("hub.mode", "subscribe")
	// q.Set("hub.lease_seconds", strconv.Itoa((1 << 31) - 1))
/*	req.Header.Add("hub_topic", t)
	req.Header.Add("hub_mode", "subscribe")
	req.Header.Add("hub_lease_seconds", strconv.Itoa((1 << 31) - 1))*/
	challenge := genKey()
	params.Add("hub.challenge",challenge);

	req.URL.RawQuery = params.Encode()
	pretty.Println(req.URL.Query())
	if (req.URL.Query().Get("hub.mode") == "") {
		log.Fatal("Query setting is broken. Somehow.")
	}
	time.Sleep(time.Second * 5);


	resp, err := client.Do(req)

	if (err != nil) {
		log.Fatal(err)
	}
	str,err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Fatal(err);
	}
	if ((string(str)) != challenge) {
		fmt.Println("Things didn't match up!")
		fmt.Println(str)
		fmt.Println(challenge)
		log.Fatal("Mismatched challenge");
	}

	log.Println("Putting " + c + " into db");
	tx, err := dbh.Begin()
	if (err != nil) {
		log.Fatal(err);
	}

	insert, err := tx.Prepare("insert into subscribers(callback,challenge) values(?,?);");

	defer insert.Close()
	if (err != nil) {
		log.Fatal(err)
	}

	_, err = insert.Exec(c, s)

	if (err != nil) {
		log.Fatal(err)
	}

	err = tx.Commit();

	if (err != nil) {
		log.Fatal(err)
	}
	
}
/*
type subunsub int



func PushStruct(action subunsub, topic string) {
	ret := Pushstruct{
		Action: action,
		Topic: topic,
		Callback: config.Server + "/push/callback",
	}
}

func PushSubscribe(PushStruct s, hub string) {
	client := http.Client{}

	request, _ := http.NewRequest("GET", string, nil)
	var action string
	switch s.Action {
	case SUBSCRIBE:
		action = "subscribe"
	case UNSUBSCRIBE:
		action = "unsubscribe"
	}
	request.Header.Add("hub.mode", s.action)
	request.Header.Add("hub.callback", s.callback)
	request.Header.Add("hub.topic", s.topic)
	request.Header.Add("From", config.Owner + "@" + config.Server)
}
*/
