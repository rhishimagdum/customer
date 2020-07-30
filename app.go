package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router       *mux.Router
	messanger    Messanger
	CustomerRepo CustomerDao
}

//Initialize ...initializes connectipn and mux
func (a *App) Initialize(repo CustomerDao, m Messanger) {
	a.CustomerRepo = repo
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.messanger = m
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
	customers := a.CustomerRepo.GetAll()
	respondWithJSON(w, http.StatusOK, customers)
}

// GetOne ...Get customer by id
func (a *App) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	cust := a.CustomerRepo.GetById(id)
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
	a.CustomerRepo.Create(cust)

	message, _ := json.Marshal(cust)
	a.messanger.PublishMessage(string(message))
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
