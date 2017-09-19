package main

import (
	"database/sql"
)

type context struct {
	dbh *sql.DB
}

