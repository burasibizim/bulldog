package handlers

import (
	"api/pkg/database"
	"api/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func TenantsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTenantHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// tenantHandler handles requests to /tenants/{tenantID}.
func TenantHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTenantHandler(w, r)
	case http.MethodPut:
		updateTenantHandler(w, r)
	case http.MethodDelete:
		deleteTenantHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createTenantHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest models.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//get user
	user, err := database.GetUser(createUserRequest.CreatorId)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	tenant := models.Tenant{
		Name: createUserRequest.OrganizationName,
	}

	//create tenant
	TenantId, err := database.CreateTenant(database.DB, &tenant)
	if err != nil {
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	//create tenant user
	err = database.CreateTenantUser(user, TenantId, "admin")
	if err != nil {
		http.Error(w, "Failed to create tenant user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tenant created"))
}

func getTenantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var db = database.DB
	tenantId := r.URL.Query().Get("id")
	tenant, err := database.GetTenantById(db, tenantId)

	if err != nil {
		http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"id":   tenant.ID,
		"name": tenant.Name,
	}
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}

func updateTenantHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := database.UpdateTenantName(database.DB, payload.Id, payload.Name)

	if err != nil {
		http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Tanant name updated: %s", payload.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func deleteTenantHandler(w http.ResponseWriter, r *http.Request) {
	var db = database.DB
	tenantId := r.URL.Query().Get("id")
	err := database.DeleteTenant(db, tenantId)

	if err != nil {
		http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Tenant id:%s deleted", tenantId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
