package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Person struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func main() {

	// register endpoints
	// - pong
	handlerPong := func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Println("GET /ping")
			fmt.Println("method:", r.Method)
			fmt.Println("url:", r.URL)
			fmt.Println("header:", r.Header)
			w.Write([]byte("pong"))
		}
	}

	// - register endpoint /ping
	http.HandleFunc("/ping", handlerPong)

	// - greetings
	handlerGreetings := func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			// request metadata
			fmt.Println("POST /greetings")
			fmt.Println("method:", r.Method)
			fmt.Println("url:", r.URL)
			fmt.Println("header:", r.Header)

			// body reader
			dc := json.NewDecoder(r.Body)

			// decode into Person struct
			var p Person
			if err := dc.Decode(&p); err != nil {
				fmt.Println(err)
				return
			}

			w.Write([]byte("Hello " + p.FirstName + " " + p.LastName))

		}
	}

	// register greetings handler
	http.HandleFunc("/greetings", handlerGreetings)

	// - unique handler
	uniqueHandler := func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case "GET":
				w.Write([]byte("Index"))
		}
	}

	// register unique handler
	http.HandleFunc("/", uniqueHandler)


	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
		return
	}

	
}