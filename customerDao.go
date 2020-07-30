package main

type CustomerDao interface {
	Create(c Customer)
	GetById(i int) Customer
	GetAll() []Customer
}
