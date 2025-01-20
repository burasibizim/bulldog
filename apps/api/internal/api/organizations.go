package api

import (
	"api/internal/models"
	"api/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

func (a *api) getOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
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

	organizations, err := a.repository.Organizations.GetAllOrganizations(userID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	utils.JSONResponse(w, http.StatusOK, organizations)
}

func (a *api) createOrganizationHandler(w http.ResponseWriter, r *http.Request) {
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

	var payload models.OrganizationBody
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	org := models.Organization{
		Name:        payload.OrganizationName,
		Description: payload.OrganizationDescription,
		CreatedBy:   userID,
	}

	organization, err := a.repository.Organizations.CreateOrganization(&org)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	user := &models.OrganizationUserCreated{
		OrganizationId: organization.Id,
		UserId:         userID,
		Role:           "owner",
	}

	organizationUser, err := a.repository.OrganizationUsers.CreateOrganizationUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	err = a.cloudService.CodebuildService.CreateCodebuild(organization.Id)
	if err != nil {
		log.Printf("Error codebuild: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, organizationUser)
}

func (a *api) updateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.OrganizationInfo
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	organization, err := a.repository.Organizations.UpdateOrganization(&payload)
	if err != nil {
		log.Printf("Error updating organization: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to update organization")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, organization)
}

func (a *api) deleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.DeleteOrganizationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := a.repository.Organizations.DeleteOrganization(&payload)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to delete organization")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, result)
}
