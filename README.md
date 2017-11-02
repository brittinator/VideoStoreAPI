# VideoStoreAPI
Go version of VideoStoreAPI


Seed in-memory DB: https://www.chazzuka.com/2015/03/load-parse-json-file-golang/

First create a router.

Match up the routes to their handlers

Call listenAndServe to actually use the router


# Calling endpoints

Get: `curl localhost:8080`

Update: `curl --request PUT -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle-oRAMA"}' http://localhost:8080/customers/1 {"id":"1","name":"brittzle-oRAMA"}`

Create `curl -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle"}' http://localhost:8080/customers/1`

Delete `curl --request DELETE -H "Content-Type: application/json" -d '{"id":"1", "name":"brittzle-oRAMA"}' http://localhost:8080/customers/1`