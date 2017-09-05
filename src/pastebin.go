package main

import (
	"net/http"
	"database/sql"
	sqlite3 "github.com/mattn/go-sqlite3"
	"math/rand"
	"log"
//	"os"
	"time"
	"fmt"
)

var sqlh *sql.DB

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
	
	fmt.Println(realret)
	return realret
}

func GetHandler(w http.ResponseWriter, request *http.Request){
	
}
func PostHandler(w http.ResponseWriter, request *http.Request) {
	if (request.Method != "POST") {
		fmt.Fprintln(w, "WRONG METHOD");
		return;
	}
	err := sqlh.Ping()
	if (err != nil) {
		log.Fatal(err)
	}	

	query := request.URL.Query();

	if (sqlh == nil) { fmt.Println("EVERYTHING IS RUINED"); return}
	tx, _ := sqlh.Begin()

	insert, _ := tx.Prepare("insert into posts(id, post) values(?,?)")
	defer insert.Close()

	insert.Exec(genKey(), query["q"]);

	tx.Commit()
}

func main() {
	sqlh, err := sqlite3.Open("pastebin.db")

	if (err != nil) {
		log.Fatal(err)
	}
	err = sqlh.Ping()
	if (err != nil) {
		log.Fatal(err)
	}	
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/get", GetHandler)
	http.ListenAndServe(":8000", nil)

	sqlh.Close()
}
