package main

import (
	"fmt"
	"net/http"
)

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

	// - unique handler
	uniqueHandler := func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case "GET":
				w.Write([]byte("Index"))
			case "POST":
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