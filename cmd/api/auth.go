package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/anmol420/Social/internal/mailer"
	"github.com/anmol420/Social/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type ActivationMailData struct {
	Username       string `json:"username"`
	ActivationLink string `json:"activationLink"`
}

// dev only
type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// dev only
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationUrl := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	// mailer
	data := ActivationMailData{
		Username:       user.Username,
		ActivationLink: activationUrl,
	}
	if err := app.mailer.Send(mailer.UserActivationTemplate, user.Username, user.Email, data); err != nil {
		// rollback user creation if mail is not send (SAGA Pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.internalServerError(w, r, err)
		}

		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
