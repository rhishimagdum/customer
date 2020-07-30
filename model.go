package main

// Customer ... customer struct
type Customer struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Address   string `json:"address,omitempty"`
}
