package main

import (
	"time"

	"github.com/anmol420/Social/internal/db"
	"github.com/anmol420/Social/internal/env"
	"github.com/anmol420/Social/internal/mailer"
	"github.com/anmol420/Social/internal/store"
	"go.uber.org/zap"
)

func main() {
	addr := env.StringGetEnv("ADDR")
	databaseAddr := env.StringGetEnv("DATABASE_ADDR")
	maxOpenConns := env.IntegerGetEnv("MAX_OPEN_CONNS")
	maxIdleConns := env.IntegerGetEnv("MAX_IDLE_CONNS")
	maxIdleTime := env.StringGetEnv("MAX_IDLE_TIME")
	mailerFromEmail := env.StringGetEnv("MAILER_FROM_EMAIL")
	mailerRegion := env.StringGetEnv("MAILER_REGION")
	frontendUrl := env.StringGetEnv("FRONTEND_URL")
	basicAuthUsername := env.StringGetEnv("AUTH_BASIC_USERNAME")
	basicAuthPassword := env.StringGetEnv("AUTH_BASIC_PASSWORD")
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
		frontendURL: frontendUrl,
		auth: authConfig{
			basic: basicConfig{
				username: basicAuthUsername,
				password: basicAuthPassword,
			},
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
	mail, err := mailer.NewSESClient(mailerFromEmail, mailerRegion)
	if err != nil {
		logger.Fatal(err)
	}
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mail,
	}
	if err := app.run(app.mount()); err != nil {
		logger.Fatal(err)
	}
}
