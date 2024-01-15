package handler

import (
	"clase4/internal"
	"clase4/platform/web/request"
	"clase4/platform/web/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewDefaultProducts(sv internal.ProductService) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
	}
}

type DefaultProducts struct {
	sv internal.ProductService
}

type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyRequestProductJSON struct {
	Name    	string  `json:"name"`
	Quantity 	int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price 		float64 `json:"price"`
}

func (d *DefaultProducts) HealthCheck() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		response.Text(w, http.StatusOK, "pong")
	}
}

func (d *DefaultProducts) GetAll() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		products, err := d.sv.FindAllProducts()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "products",
			"data": products,
		})	
	}
}

func (d *DefaultProducts) GetById() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		
		// Request
		id, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// Process
		p, err := d.sv.FindProductById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Error(w, http.StatusNotFound, "Product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// Response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product",
			"data": ProductJSON{
				ID: p.ID,
				Name: p.Name,
				Quantity: p.Quantity,
				CodeValue: p.CodeValue,
				IsPublished: p.IsPublished,
				Expiration: p.Expiration,
				Price: p.Price,
			},
		})
	}
}

func (d *DefaultProducts) Search() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Request
		priceGt, err := strconv.ParseFloat(r.URL.Query().Get("priceGt"), 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// Process 
		products, err := d.sv.FindProductsFilteredByPrice(priceGt)
		if err != nil {
			switch {
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// Response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "products",
			"data": products,
		})
	}
}

func (d *DefaultProducts) Create() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var body BodyRequestProductJSON
		if err := request.JSON(r, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		product := internal.Product{
			Name: body.Name,
			Quantity: body.Quantity,
			CodeValue: body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration: body.Expiration,
			Price: body.Price,
		}

		if err := d.sv.SaveProduct(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrFieldQuality):
				response.Error(w, http.StatusBadRequest, "invalid body")
			case errors.Is(err, internal.ErrProductAlreadyExists):
				response.Error(w, http.StatusConflict, "movie already exists")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		//response 
		data := ProductJSON{
			ID: product.ID,
			Quantity: product.Quantity,
			CodeValue: product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration: product.Expiration,
			Price: product.Price,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "product saved",
			"data": data,
		})
	}
}

func (d *DefaultProducts) Update() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// request
		// - get id from path
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// - get body to []byte
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		// - get body to map[string]any
		var bodyMap map[string]any
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		// - validate body
		if err := validateKeyExistence(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		// - get body
		var body BodyRequestProductJSON
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
		}

		// process
		product := internal.Product{
			ID:          id,
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price: 		 body.Price,
		}

		// - update movie
		if err := d.sv.UpdateProduct(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrFieldQuality):
				response.Text(w, http.StatusBadRequest, "invalid body")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		
		// response
		// - map to ProductJSON
		data := ProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price: 		 product.Price,
		}

		// - response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product updated",
			"data": data,
		})
	}
}

func (d *DefaultProducts) UpdatePartial() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		// request
		// - get id from path
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// - get product
		product, err := d.sv.FindProductById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
				
			}
			return
		}

		// process
		// - map internal.Product to BodyRequestProductJSON
		reqBody := BodyRequestProductJSON{
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price: 		 product.Price,
		}

		// - get body
		if err := request.JSON(r, &reqBody); err != nil {
			response.Text(w, http.StatusBadRequest, " invalid body")
			return
		}

		// - map internal.Product
		product = internal.Product{
			ID: id,
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  reqBody.Expiration,
			Price: 		 reqBody.Price,
		}

		// - update product
		if err := d.sv.SaveProduct(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrFieldQuality):
				response.Text(w, http.StatusBadRequest, "invalid body")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - map to ProductJSON
		data := ProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price: 		 product.Price,
		}

		// - response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product updated",
			"data": data,
		})

	}
}

func (d *DefaultProducts) Delete() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// request
		// - get id from path
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "invalid id",
				"data": nil,
			})
			return
		}

		// process
		// -delete product
		if err := d.sv.DeleteProduct(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		response.Text(w, http.StatusNoContent, "product deleted")
	}
}

func validateKeyExistence(mp map[string]any, keys ...string) (err error) {
	for _, k := range keys {
		if _, ok := mp[k]; !ok {
			return fmt.Errorf("key %s not found", k)
		}
	}
	return
}