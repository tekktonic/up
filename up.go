package main;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"net/http"
	"github.com/husobee/vestigo"

	"log"
)


var ctx context

func main() {

	// Dummy right now, soon™
	readConfig()
	
	dbh, err := sql.Open("sqlite3", "up.db");

	if (err != nil) {
		log.Fatal(err)
	}
	// Our context needs the handle so that it can get pushed around in our 
	ctx.dbh = dbh
	
	router := vestigo.NewRouter()
	
	// Set up our URL handlers. See <<mapping.go>>
	for path, cb := range getCallbacks {
		router.Get(path, cb)
	}

	// Set up our URL handlers. See <<mapping.go>>
	for path, cb := range postCallbacks {
		router.Post(path, cb)
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}
