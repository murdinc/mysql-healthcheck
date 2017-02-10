package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

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
		CheckSlaveStatus bool
		MaxQueries       int
		Port             string
	}
}

var cnf Config

func main() {

	err := gcfg.ReadFileInto(&cnf, "/etc/mysql/mysql-healthcheck.cnf")
	if err != nil {
		log.Printf("%v", err)
		return
	}

	http.HandleFunc("/", healthcheck)                  // redirect all urls to the handler function
	http.ListenAndServe(":"+cnf.HealthCheck.Port, nil) // listen for connections at port 9999 on the local machine

	//log.Printf("%v", cnf.Mysql.UserName)
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
		//log.Printf("%v", err)
		log.Print("MySQL Running: [false], Query Count: [0],  Slave Running: [0]")
		w.Header().Set("Error", "Unable to connect to mysql")
		w.WriteHeader(500)
		return
	}

	openQueries, slaveRunning := queryCheck(db)

	log.Printf("MySQL Running: [true], Query Count: [%v],  Slave Running: [%v]", openQueries, slaveRunning)

	w.Header().Set("Server", "Mysql-Health-Check")
	w.Header().Set("Queries", strconv.Itoa(openQueries))
	w.Header().Set("Slave-Running", strconv.FormatBool(slaveRunning))

	if openQueries >= cnf.HealthCheck.MaxQueries || (cnf.HealthCheck.CheckSlaveStatus && slaveRunning == false) {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}

}

func queryCheck(db *sql.DB) (int, bool) {

	var queryCount int
	var slaveRunning bool
	slaveRunning = false

	globalStats, err := db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		log.Printf("%v", err)
	}
	defer globalStats.Close()

	for globalStats.Next() {
		var name string
		var value string
		err = globalStats.Scan(&name, &value)

		// Check current queries
		if name == "Threads_connected" {
			queryCount, err = strconv.Atoi(value)

		}
		// Check if slave running
		if name == "Slave_running" && value == "ON" {
			slaveRunning = true
		}
	}

	return queryCount, slaveRunning
}
