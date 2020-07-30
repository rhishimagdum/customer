package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(CustomerImplCassandra{}, KafkaMessanger{})
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if err := Session.Query(tableCreationQuery).Exec(); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	Session.Query("TRUNCATE TABLE customers").Exec()
}

func addTestCustomers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		suffix := count + i
		if err := Session.Query("INSERT INTO customer.customers(id, first_name, last_name, address) VALUES(?, ?, ?, ?)", suffix, "First-"+strconv.Itoa(suffix), "Last-"+strconv.Itoa(suffix), "Address-"+strconv.Itoa(suffix)).Exec(); err != nil {
			fmt.Println("Not inserted")
		}
	}
}

const tableCreationQuery = "CREATE TABLE IF NOT EXISTS customers (id int PRIMARY KEY,first_name text,last_name text,address text)"

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/customers", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	body := response.Body.String()
	if body != "null" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"id": 2, "firstName":"John", "lastName": "snow", "Address":"Winterfel"}`)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["firstName"] != "John" {
		t.Errorf("Expected product name to be 'John'. Got '%v'", m["firstName"])
	}

	if m["lastName"] != "snow" {
		t.Errorf("Expected product price to be 'snow'. Got '%v'", m["lastName"])
	}

	// the id is compared to 2.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 2.0 {
		t.Errorf("Expected product ID to be '2'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentCustomer(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/customers/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Customer not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addTestCustomers(1)
	req, _ := http.NewRequest("GET", "/customers/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
