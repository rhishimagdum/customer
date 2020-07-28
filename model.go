package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

// Customer ... customer struct
type Customer struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Address   string `json:"address,omitempty"`
}

func (cust *Customer) getCustomer(session *gocql.Session) {
	m := map[string]interface{}{}

	iter := session.Query("SELECT * FROM customers where id=?", cust.ID).Iter()
	for iter.MapScan(m) {
		cust.ID = m["id"].(int)
		cust.FirstName = m["first_name"].(string)
		cust.LastName = m["last_name"].(string)
		cust.Address = m["address"].(string)
	}
}

func (cust *Customer) createCustomer(session *gocql.Session) {
	if err := session.Query("INSERT INTO customers(id, first_name, last_name, address) VALUES(?, ?, ?, ?)",
		cust.ID, cust.FirstName, cust.LastName, cust.Address).Exec(); err != nil {
		fmt.Println("Error while inserting customer")
		fmt.Println(err)
	}
}

func getCustomers(session *gocql.Session) []Customer {
	fmt.Println("Getting all customers")
	var customers []Customer
	m := map[string]interface{}{}
	iter := session.Query("SELECT * FROM customers").Iter()
	for iter.MapScan(m) {
		customers = append(customers, Customer{
			ID:        m["id"].(int),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Address:   m["address"].(string),
		})
		m = map[string]interface{}{}
	}
	return customers
}
