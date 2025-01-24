package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	PlainToken string `json:"token"`
}

func (app *application) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(payload); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	user := store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	exp := time.Duration(time.Hour * 2)
	plaintoken := uuid.New().String()
	hash := sha256.Sum256([]byte(plaintoken))
	hashtoken := hex.EncodeToString(hash[:])

	userwithToken := UserWithToken{
		User:       &user,
		PlainToken: plaintoken,
	}
	if err := app.store.Users.CreateAndInvite(r.Context(), &user, exp, hashtoken); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.mail.Send(user.Username, user.Email); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]any{"user created successfully": userwithToken})
}
