package api

import (
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/utils"
	"log"
	"net/http"
	"strings"
)

func (a *api) addUserToOrganization(w http.ResponseWriter, r *http.Request) {
	var payload models.AddUserOrganizationPayload
	if !utils.DecodeJSONBody(w, r, &payload) {
		return
	}
	organizationID := r.Header.Get("organization_id")

	user, err := a.repository.Users.GetUserWithEmail(payload.Email)
	if err != nil {
		log.Printf("Error fetching user by email: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	organizationUser := &models.OrganizationUserCreated{
		OrganizationId: organizationID,
		UserId:         user.Id,
		Role:           payload.Role,
	}

	createdOrganizationUser, err := a.repository.OrganizationUsers.CreateOrganizationUser(organizationUser)
	if a.handleOrganizationUserError(w, err, "creating organization user") {
		return
	}

	response := map[string]interface{}{
		"user_id":    createdOrganizationUser.UserId,
		"email":      payload.Email,
		"role":       payload.Role,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}

	utils.JSONResponse(w, http.StatusCreated, response)
}

func (a *api) updateUserRoleInOrganization(w http.ResponseWriter, r *http.Request) {
	var payload models.UpdateUserRolePayload
	if !utils.DecodeJSONBody(w, r, &payload) {
		return
	}

	organizationID := r.Header.Get("organization_id")
	err := a.repository.OrganizationUsers.UpdateOrganizationUserRole(organizationID, payload.UserID, payload.Role)
	if a.handleOrganizationUserError(w, err, "updating user role") {
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User role updated successfully"})
}

func (a *api) deleteUserFromOrganization(w http.ResponseWriter, r *http.Request) {
	var payload models.DeleteOrganizationUserPayload
	if !utils.DecodeJSONBody(w, r, &payload) {
		return
	}

	organizationID := r.Header.Get("organization_id")
	err := a.repository.OrganizationUsers.DeleteOrganizationUser(organizationID, payload.UserID)
	if a.handleOrganizationUserError(w, err, "deleting user from organization") {
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted from organization successfully"})
}

func (a *api) getOrganizationUserInfo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	userID := parts[len(parts)-1]

	if userID == "" {
		utils.JSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	organizationID := r.Header.Get("organization_id")
	userInfo, err := a.repository.OrganizationUsers.GetOrganizationUserInfo(userID, organizationID)
	if a.handleOrganizationUserError(w, err, "getting user information") {
		return
	}

	utils.JSONResponse(w, http.StatusOK, userInfo)
}

func (a *api) listOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	organizationID := r.Header.Get("organization_id")
	if organizationID == "" {
		utils.JSONError(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	users, err := a.repository.OrganizationUsers.ListOrganizationUsers(organizationID)
	if a.handleOrganizationUserError(w, err, "listing organization users") {
		return
	}

	if len(users) == 0 {
		utils.JSONResponse(w, http.StatusOK, []models.OrganizationUserInfo{})
		return
	}

	utils.JSONResponse(w, http.StatusOK, users)
}

func (a *api) handleOrganizationUserError(w http.ResponseWriter, err error, operation string) bool {
	if err != nil {
		if _, ok := err.(repository.ErrOrganizationUserNotFound); ok {
			utils.JSONError(w, http.StatusNotFound, "User not found in the organization")
		} else {
			log.Printf("Error %s: %v", operation, err)
			utils.JSONError(w, http.StatusInternalServerError, "Failed to "+operation)
		}
		return true
	}
	return false
}
