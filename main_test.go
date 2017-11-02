package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB() []Customer {
	return []Customer{Customer{ID: "1", Name: "B", City: "Honu"}, Customer{ID: "2", Name: "C", City: "Zandu"}}
}

func TestGetCustomers(t *testing.T) {
	Customers = setupDB()

	router := mux.NewRouter()
	addCustomerRoutes(router)

	req, err := http.NewRequest("GET", "/customers", nil)
	require.NoError(t, err, "Creating Get /customers failed")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := "[{\"id\":\"1\",\"name\":\"B\",\"city\":\"Honu\"},{\"id\":\"2\",\"name\":\"C\",\"city\":\"Zandu\"}]\n"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestGetNonExistentCustomer(t *testing.T) {

}
