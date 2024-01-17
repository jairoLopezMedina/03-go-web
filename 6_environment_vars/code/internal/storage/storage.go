package storage

import (
	"clase4/internal"
	"errors"
)

var (
	ErrStorageProductTimeLayout = errors.New("storage: time layout invalid")
)

type StorageProduct interface {
	// ReadAll reads and returns all products from a file
	ReadAll() (p []internal.Product, err error)
	// WriteAll writes all products to a file
	WriteAll(p []internal.Product) (err error)
}