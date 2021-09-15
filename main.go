package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Mysql struct {
		UserName string
		Password string
		Database string
		Port     string
	}
	HealthCheck struct {
		Port             string
	}
}

var cnf Config

func main() {

	err := gcfg.ReadFileInto(&cnf, "/etc/mysql-healthcheck.cnf")

	if err != nil {
		log.Printf("%v", err)
		return
	}

	http.HandleFunc("/", healthcheck)                  // redirect all urls to the handler function
	http.ListenAndServe(":"+cnf.HealthCheck.Port, nil) // listen for connections at port 9999 on the local machine

}

func healthcheck(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", cnf.Mysql.UserName+":"+cnf.Mysql.Password+"@/"+cnf.Mysql.Database)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()

	if err != nil {
		log.Print("ERROR: Unable to connect to mysql")
		log.Print(err)
		w.Header().Set("Error", "Unable to connect to mysql")
		w.WriteHeader(500)
		return
	}
	
	log.Print("SUCCESS: Connected to Mysql")

	var status bool
	
	status = queryCheck(db)
	
	if status == false {
		w.WriteHeader(500)
	} else{
		w.WriteHeader(200)
	}

	return
}

func queryCheck(db *sql.DB) (bool) {
	var val int
	err := db.QueryRow("SELECT 1;").Scan(&val)

	if err != nil {
		log.Print("ERROR:  Query failure: ", err)
		return false
	}
	log.Print("SUCCESS: Query Successful: ", val)
    return true
}
