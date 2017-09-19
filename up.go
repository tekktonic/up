package main;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
)


func main() {
	var ctx context

	dbh, err := sql.Open("sqlite3", "up.db");

	if (err != nil) {
		log.Fatal(err)
	}
	ctx.dbh = dbh
	
	post := post{id: "", author: "tekk@up.tekk.in", post: "please fix everything", favorites: 0, replyto: ""}

	fmt.Println(String(&post));
	
	id, _ := Put(dbh, &post);
	fmt.Println("trying to retrieve");
	gotten,_ := Get(dbh,id);
	fmt.Print(String(gotten))
}
