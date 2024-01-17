package service

import (
	"clase4/internal"
	"fmt"
	"time"
)

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

type ProductDefault struct {
	rp internal.ProductRepository
}

func (p *ProductDefault) FindAllProducts() (products []internal.Product, err error) {
	products, err = p.rp.FindAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (p *ProductDefault) FindProductById(id int) (product internal.Product, err error) {
	product, err = p.rp.FindById(id)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w:id", internal.ErrProductNotFound)
		}
		return
	}

	return
}

func (p *ProductDefault) FindProductsFilteredByPrice(priceGt float64) (products []internal.Product, err error) {
	products, err = p.rp.FindByPriceGreaterThan(priceGt)
	if err != nil {
		return
	}
	return
}

func (p *ProductDefault) SaveProduct(product *internal.Product) (err error){
	
	if err = validateProduct(product); err != nil{
		return
	}

	// save movie
	err = p.rp.Save(product)
	if err != nil {
		switch err {
		case internal.ErrProductCodeValueAlreadyExists:
			err = fmt.Errorf("%w: title", internal.ErrProductCodeValueAlreadyExists)
		}
		return err
	}

	return
}


func (p *ProductDefault) UpdateProduct(product *internal.Product) (err error) {
	
	if err = validateProduct(product); err != nil {
		return
	}
	
	err = p.rp.Upsert(product)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	
	return
}

func (p *ProductDefault) DeleteProduct(id int) (err error) {
	
	err = p.rp.Delete(id)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	
	return
}

func validateProduct(product *internal.Product) (err error) {
	if (*product).Name == "" {
		return fmt.Errorf("%w: name", internal.ErrFieldRequired)
	}
	
	if (*product).CodeValue == "" {
		return fmt.Errorf("%w: code value", internal.ErrFieldRequired)
	}

	if (*product).Expiration == "" {
		return fmt.Errorf("%w: expiration", internal.ErrFieldRequired)
	}

	_, err = time.Parse("02/01/2006", (*product).Expiration)
	if err != nil {
		return fmt.Errorf("%w: expiration", internal.ErrFieldQuality)
	}

	if (*product).Quantity < 0 {
		return fmt.Errorf("%w: price", internal.ErrFieldQuality)
	}

	if (*product).Price < 0 {
		return fmt.Errorf("%w: quantity", internal.ErrFieldQuality)
	}

	return
}