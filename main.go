package main

import (
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

var Session *gocql.Session

func init() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "customer"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra connection done")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers", GetAll).Methods("GET")
	r.HandleFunc("/customers", Insert).Methods("Post")
	r.HandleFunc("/customers/{id}", GetOne).Methods("GET")

	http.ListenAndServe(":80", r)
}
