package api

import (
	"api/pkg/utils"
	"log"
	"net/http"
)

func (a *api) getOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	organizationID := r.Header.Get("organization_id")
	organization, err := a.repository.Organizations.GetOrganization(organizationID)

	if err != nil {
		log.Printf("Error getting organization: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to get organization")
		return
	}

	utils.JSONResponse(w, http.StatusOK, organization)
}
