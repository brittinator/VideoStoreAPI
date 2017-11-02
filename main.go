package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Customer struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	RegisteredAt  string `json:"registered_at,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Phone         string `json:"phone,omitempty"`
	AccountCredit int    `json:"account_credit,float,omitempty"`
}

func (p Customer) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func seedCustomers() []Customer {
	raw, err := ioutil.ReadFile("./seeds/customers.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Customer
	json.Unmarshal(raw, &c)
	return c
}

var Customers []Customer

// Log spits out logs before returning a handler function
func Log(l *log.Logger, h http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, req)
		l.Printf("%s\t%s\t%s\t%s",
			req.Method, req.RequestURI, name, time.Since(start),
		)
	})
}

func getCustomersHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	// b/c there's no DB
	for _, c := range Customers {
		if c.ID == params["id"] {
			json.NewEncoder(w).Encode(c)
			break
		}
	}

	json.NewEncoder(w).Encode(&Customer{})

}

func getAllCustomersHandler(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(&Customers)
}

func createCustomersHandler(w http.ResponseWriter, req *http.Request) {
	var c Customer

	json.NewDecoder(req.Body).Decode(&c)
	Customers = append(Customers, c)

	json.NewEncoder(w).Encode(http.StatusCreated)
}

func deleteCustomerHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for i, c := range Customers {
		if c.ID == params["id"] {
			json.NewEncoder(w).Encode(c)
			Customers = append(Customers[:i], Customers[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}

func updateCustomerHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var newC Customer
	json.NewDecoder(req.Body).Decode(&newC)

	for i, c := range Customers {
		if c.ID == params["id"] {
			Customers[i] = newC
			json.NewEncoder(w).Encode(newC)
			break
		}
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}

func main() {
	Customers = seedCustomers()

	// logger := log.New(os.Stdout, "[VideoStoreAPI]", 0)
	router := mux.NewRouter()

	cusRouter := router.PathPrefix("/customers").Subrouter()
	cusRouter.HandleFunc("", getAllCustomersHandler).Methods("GET")
	cusRouter.HandleFunc("/{id:[0-9]+}", getCustomersHandler).Methods("GET")
	cusRouter.HandleFunc("/{id:[0-9]+}", updateCustomerHandler).Methods("PUT")
	cusRouter.HandleFunc("/{id:[0-9]+}", createCustomersHandler).Methods("POST")
	cusRouter.HandleFunc("/{id:[0-9]+}", deleteCustomerHandler).Methods("DELETE")

	// fallthrough if no paths matched
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("I'm sorry Dave, I can't help you.\n'%s' not found\n", req.URL)))
	})
	// actually activate the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
