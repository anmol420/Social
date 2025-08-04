package main

import (
	"log"

	"github.com/anmol420/Social/internal/env"
	"github.com/anmol420/Social/internal/store"
)

func main() {
	addr := env.GetEnv("ADDR")
	if addr == "" {
		log.Fatal("ADDR environment variable is not set")
	}
	cfg := config{
		addr: addr,
	}
	store := store.NewStorage(nil)
	app := &application{
		config: cfg,
		store:  store,
	}
	if err := app.run(app.mount()); err != nil {
		log.Fatal(err)
	}
}