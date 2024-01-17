package repository

import (
	"clase4/internal"
	"clase4/internal/storage"
	"time"
)

func NewProductStore(st storage.StorageProduct, layoutDate string) *ProductStore{

	defaultLaoutDate := time.DateOnly
	if layoutDate == "" {
		layoutDate = defaultLaoutDate
	}
	return &ProductStore{
		st: st,
		layoutDate: layoutDate,
	}
}

type ProductStore struct {
	// storage 
	st storage.StorageProduct
	lastId int
	layoutDate string
}

func (p *ProductStore) FindAll() (products []internal.Product, err error) {
	// get all products
	products, err = p.st.ReadAll()
	return 
}

func (p *ProductStore) FindById(id int) (product internal.Product, err error) {
	// get all products
	products, err := p.st.ReadAll()

	// verify if product exists
	for _, pr := range products {
		if pr.ID == id {
			product = pr
			return
		}
	}

	err = internal.ErrProductNotFound

	return
}

func (p *ProductStore) FindByPriceGreaterThan(price float64) (products []internal.Product, err error) {

	ps, err := p.st.ReadAll()
	for _, product := range ps {
		if product.Price > price {
			products = append(products, product)
		}
	}

	return
}

func (p *ProductStore) Save(product *internal.Product) (err error) {

	// get all products
	products, err := p.st.ReadAll()
	if err != nil {
		return
	}

	// validate product
	if err = p.validateProductCodeValue(products ,(*product).CodeValue); err != nil {
		return
	}

	// set last id based on products length
	// - verify if it has already been set
	if (*p).lastId == 0 {
		(*p).lastId = len(products)
	}

	// autoincrement
	// - increment id
	(*p).lastId++
	// - set id
	(*product).ID = (*p).lastId

	// store movie
	products = append(products, *product)

	// write all products
	err = p.st.WriteAll(products)
	if err != nil {
		return
	}

	return

}

func (p *ProductStore) Upsert(product *internal.Product) (err error) {
	// get all products
	ps, err := p.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	var exists bool; var index int
	for idx, pr := range ps {
		if pr.ID == product.ID {
			index = idx
			exists = true
			break
		}
	}

	// update or create product
	switch exists {
	case true:
		ps[index] = *product
	default:
		// set id
		if p.lastId == 0 { (*p).lastId = len(ps) }
		p.lastId++
		(*product).ID = p.lastId
		ps =append(ps, *product)
	}

	// write all products
	err = p.st.WriteAll(ps)
	if err != nil {
		return
	}

	return
}

func (p *ProductStore) Delete(id int) (err error) {
	// get all products
	products, err := p.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	var exists bool; var idx int
	for index, product := range products {
		if product.ID == id {
			idx = index
			exists = true
			break
		}
	}

	// check if product exists
	if !exists {
		err = internal.ErrProductNotFound
		return
	}

	// delete product
	products = append(products[:idx], products[idx+1:]...)

	// write all products
	err = p.st.WriteAll(products)
	if err != nil {
		return
	}

	return
}

func (p *ProductStore) validateProductCodeValue(products []internal.Product, codeValue string) (err error){
	for _, p := range products {
		if p.CodeValue == codeValue {
			return internal.ErrProductCodeValueAlreadyExists
		}
	}

	return
}