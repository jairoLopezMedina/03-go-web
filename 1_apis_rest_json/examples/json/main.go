package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	stringJSON := `{"name": "John", "age": 30, "cars":["Ford", "BMW", "Fiat"]}`

	bytesJSON := []byte(stringJSON)


	// Convert into maps
	var mapJSON map[string]any
	err := json.Unmarshal(bytesJSON, &mapJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", mapJSON)

	value, ok := mapJSON["age"]
	if !ok {
		fmt.Println("age key not found")
		return
	}

	ageFloat, ok := value.(float64)
	if !ok{
		fmt.Println()
	}

	fmt.Println(ageFloat)
}