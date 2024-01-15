package internal

import "errors"

var (
	ErrProductCodeValueAlreadyExists = errors.New("product code value already exists")
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	FindAll() (products map[int]Product, err error)
	FindById(id int) (product Product, err error)
	FindByPriceGreaterThan(price float64) (products map[int]Product, err error)
	Save(product *Product) (err error)
	Update(product *Product) (err error)
	Delete(id int) (err error)
}