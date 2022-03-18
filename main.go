package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	PostgresHostName     = "ec2-54-228-97-176.eu-west-1.compute.amazonaws.com"
	PostgresPort         = 5432
	PostgresUsername     = "idelciuumoufmw"
	PostgresPassword     = "58aab6636158c74dc63c7fcd4655701a65afc4cbbea04fddccde96deacc1bae5"
	PostgresDatabaseName = "d74pjn9g01ncq1"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")

	pgConnString := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		PostgresHostName, PostgresPort, PostgresUsername, PostgresPassword, PostgresDatabaseName,
	)

	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		io.WriteString(w, "postgres: connected failed")
	}

	err = db.Ping()
	if err != nil {
		io.WriteString(w, "postgres: connected failed")
	}

	io.WriteString(w, "postgres: connected to database has been successfully")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
