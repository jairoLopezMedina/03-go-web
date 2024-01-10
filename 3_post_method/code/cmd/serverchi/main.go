package main

import (
	"clase3/code/internal/application"
	"fmt"
)

func main() {
	
	app := application.NewServerChi("")

	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
	
}