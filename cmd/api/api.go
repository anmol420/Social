package main

import (
	"context"
	"errors"
	"expvar"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anmol420/Social/internal/auth"
	"github.com/anmol420/Social/internal/env"
	"github.com/anmol420/Social/internal/mailer"
	"github.com/anmol420/Social/internal/ratelimiter"
	"github.com/anmol420/Social/internal/store"
	"github.com/anmol420/Social/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	ratelimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redis       redisConfig
	ratelimiter ratelimiter.Config
}

type redisConfig struct {
	addr    string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	username string
	password string
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.StringGetEnv("FRONTEND_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(app.ratelimiterMiddleware)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// r.With(app.basicAuthMiddleware()).Get("/health", app.healthCheckHandler)
		r.Get("/health", app.healthCheckHandler)
		r.With(app.basicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.authTokenMiddleware)
			r.Post("/create", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostByIDHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.authTokenMiddleware)
				r.Get("/", app.getUserByIDHandler)
				r.Post("/follow", app.followUserHandler)
				r.Post("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.authTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	// graceful shutdown
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.logger.Infow("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()
	app.logger.Infow("Server Started", "addr", app.config.addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	app.logger.Infow("server has stopped", "addr", app.config.addr)
	return nil
}
