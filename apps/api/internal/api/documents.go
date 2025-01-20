package api

import (
	"api/internal/models"
	"api/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (a *api) createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	organization_id := r.Header.Get("organization_id")
	if organization_id == "" {
		utils.JSONError(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	title := r.FormValue("title")
	preview := r.FormValue("preview")
	statusStr := r.FormValue("status")
	tagsJSON := r.FormValue("tags")

	var tags []string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		http.Error(w, "Invalid tags format", http.StatusBadRequest)
		return
	}

	var status models.DocumentStatus
	switch statusStr {
	case string(models.StatusPublished), string(models.StatusDraft), string(models.StatusScheduled):
		status = models.DocumentStatus(statusStr)
	default:
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	payload := models.DocumentCreateRequest{
		Title:   title,
		Preview: preview,
		Status:  status,
		Tags:    tags,
	}

	err := a.cloudService.StorageService.UploadMarkdownHandler(w, r, organization_id, "/documents")
	if err != nil {
		log.Println("Failed to create document", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create document")
		return
	}

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

	document := &models.Document{
		OrganizationId: organization_id,
		Title:          payload.Title,
		Preview:        payload.Preview,
		Status:         payload.Status,
		Tags:           payload.Tags,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdDocument, err := a.repository.Documents.CreateDocument(document)
	if err != nil {
		log.Println("Failed to create document", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create document")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, createdDocument)
}

func (a *api) getDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	organizationId := r.Header.Get("organization_id")
	documents, err := a.repository.Documents.GetDocuments(organizationId)
	if err != nil {
		log.Println("Failed to get documents", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to get documents")
		return
	}

	utils.JSONResponse(w, http.StatusOK, documents)
}

func (a *api) updateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	organizationId := r.Header.Get("organization_id")
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

	var payload models.DocumentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	document := &models.Document{
		OrganizationId: organizationId,
		Id:             payload.Id,
		Title:          payload.Title,
		Preview:        payload.Preview,
		Status:         payload.Status,
		Tags:           payload.Tags,
		UpdatedBy:      userID,
		UpdatedAt:      time.Now(),
	}

	err := a.repository.Documents.UpdateDocument(document)
	if err != nil {
		log.Println("Failed to update document", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to update document")
		return
	}

	utils.JSONResponse(w, http.StatusOK, document)
}

func (a *api) getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	organizationId := r.Header.Get("organization_id")
	vars := mux.Vars(r)
	documentId := vars["document_id"]

	document, err := a.repository.Documents.GetDocument(organizationId, documentId)
	if err != nil {
		log.Println("Failed to get document", err)
		utils.JSONError(w, http.StatusInternalServerError, "Failed to get document")
		return
	}

	utils.JSONResponse(w, http.StatusOK, document)
}
