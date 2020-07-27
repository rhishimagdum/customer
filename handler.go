package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAll ...Get all customers
func GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all customers")
	var customers []Customer
	m := map[string]interface{}{}
	iter := Session.Query("SELECT * FROM customers").Iter()
	for iter.MapScan(m) {
		customers = append(customers, Customer{
			ID:        m["id"].(int),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Address:   m["address"].(string),
		})
		m = map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(customers)
}

// GetOne ...Get customer by id
func GetOne(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	m := map[string]interface{}{}
	params := mux.Vars(r)

	iter := Session.Query("SELECT * FROM customers where id=?", params["id"]).Iter()
	for iter.MapScan(m) {
		customer = Customer{
			ID:        m["id"].(int),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Address:   m["address"].(string),
		}
	}

	if customer != (Customer{}) {
		json.NewEncoder(w).Encode(customer)
	} else {
		http.Error(w, "Customer not found",
			http.StatusNotFound)
	}
}

// Insert ...Add new customer
func Insert(w http.ResponseWriter, r *http.Request) {
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
	if err := Session.Query("INSERT INTO customers(id, first_name, last_name, address) VALUES(?, ?, ?, ?)",
		cust.ID, cust.FirstName, cust.LastName, cust.Address).Exec(); err != nil {
		fmt.Println("Error while inserting customer")
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(cust)
}
