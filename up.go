package main;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"net/http"
	"github.com/husobee/vestigo"
	"time"
	"math/rand"
	"log"
	"golang.org/x/sys/unix"
)


func pledge() {
	unix.Pledge("stdio inet", []string{})
}
var ctx context

func main() {
	configfile := "config.json"
	if (len(os.Args) > 1) {
		configfile = os.Args[1]
	}
	readConfig(configfile)
	
	dbh, err := sql.Open("sqlite3", config.DbFile);


	if (err != nil) {
		log.Fatal(err)
	}
	// Our context needs the handle so that it can get pushed around in our 
	ctx.dbh = dbh

	rand.Seed(time.Now().UnixNano())
	router := vestigo.NewRouter()

	vestigo.CustomNotFoundHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Failed to find page " + r.URL.String())
		log.Println(r)
		r.ParseForm();
		log.Println("Headers")
		log.Println(r.Header)
	});
	// Set up our URL handlers. See <<mapping.go>>
	for path, cb := range getCallbacks {
		router.Get(path, cb)
	}

	// Set up our URL handlers. See <<mapping.go>>
	for path, cb := range postCallbacks {
		router.Post(path, cb)
	}

	unix.Chroot(".")
	log.Fatal(http.ListenAndServe(":" + config.Port, router))
}
