# VideoStoreAPI
Go version of VideoStoreAPI


Seed in-memory DB: https://www.chazzuka.com/2015/03/load-parse-json-file-golang/

First create a router.

Match up the routes to their handlers

Call listenAndServe to actually use the router

# How to run
` go run main.go`

# Calling endpoints

Get: `curl localhost:8080`

Update: `curl --request PUT -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle-oRAMA"}' http://localhost:8080/customers/1 {"id":"1","name":"brittzle-oRAMA"}`

Create `curl -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle"}' http://localhost:8080/customers/1`

Delete `curl --request DELETE -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle-oRAMA"}' http://localhost:8080/customers/1`

Filtering
Currently only implemented filter by city. `curl localhost:8080/customers/filter_by=city/"Anchorage"` Should return 4 results.


## Reading Resources

Here are a few of the great resources I used while making this

* JSON parsing: https://www.chazzuka.com/2015/03/load-parse-json-file-golang/
* API https://thenewstack.io/make-a-restful-json-api-go/
* API https://www.thepolyglotdeveloper.com/2016/07/create-a-simple-restful-api-with-golang/
* API + Mux: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
* Mux Routes: https://gowebexamples.com/routes-using-gorilla-mux/


