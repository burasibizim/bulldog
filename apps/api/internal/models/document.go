package models

import (
	"time"
)

type DocumentStatus string

const (
	StatusPublished DocumentStatus = "published"
	StatusDraft     DocumentStatus = "draft"
	StatusScheduled DocumentStatus = "scheduled"
)

type DocumentCreateRequest struct {
	Title   string         `json:"title" validate:"required"`
	Preview string         `json:"preview"`
	Status  DocumentStatus `json:"status" validate:"required,oneof=published draft scheduled"`
	Tags    []string       `json:"tags"`
	File    []byte         `json:"file"`
}

type DocumentUpdateRequest struct {
	Id      string         `json:"id" validate:"required"`
	Name    string         `json:"name,omitempty"`
	Title   string         `json:"title,omitempty"`
	Preview string         `json:"preview,omitempty"`
	Status  DocumentStatus `json:"status,omitempty" validate:"omitempty,oneof=published draft scheduled"`
	Tags    []string       `json:"tags,omitempty"`
}

type Document struct {
	OrganizationId string         `json:"organization_id"`
	Id             string         `json:"id"`
	Title          string         `json:"title"`
	Preview        string         `json:"preview"`
	Status         DocumentStatus `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Tags           []string       `json:"tags"`
	StoragePath    string         `json:"storage_path"`
	CreatedBy      string         `json:"created_by"`
	UpdatedBy      string         `json:"updated_by"`
}
