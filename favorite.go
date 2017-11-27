package main;

import (
	"log"
	"github.com/husobee/vestigo"
	"net/http"
)

func FavoriteCB(w http.ResponseWriter,r *http.Request) {
	pledge()
	dbh := ctx.dbh;

	id := vestigo.Param(r, "id");

	str := auth(r.Header.Get("X-Up-Auth"));

	if (str == "") {
		tx, err := dbh.Begin()
		if (err != nil) {
			log.Fatal(err);
		}


		insert, err := tx.Prepare("update posts set favorites = favorites + 1 where id = ?");
		if (err != nil) {
			log.Fatal(err)
		}

		result, err := insert.Exec(id);

		if (err != nil) {
			log.Fatal(err);
		}
		affected,_ := result.RowsAffected();

		if (affected != 1) {
			w.WriteHeader(400)
			w.Write(([]byte)("Bad post ID"))
			return;
		}

		err = tx.Commit()

		if (err != nil) {
			log.Fatal(err)
		}

		insert.Close()

		post := NewPost(config.Owner + " liked the post https://" + config.Server + "/post/" + id);

		Put(ctx.dbh, &post);

		DistributePost(post);

	}
	w.WriteHeader(401)
	w.Write(([]byte)(str))
}
