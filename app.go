package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type App struct {
	Router  *mux.Router
	Session *gocql.Session
}

//Initialize ...initializes connectipn and mux
func (a *App) Initialize(keyspace string) {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	a.Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra connection done")
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/customers", a.getAll).Methods("GET")
	a.Router.HandleFunc("/customers", a.insert).Methods("Post")
	a.Router.HandleFunc("/customers/{id}", a.getOne).Methods("GET")
}

func (a *App) Run(addr string) {
	http.ListenAndServe(":80", a.Router)
}

func (a *App) getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getAll(): Get all customers")
	customers := getCustomers(a.Session)
	json.NewEncoder(w).Encode(customers)
}

// GetOne ...Get customer by id
func (a *App) getOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	fmt.Println("getOne(): Get customer", id)
	c := Customer{ID: id}
	c.getCustomer(a.Session)
	if c != (Customer{}) {
		json.NewEncoder(w).Encode(c)
	} else {
		http.Error(w, "Customer not found",
			http.StatusNotFound)
	}
}

// Insert ...Add new customer
func (a *App) insert(w http.ResponseWriter, r *http.Request) {
	var cust Customer
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	err = json.Unmarshal(body, &cust)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	cust.createCustomer(a.Session)
	fmt.Println("insert(): Inserting new customer", cust)
	json.NewEncoder(w).Encode(cust)
}
