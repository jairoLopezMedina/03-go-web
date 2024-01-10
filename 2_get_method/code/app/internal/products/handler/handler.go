package handler

import (
	"clase2/internal/products/storage"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	data map[int]storage.Product
}

type Response struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func NewHandler() (*ProductHandler, error) {
	f, err := os.Open("./../products.json")
	if err != nil {
		return nil, err
	}

	defer f.Close()
	
	jsonFile := storage.NewJSON(f)
	pArr, err := jsonFile.Read()
	if err != nil {
		return nil, err
	}

	productsM := make(map[int]storage.Product)
	for _, p := range pArr {
		productsM[p.ID] = p
	}

	return &ProductHandler{
		data: productsM,
	}, nil
}

// Method HealthCheck for GET /ping route
func (h *ProductHandler) HealthCheck() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		w.WriteHeader(code)
		w.Write([]byte("pong"))
	}
}

func (h *ProductHandler) GetProducts() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// jsonData, err := json.Marshal(h.data)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h.data)
		//w.Write(jsonData)
	}
}

func (h *ProductHandler) GetProductById() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "productId")
		id, _ := strconv.Atoi(idStr)

		p, ok := h.data[id]
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{ Message: "Product not found", Data: nil })
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}
}

func (h *ProductHandler) Search() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		pgtParam := r.URL.Query().Get("priceGt")

		priceGt, err := strconv.ParseFloat(pgtParam, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{ Message: "Price not valid", Data: nil })
			return
		}

		result := []storage.Product{}

		for _, p := range h.data {
			if p.Price > priceGt {
				result = append(result, p)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
