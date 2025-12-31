package main

import (
	"log"

	"github.com/anmol420/Social/internal/db"
	"github.com/anmol420/Social/internal/env"
	"github.com/anmol420/Social/internal/store"
)
func main() {
	addr := env.StringGetEnv("ADDR")
	databaseAddr := env.StringGetEnv("DATABASE_ADDR")
	maxOpenConns := env.IntegerGetEnv("MAX_OPEN_CONNS")
	maxIdleConns := env.IntegerGetEnv("MAX_IDLE_CONNS")
	maxIdleTime := env.StringGetEnv("MAX_IDLE_TIME")
	cfg := config{
		addr:   addr,
		db: dbConfig{
			addr:         databaseAddr,
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  maxIdleTime,
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database Connection Established!")
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}
	if err := app.run(app.mount()); err != nil {
		log.Fatal(err)
	}
}
