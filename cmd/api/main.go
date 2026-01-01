package main

import (
	"time"

	"github.com/anmol420/Social/internal/db"
	"github.com/anmol420/Social/internal/env"
	"github.com/anmol420/Social/internal/store"
	"go.uber.org/zap"
)

func main() {
	addr := env.StringGetEnv("ADDR")
	databaseAddr := env.StringGetEnv("DATABASE_ADDR")
	maxOpenConns := env.IntegerGetEnv("MAX_OPEN_CONNS")
	maxIdleConns := env.IntegerGetEnv("MAX_IDLE_CONNS")
	maxIdleTime := env.StringGetEnv("MAX_IDLE_TIME")
	cfg := config{
		addr: addr,
		db: dbConfig{
			addr:         databaseAddr,
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  maxIdleTime,
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
	}
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Database Connection Established!")
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}
	if err := app.run(app.mount()); err != nil {
		logger.Fatal(err)
	}
}
