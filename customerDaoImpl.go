package main

import (
	"fmt"
)

type CustomerImplCassandra struct {
}

func (dao CustomerImplCassandra) Create(c Customer) {
	fmt.Println("Create()", c)
	if err := Session.Query("INSERT INTO customers(id, first_name, last_name, address) VALUES(?, ?, ?, ?)",
		c.ID, c.FirstName, c.LastName, c.Address).Exec(); err != nil {
		fmt.Println("Error while inserting customer")
		fmt.Println(err)
	}
}

func (dao CustomerImplCassandra) GetById(i int) Customer {
	fmt.Println("GetById", i)
	m := map[string]interface{}{}
	var cust Customer
	iter := Session.Query("SELECT * FROM customers where id=?", i).Iter()
	for iter.MapScan(m) {
		cust.ID = m["id"].(int)
		cust.FirstName = m["first_name"].(string)
		cust.LastName = m["last_name"].(string)
		cust.Address = m["address"].(string)
	}
	return cust
}

func (dao CustomerImplCassandra) GetAll() []Customer {
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
	return customers
}
