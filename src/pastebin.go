package main

import (
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/husobee/vestigo"
	"github.com/kr/pretty"
	"math/rand"
	"mime/multipart"
	"log"
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

	defer rows.Close()
	// Since we selected on the key the result is unambiguous...
	var post string
	// Safe because we selected on primary key, only 1 post
	rows.Next()
	err = rows.Scan(&post)
	if (err != nil) {
		w.WriteHeader(404);
		w.Write(([]byte)("Paste not found."))
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200);

	w.Write([]byte(post))
}

func PostHandler(w http.ResponseWriter, request *http.Request) {

	maxsize := (int64)(2 << 21)
	request.ParseMultipartForm(maxsize)

	var file multipart.File
	var h *multipart.FileHeader
	var err error
	file, h, err = request.FormFile("post")

	if (h.Size > (2 << 21)) {
		w.WriteHeader(500)
		w.Write(([]byte)("Too big"))
		return;
	}
	rawdata := make([]byte, h.Size)

	file.Read(rawdata)
	
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
	
	// Only one key we need to extract, apparently best way.
//	for k := range request.PostMultipartForm {
//		post = k
	//	}

	data := (string)(rawdata)
	pretty.Println(data)
//	pretty.Println((string)(request.PostFormValue("post")))
	key := genKey()
	_, err = insert.Exec(key, data);

	if (err != nil) {
		log.Fatal(err)
	}
	
	err = tx.Commit()
	if (err != nil) {
		log.Fatal(err)
	}

	w.Write(([]byte)(key + "\r\n"))

	
}

func main() {
	var err error
	fmt.Println("Opening db");
	sqlh, err = sql.Open("sqlite3", "pastebin.db")

	if (err != nil) {
		log.Fatal(err)
	}

		fmt.Println("Creating vestigo")
	router := vestigo.NewRouter()
	fmt.Println("Registering get..")
	router.Get("/:postid", GetHandler)
		fmt.Println("Registering post..")
	router.Post("/", PostHandler)

	fmt.Println("Trying to serve..")
	http.ListenAndServe(":8080", router)

	sqlh.Close()
}
