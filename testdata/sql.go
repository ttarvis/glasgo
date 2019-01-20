package main

import (
	"database/sql"
)

var db *sql.DB;

func sqlInject1(userInput string) {
	var args []string;

	_, _ = db.Exec(userInput, args);
	_, _ = db.Query(userInput, args);
}

func sqlInject2(userInput string) {
	var args []string;
	q := "SELECT user WHERE user =";
	q = q + "test";

	_, _ = db.Exec(q, args);
}

func main() {
	sqlInject1("test");
	sqlInject2("test");
}
