package main

import (
	"log"
	"net/http"
	"strconv"
	"code.google.com/p/gcfg"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type Config struct {
	Mysql struct {
		UserName string
		Password string
		Database string
		CheckSlaveStatus bool
		Port string
	}
	HealthCheck struct {
		MaxQueries	int
		MaxDelay	int
		Port string
	}
}

var cnf Config

func main() {

	//var cnf Config
	err := gcfg.ReadFileInto(&cnf, "/etc/mysql/mysql-healthcheck.cnf")
	if err != nil {
		log.Printf("%v", err)
		return
	}

	http.HandleFunc("/", healthcheck) // redirect all urls to the handler function
	http.ListenAndServe(":" + cnf.HealthCheck.Port, nil) // listen for connections at port 9999 on the local machine

	//log.Printf("%v", cnf.Mysql.UserName)
}


func healthcheck(w http.ResponseWriter, r *http.Request) { 

	db, err := sql.Open("mysql", cnf.Mysql.UserName + ":" + cnf.Mysql.Password + "@/" + cnf.Mysql.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Printf("%v", err)
	}

	slaveCheck(db)
	openQueries, slaveRunning := queryCheck(db)

	log.Printf("Query Count: [%v],  Slave Running: [%v]", openQueries, slaveRunning)


}

func queryCheck(db *sql.DB) (int, bool) {

	var queryCount int
	var slaveRunning bool

	globalStats, err := db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		log.Printf("%v", err)
	}
	defer globalStats.Close()

	for globalStats.Next() {
		var name string
		var value string
		err = globalStats.Scan(&name, &value )
		//log.Printf("Name: %v, Value: %v", name, value)

		// Check current queries
		if name == "Threads_connected" {
			queryCount, err = strconv.Atoi(value)

		}
		// Check if slave running
		if name == "Slave_running" && value == "OFF" {
			slaveRunning = false

		}
	}

	return queryCount, slaveRunning
}

func slaveCheck(db *sql.DB) {

	slaveStats, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		log.Printf("%v", err)
	}
	defer slaveStats.Close()

	log.Printf("%+v", slaveStats)

	for slaveStats.Next() {
		var name string
		var value string
		err = slaveStats.Scan(&name, &value )
		log.Printf("Name: %v, Value: %v", name, value)
	}

}

