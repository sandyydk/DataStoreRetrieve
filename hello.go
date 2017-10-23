package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type Entity struct {
	Values string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	kind := "input"
	var entity *Entity
	value := r.URL.Query().Get("input")

	projectID := "mytestproject-183711"
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	//	query := datastore.NewQuery("input")

	taskkey := datastore.NameKey(kind, "inputvalues", nil)

	entity = new(Entity)
	if err := client.Get(ctx, taskkey, entity); err != nil {
		log.Println("Failed to obtain data from datastore")
	}

	// There can be two approaches. First one is one where we will have only a single entity(row) with the key as "inputvalue"
	// to which all input values will be appended with ',' as a separator
	// Second, we will store a key called "inputCount" which will be incremented each time input is used as a param. Using this we know
	// how many entities (rows) are present for 'input'. Then we can stick to a pattern for key name like 'key1,key2,..key{inputCount}'.
	// This way we retrieve all the values for 'input'.
	// I chose first one becauese number of API calls to datastore is reduced -- single call in the first case.

	existingValues := entity.Values

	entity.Values = existingValues + "," + value

	if _, err := client.Put(ctx, taskkey, entity); err != nil {
		log.Println("Failed to update key")
	}

}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	kind := "input"
	var entity *Entity

	projectID := "mytestproject-183711"

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	entityKey := datastore.NameKey(kind, "inputvalues", nil)

	entity = new(Entity)
	if err := client.Get(ctx, entityKey, entity); err != nil {
		log.Println("Failed to obtain data from datastore")
	}

	if entity != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, entity.Values)
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/save", saveHandler)
	r.HandleFunc("retrieve", retrieveHandler)
	appengine.Main()
}
