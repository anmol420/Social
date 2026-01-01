package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/anmol420/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userContextKey userKey = "user"

type FollowUserPayload struct {
	UserID int64 `json:"user_id"`
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user, err := app.store.Users.GetByID(r.Context(), userId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r.Context())

	// TODO: Revert back to auth userID from ctx
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, payload.UserID, followerUser.ID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := getUserFromContext(r.Context())

	// TODO: Revert back to auth userID from ctx
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, payload.UserID, unFollowedUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) *store.User {
	user := ctx.Value(userContextKey).(*store.User)
	return user
}
