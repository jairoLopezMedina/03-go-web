package internal

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrFieldQuality = errors.New("field quality")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type ProductService interface {
	FindAllProducts() (products map[int]Product, err error)
	FindProductById(id int) (product Product, err error)
	FindProductsFilteredByPrice(priceGt float64) (products map[int]Product, err error)
	SaveProduct(product *Product) (err error)
	UpdateProduct(product *Product) (err error)
	DeleteProduct(id int) (err error)
}