package main

import (
	"clase2/internal/products/handler"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	
	h, err := handler.NewHandler()
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Get("/ping", h.HealthCheck())

	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.GetProducts())

		r.Get("/{productId}", h.GetProductById())

		r.Get("/search", h.Search())
	})

	http.ListenAndServe(":8080", r)
}