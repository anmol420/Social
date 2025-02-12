package main

import (
	"log"
	"os"

	"github.com/anmol420/twitterbackend/internal/env"
)

func main() {
	env.LoadEnv()
	cfg := config{
		addr: os.Getenv("API_ADDR"),
	}
	app := application{
		config: cfg,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}