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

func getCustomerHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	// b/c there's no DB
	for _, c := range Customers {
		if c.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	fmt.Println("Here")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&Customer{})

}

func getAllCustomersHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println(Customers)
	json.NewEncoder(w).Encode(&Customers)
}

func createCustomerHandler(w http.ResponseWriter, req *http.Request) {
	var c Customer
	err := json.NewDecoder(req.Body).Decode(&c)
	// https://stackoverflow.com/questions/33238518/what-could-happen-if-i-dont-close-response-body-in-golang
	defer req.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Customers = append(Customers, c)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fmt.Sprintf("Customer %v successfully created", c.Name))
}

func deleteCustomerHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for i, c := range Customers {
		if c.ID == params["id"] {
			Customers = append(Customers[:i], Customers[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fmt.Sprintf("Customer %v successfully deleted", c.Name))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(fmt.Sprintf("No customer with id %v found", params["id"]))
}

func updateCustomerHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for i, c := range Customers {
		if c.ID == params["id"] {
			var updatedCust Customer
			err := json.NewDecoder(req.Body).Decode(&updatedCust)
			defer req.Body.Close()

			if err != nil {
				json.NewEncoder(w).Encode(http.StatusBadRequest)
				return
			}
			Customers[i] = updatedCust
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fmt.Sprintf("Customer %v updated to: %v", c.Name, updatedCust))
			return
		}
	}

	w.WriteHeader(http.StatusNotModified)
	json.NewEncoder(w).Encode(fmt.Sprintf("No customer with id %v found", params["id"]))
}

func addCustomerRoutes(r *mux.Router) {
	cusRouter := r.PathPrefix("/customers").Subrouter()
	cusRouter.HandleFunc("", getAllCustomersHandler).Methods("GET")
	cusRouter.HandleFunc("/{id:[0-9]+}", getCustomerHandler).Methods("GET")
	cusRouter.HandleFunc("/{id:[0-9]+}", updateCustomerHandler).Methods("PUT")
	cusRouter.HandleFunc("/{id:[0-9]+}", createCustomerHandler).Methods("POST")
	cusRouter.HandleFunc("/{id:[0-9]+}", deleteCustomerHandler).Methods("DELETE")

}

func main() {
	Customers = seedCustomers()

	// logger := log.New(os.Stdout, "[VideoStoreAPI]", 0)
	router := mux.NewRouter()

	addCustomerRoutes(router)

	// fallthrough if no paths matched
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("I'm sorry Dave, I can't help you.\n'%s' not found\n", req.URL)))
	})
	// actually activate the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
