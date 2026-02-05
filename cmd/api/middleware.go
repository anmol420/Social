package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anmol420/Social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) basicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, fmt.Errorf("AuthHeader is missing"))
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedError(w, r, fmt.Errorf("AuthHeader is malformed"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}
			username := app.config.auth.basic.username
			pass := app.config.auth.basic.password
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedError(w, r, fmt.Errorf("Invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("AuthHeader is missing"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("AuthHeader is malformed"))
			return
		}
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r.Context())
		post := getPostFromCtx(r)
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redis.enabled {
		return app.store.Users.GetByID(ctx, userID)
	}
	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (app *application) ratelimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.ratelimiter.Enabled {
			if allow, retryAfter := app.ratelimiter.Allow(r.RemoteAddr); !allow {
				app.ratelimitExceedError(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}