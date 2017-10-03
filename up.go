package main;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"net/http"
	"github.com/husobee/vestigo"
	"fmt"
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
	
	post := Post{id: "", author: "tekk@up.tekk.in", post: "all is bad", favorites: 0, replyto: ""}

	fmt.Println(post);
	
	id, _ := Put(dbh, &post);
	fmt.Println("trying to retrieve");
	gotten,_ := Get(dbh,id);
	fmt.Print(gotten)

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
