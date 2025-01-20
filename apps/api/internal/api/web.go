package api

import (
	"api/pkg/utils"
	"net/http"
)

func (a *api) publishWebHandler(w http.ResponseWriter, r *http.Request) {
	organization_id := r.Header.Get("organization_id")
	if organization_id == "" {
		utils.JSONError(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	a.cloudService.CodebuildService.StartCodebuild(organization_id)
	utils.JSONResponse(w, http.StatusOK, "Build started")
}
