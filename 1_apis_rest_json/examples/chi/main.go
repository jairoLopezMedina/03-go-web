package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	// create a router with chi
	router := chi.NewRouter()

	// Create a new endpoitn GET "/hello-world"
	router.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		//set code
		w.WriteHeader(200)

		// set body
		w.Write([]byte("Hello world!"))
	})

	http.ListenAndServe(":8000", router)
}