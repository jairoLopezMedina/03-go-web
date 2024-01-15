package repository

import (
	"clase4/internal"
	"clase4/platform/file"
	"fmt"
	"os"
)

func NewProductMap() *ProductMap{

	f, err := os.Open("../../products.json")
	if err != nil {
		fmt.Println(err)
		return &ProductMap{
			db: make(map[int]internal.Product),
			lastId: 0,
		}
	}

	defer f.Close()

	jsonFile := storage.NewJSON(f)
	productsArr, err := jsonFile.Read()
	if err != nil {
		fmt.Println(err)
		return &ProductMap{
			db: make(map[int]internal.Product),
			lastId: 0,
		}
	}

	products := make(map[int]internal.Product)
	for _, p := range productsArr {
		products[p.ID] = p
	}

	return &ProductMap{
		db: products,
		lastId: len(productsArr),
	}
}

type ProductMap struct {
	db map[int]internal.Product
	lastId int
}

func (p *ProductMap) FindAll() (products map[int]internal.Product, err error) {
	products = p.db
	return 
}

func (p *ProductMap) FindById(id int) (product internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}

	return
}

func (p *ProductMap) FindByPriceGreaterThan(price float64) (products map[int]internal.Product, err error) {

	products = make(map[int]internal.Product)
	for _, pr := range p.db {
		if pr.Price > price {
			products[pr.ID] = pr
		}
	}

	return
}

func (p *ProductMap) Save(product *internal.Product) (err error) {
	// validate product
	if err = p.validateProductCodeValue((*product).CodeValue); err != nil {
		return
	}
	// autoincrement
	// - increment id
	(*p).lastId++
	// - set id
	(*product).ID = (*p).lastId

	// store movie
	(*p).db[(*product).ID] = *product

	return

}

func (p *ProductMap) Update(product *internal.Product) (err error) {
	_, ok := p.db[(*product).ID]
	if !ok {
		err = internal.ErrProductNotFound
	}

	p.db[(*product).ID] = *product
	return
}

func (p *ProductMap) Delete(id int) (err error) {
	_, ok := p.db[id]
	if !ok {
		err = internal.ErrProductNotFound
	}

	delete(p.db, id)

	return
}

func (p *ProductMap) validateProductCodeValue(codeValue string) (err error){
	for _, p := range p.db {
		if p.CodeValue == codeValue {
			return internal.ErrProductCodeValueAlreadyExists
		}
	}

	return
}