package application

import (
	"clase4/internal/handler"
	"clase4/internal/repository"
	"clase4/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewDefaultHTTP(addr string) *DefaultHTTP {
	return &DefaultHTTP{
		addr: addr,
	}	
}

type DefaultHTTP struct {
	addr string
}

func (h *DefaultHTTP) Run() (err error) {
	rp := repository.NewProductMap()

	sv := service.NewProductDefault(rp)

	hd := handler.NewDefaultProducts(sv)

	rt := chi.NewRouter()

	rt.Get("/ping", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll())
		r.Get("/{productId}", hd.GetById())
		r.Get("/search", hd.Search())
		r.Post("/", hd.Create())
	})

	err = http.ListenAndServe(h.addr, rt)
	return
}

