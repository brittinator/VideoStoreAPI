package customer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Customers is out in-memory "DB"
var Customers []Customer

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

func CreateCustomer(c Customer) {
	Customers = append(Customers, c)
}

func GetAll() []Customer {
	return Customers
}

func GetCustomer(id string) Customer {
	for _, c := range Customers {
		if c.ID == id {
			return c
		}
	}
	return Customer{}
}

func DeleteCustomer(id string) bool {
	for i, c := range Customers {
		if c.ID == id {
			Customers = append(Customers[:i], Customers[i+1:]...)
			return true
		}
	}
	return false
}

func UpdateCustomer(updatedCust Customer) bool {
	for i, c := range Customers {
		if c.ID == updatedCust.ID {
			Customers[i] = updatedCust
			return true
		}
	}
	return false
}

// FilterBy is the entry point into interacting and filtering results with the Customer model.
func FilterBy(filter, variable string) ([]Customer, error) {
	switch filter {
	case "city":
		return filterByCity(variable)
	case "name":
		return filterByName(variable)
	case "id":
		return filterByID(variable)
	case "state":
		return filterByState(variable)
	case "phone":
		return filterByPhone(variable)
	default:
		return nil, fmt.Errorf("Filter %v not a valid Customer filter", filter)
	}
}

func filterByName(name string) ([]Customer, error) {
	var customers []Customer
	for _, c := range Customers {
		if c.Name == name {
			customers = append(customers, c)
		}
	}
	return customers, nil
}

func filterByID(id string) ([]Customer, error) {
	var customers []Customer
	for _, c := range Customers {
		if c.ID == id {
			customers = append(customers, c)
		}
	}
	return customers, nil
}

func filterByCity(city string) ([]Customer, error) {
	var customers []Customer
	for _, c := range Customers {
		if c.City == city {
			customers = append(customers, c)
		}
	}
	return customers, nil
}

func filterByState(state string) ([]Customer, error) {
	var customers []Customer
	for _, c := range Customers {
		if c.State == state {
			customers = append(customers, c)
		}
	}
	return customers, nil
}

func filterByPhone(phone string) ([]Customer, error) {
	var customers []Customer
	for _, c := range Customers {
		if c.Phone == phone {
			customers = append(customers, c)
		}
	}
	return customers, nil
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

func SeedCustomers() []Customer {
	raw, err := ioutil.ReadFile("./seeds/customers.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Customer
	json.Unmarshal(raw, &c)
	return c
}
