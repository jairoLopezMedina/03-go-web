package storage

import (
	"clase4/internal"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type StorageProductJSON struct {
	// filePath is the path of the file to read
	filePath string
	// layoutDate is the layout of the expiration date
	layoutDate string

}

// NewStorageProductJSON creates and returns a new storage product json instance
func NewStorageProductJSON(filePath, layoutDate string) *StorageProductJSON {
	// default config
	defaultLayoutDate := time.DateOnly
	if layoutDate != "" {
		defaultLayoutDate = layoutDate
	}

	return &StorageProductJSON{filePath: filePath, layoutDate: defaultLayoutDate}
}

// ProductAttributesJSON is an abstract representation of a product properties
type ProductAttributesJSON struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	CodeValue string `json:"code_value"`
	IsPublished bool `json:"is_published"`
	Expiration string `json:"expiration"`
	Price float64 `json:"price"`
}

// ReadAll reads and returns all products from a file
func (s *StorageProductJSON) ReadAll() (p []internal.Product, err error) {
	// open file
	f, err := os.Open(s.filePath)
	if err != nil {
		return
	}

	defer f.Close()

	// decode
	pr := []ProductAttributesJSON{}
	err = json.NewDecoder(f).Decode(&pr)
	if err != nil {
		return
	}

	// serialization
	for _, product := range pr {
		_, err = time.Parse(s.layoutDate, product.Expiration)
		if err != nil {
			fmt.Println(product.Expiration)
			fmt.Println(s.layoutDate)
			err = ErrStorageProductTimeLayout
			return
		}

		p = append(p, internal.Product{
			ID: 		 product.ID, 
			Name: 		 product.Name, 
			Quantity: 	 product.Quantity, 
			CodeValue: 	 product.CodeValue, 
			IsPublished: product.IsPublished, 
			Expiration:  product.Expiration, 
			Price : 	 product.Price,
		})
	}

	return
}


// WriteAll writes all products to a file
func (s *StorageProductJSON) WriteAll(p []internal.Product) (err error) {
	
	// open file
	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}

	// encode
	err = json.NewEncoder(f).Encode(p)
	if err != nil {
		return
	}
	
	return
}

