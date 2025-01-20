package api

import (
	"api/internal/models"
	"api/pkg/utils"
	"log"
	"net/http"
	"time"
)

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var payload models.User
	if !utils.DecodeJSONBody(w, r, &payload) {
		return
	}

	now := time.Now()
	user := &models.User{
		Id:        userID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	success, err := a.repository.Users.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, success)
}

func (a *api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := a.repository.Users.GetUserWithID(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

func (a *api) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var payload models.User
	if !utils.DecodeJSONBody(w, r, &payload) {
		return
	}

	success, err := a.repository.Users.UpdateUser(userID, payload.FirstName, payload.LastName)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.JSONResponse(w, http.StatusOK, success)
}
