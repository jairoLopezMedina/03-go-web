package main

import (
	"clase4/internal/application"
	"fmt"
	"os"
)

func main() {
	// env

	// application
	// - config
	cfg := &application.ConfigAppJSON{
		Addr:  os.Getenv("ENV_ADDR"),
		Token: os.Getenv("API_TOKEN"),
		FilePath: os.Getenv("ENV_DB_FILE_PATH"),
		LayoutDate: os.Getenv("LAYOUT_DATE"),
	}
	app := application.NewApplicationJSON(cfg)
	if err := app.SetUp();err != nil {
		fmt.Println(err)
		return
	}

	// run 
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}