package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

const (
	kafkaConn        = "localhost:9092"
	topic            = "cust"
	cassandraCluster = "127.0.0.1"
)

type App struct {
	Router   *mux.Router
	Session  *gocql.Session
	Producer sarama.AsyncProducer
	Consuner sarama.Consumer
}

//Initialize ...initializes connectipn and mux
func (a *App) Initialize(keyspace string) {
	var err error
	cluster := gocql.NewCluster(cassandraCluster)
	cluster.Keyspace = keyspace
	a.Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra connection done")
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.initializeProducer()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/customers", a.getAll).Methods("GET")
	a.Router.HandleFunc("/customers", a.insert).Methods("Post")
	a.Router.HandleFunc("/customers/{id}", a.getOne).Methods("GET")
}

func (a *App) initializeProducer() {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	a.Producer, _ = sarama.NewAsyncProducer([]string{kafkaConn}, config)
	fmt.Println("KAFKA connection done")
}

func (a *App) Run(addr string) {
	http.ListenAndServe(":80", a.Router)
}

func (a *App) getAll(w http.ResponseWriter, r *http.Request) {
	customers := getCustomers(a.Session)
	respondWithJSON(w, http.StatusOK, customers)
}

// GetOne ...Get customer by id
func (a *App) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	cust := Customer{ID: id}
	cust.getCustomer(a.Session)

	if len(cust.Address) > 0 {
		respondWithJSON(w, http.StatusOK, cust)
	} else {
		respondWithError(w, http.StatusNotFound, "Customer not found")
	}
}

// Insert ...Add new customer
func (a *App) insert(w http.ResponseWriter, r *http.Request) {
	var cust Customer
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading request body")
	}
	err = json.Unmarshal(body, &cust)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading request body")
	}
	cust.createCustomer(a.Session)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(body),
	}
	a.Producer.Input() <- msg
	respondWithJSON(w, http.StatusOK, cust)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
