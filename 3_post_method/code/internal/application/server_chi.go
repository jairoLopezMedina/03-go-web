package application

import (
	"clase3/code/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewServerChi(address string) *ServerChi {
	defaultAddress := ":8080"
	if address != "" {
		defaultAddress = address
	}

	return &ServerChi{
		address: defaultAddress,
	}
}

type ServerChi struct {
	address string
}

func (s *ServerChi) Run() error {
	rt := chi.NewRouter()

	hd, err := handler.NewHandler()
	if err != nil {
		return err
	}

	rt.Get("/ping", hd.HealthCheck())

	rt.Route("/products", func(rt chi.Router) {
		rt.Get("/", hd.GetProducts())

		rt.Get("/{productId}", hd.GetProductById())

		rt.Get("/search", hd.Search())

		rt.Post("/", hd.Create())
	})

	return http.ListenAndServe(s.address, rt)
}