package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	PostgresHostName     = "ec2-54-228-97-176.eu-west-1.compute.amazonaws.com"
	PostgresPort         = 5432
	PostgresUsername     = "idelciuumoufmw"
	PostgresPassword     = "58aab6636158c74dc63c7fcd4655701a65afc4cbbea04fddccde96deacc1bae5"
	PostgresDatabaseName = "d74pjn9g01ncq1"
)

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := NewClient()
	if err != nil {
		io.WriteString(w, "failed")
	}

	io.WriteString(w, "success")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func NewClient() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://idelciuumoufmw:58aab6636158c74dc63c7fcd4655701a65afc4cbbea04fddccde96deacc1bae5@ec2-54-228-97-176.eu-west-1.compute.amazonaws.com:5432/d74pjn9g01ncq1")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("postgres: connected to database has been successfully")
	return db, nil
}
