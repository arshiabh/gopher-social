package main

import (
	"net/http"

	"github.com/arshiabh/gopher-social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
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
	createdUser := app.store.Users.CreateAndInvite(r.Context(), &user, "token-123")
	if err := jsonResponse(w, http.StatusCreated, createdUser); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
