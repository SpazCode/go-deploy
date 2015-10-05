package main

import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"net/http"
)

type Job struct {
	title string
	description string
	repo_link string
	test_command string
	deploy_location string
	deploy_password string
	deploy_username string
	current_build string
	last_run string
}

type Build struct {
	id string
	last_run string
	run_start_time time.Time
	log string
	status string
}

func setup() {
	session, error := mgo.Dial("localhost")
	if error != nil {
		panic(error)	
	} 

	defer session.Close()
}

func rootHandler(w http.ResponseWriter, r *http .Request) {
	fmt.Fprintf(w, "Welcome to the Deployment APIs")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":9000", nil)
}
