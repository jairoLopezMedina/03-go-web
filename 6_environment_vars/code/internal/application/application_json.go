package application

import (
	"clase4/internal/handler"
	"clase4/internal/repository"
	"clase4/internal/service"
	"clase4/internal/storage"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type ConfigAppJSON struct {
	Addr       string
	Token      string
	FilePath   string
	FilePathLastId int
	LayoutDate string
}

func NewApplicationJSON(cfg *ConfigAppJSON) *ApplicationJSON {
	// default configuration
	defaultCfg := &ConfigAppJSON{
		Addr:           ":8080",
		Token:          "",
		FilePath:       "",
		FilePathLastId: 0,
		LayoutDate:     "",
	}

	if cfg != nil {
		if cfg.Addr != "" {
			defaultCfg.Addr = cfg.Addr
		}
		if cfg.Token != "" {
			defaultCfg.Token = cfg.Token
		}
		if cfg.FilePath != "" {
			defaultCfg.FilePath = cfg.FilePath
		}
		if cfg.FilePathLastId != 0 {
			defaultCfg.FilePathLastId = cfg.FilePathLastId
		}
		if cfg.LayoutDate != "" {
			defaultCfg.LayoutDate = cfg.LayoutDate
		}
	}
	return &ApplicationJSON{
		rt:       chi.NewRouter(),
		addr:     defaultCfg.Addr,
		token:    defaultCfg.Token,
		filePath: defaultCfg.FilePath,
		layoutDate: defaultCfg.LayoutDate,
	}	
}

type ApplicationJSON struct {
	rt *chi.Mux
	addr string
	token string
	filePath string
	layoutDate string
}

func (a *ApplicationJSON) SetUp() (err error) {
	// dependencies

	// - check if file exists or can be opened
	f, err := os.Open(a.filePath)
	if err != nil {
		return
	}
	f.Close()

	st := storage.NewStorageProductJSON(a.filePath, a.layoutDate)

	rp := repository.NewProductStore(st, a.layoutDate)

	sv := service.NewProductDefault(rp)

	hd := handler.NewDefaultProducts(sv, a.token)

	// server
	// - routes
	a.rt.Get("/ping", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	a.rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll())
		r.Get("/{productId}", hd.GetById())
		r.Get("/search", hd.Search())
		r.Post("/", hd.Create())
		r.Put("/{id}", hd.Update())
		r.Patch("/{id}", hd.UpdatePartial())
		r.Delete("/{id}", hd.Delete())
	})
	return
}

func (a *ApplicationJSON) Run() (err error) {
	err = http.ListenAndServe(a.addr, a.rt)
	return
}

