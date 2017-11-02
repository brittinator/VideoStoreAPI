package main

import (
	"fmt"
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

func TestGetCustomer(t *testing.T) {

	testCases := []struct {
		name string
		id   int

		expected    string
		code        int
		expectError bool
	}{
		{
			"does exist",
			1,
			"{\"id\":\"1\",\"name\":\"B\",\"city\":\"Honu\"}\n",
			200,
			false,
		},
		{
			"does not exist",
			100,
			"{}\n",
			404,
			true,
		},
	}

	Customers = setupDB()

	router := mux.NewRouter()
	addCustomerRoutes(router)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/customers/%v", tt.id), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.expected, w.Body.String())
		})
	}

}
