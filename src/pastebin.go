package main

import (
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/husobee/vestigo"
	"github.com/kr/pretty"
	"math/rand"
	"log"
	//	"os"
//	"strings"
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

func errCheck(e error) {
	if (e != nil) {
		log.Fatal(e)
	}
}

func GetHandler(w http.ResponseWriter, request *http.Request){
	id := vestigo.Param(request, "postid");
	fmt.Println(id)
	rows, err := sqlh.Query("select post from posts where id = ?", id)
	errCheck(err)

	// Since we selected on the key the result is unambiguous...
	var post string
	// Safe because we selected on primary key, only 1 post
	rows.Next()
	err = rows.Scan(&post)
	errCheck(err)
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200);

	w.Write([]byte(post))
}

func PostHandler(w http.ResponseWriter, request *http.Request) {

	request.ParseForm()

//	fmt.Println(request.PostForm)
	
	if (sqlh == nil) { fmt.Println("EVERYTHING IS RUINED"); return}
	
	tx, err := sqlh.Begin();
	if (err != nil) {
		log.Fatal(err)
	}
	
	insert, err := tx.Prepare("insert into posts(id, post) values(?,?)")
	if (err != nil) {
		log.Fatal(err)
	}
	
	defer insert.Close()

//	data := strings.Join(request.PostForm[0], "")

	pretty.Println(request.PostForm)
	key := genKey()
//	_, err = insert.Exec(key, data);

	if (err != nil) {
		log.Fatal(err)
	}
	
	err = tx.Commit()
	if (err != nil) {
		log.Fatal(err)
	}

	w.Header().Set("Location", "/" + key)
	w.WriteHeader(303)
	
}

func main() {
	var err error
	sqlh, err = sql.Open("sqlite3", "pastebin.db")

	if (err != nil) {
		log.Fatal(err)
	}
	err = sqlh.Ping()
	if (err != nil) {
		log.Fatal(err)
	}

	router := vestigo.NewRouter()

	router.Get("/:postid", GetHandler)
	router.Put("/", PostHandler)
	http.ListenAndServe(":8000", router)

	sqlh.Close()
}
