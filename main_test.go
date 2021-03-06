package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customer "github.com/VideoStoreAPI/models/customers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB() []customer.Customer {
	return []customer.Customer{customer.Customer{ID: "1", Name: "B", City: "Honu"}, customer.Customer{ID: "2", Name: "C", City: "Zandu"}}
}

func TestGetCustomers(t *testing.T) {
	customer.Customers = setupDB()

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
			200,
			true,
		},
	}

	customer.Customers = setupDB()

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

func TestCreateCustomerHandler(t *testing.T) {
	testCases := []struct {
		name    string
		id      int
		payload []byte

		expected    string
		code        int
		expectError bool
	}{
		{
			"overwriting works",
			1,
			[]byte(`{"id":"1","name":"Bonobo"}`),

			`"Customer Bonobo successfully created"
`,
			http.StatusCreated,
			false,
		},
		{
			"net new customer",
			100,
			[]byte(`{"id":"100","name":"Mjobo"}`),

			`"Customer Mjobo successfully created"
`,
			http.StatusCreated,
			true,
		},
		{
			"malformed request",
			100,
			[]byte(`{"i""""`),

			"",
			http.StatusBadRequest,
			true,
		},
	}

	customer.Customers = setupDB()

	router := mux.NewRouter()
	addCustomerRoutes(router)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", fmt.Sprintf("/customers/%v", tt.id), bytes.NewBuffer(tt.payload))
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.expected, w.Body.String())
		})
	}
}

func TestDeleteCustomerHandler(t *testing.T) {
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
			`"Customer with id 1 successfully deleted"
`,
			200,
			false,
		},
		{
			"does not exist",
			100,
			`"No customer with id 100 found"
`,
			404,
			true,
		},
	}

	customer.Customers = setupDB()

	router := mux.NewRouter()
	addCustomerRoutes(router)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/customers/%v", tt.id), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.expected, w.Body.String())
		})
	}
}

func TestUpdateCustomerHandler(t *testing.T) {
	testCases := []struct {
		name    string
		id      int
		payload []byte

		expected    string
		code        int
		expectError bool
	}{
		{
			"does exist",
			1,
			[]byte(`{"id":"1","name":"Bonobo"}`),

			`"Customer Bonobo updated to: {1 Bonobo       0}"
`,
			200,
			false,
		},
		{
			"does not exist",
			100,
			[]byte(`{"id":"100","name":"Bonobo"}`),

			`"No customer with id 100 found"
`,
			304,
			true,
		},
		{
			"bad request",
			100,
			[]byte(`{"id":"100","name:"Bonobo"}`),
			`400
`,
			400,
			true,
		},
	}

	customer.Customers = setupDB()

	router := mux.NewRouter()
	addCustomerRoutes(router)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("PUT", fmt.Sprintf("/customers/%v", tt.id), bytes.NewBuffer(tt.payload))
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.expected, w.Body.String())
		})
	}
}
