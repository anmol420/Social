package main

import (
	"log"

	"github.com/anmol420/Social/internal/env"
)

func main() {
	addr := env.GetEnv("ADDR")
	if addr == "" {
		log.Fatal("ADDR environment variable is not set")
	}
	cfg := config{
		addr: addr,
	}
	app := &application{
		config: cfg,
	}
	if err := app.run(app.mount()); err != nil {
		log.Fatal(err)
	}
}