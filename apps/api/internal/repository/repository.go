package repository

import "api/pkg/database"

type Repository struct {
	DB                *database.DB
	Documents         *DocumentRepository
	Organizations     *OrganizationRepository
	OrganizationUsers *OrganizationUserRepository
	Users             *UserRepository
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{
		DB:                db,
		Documents:         NewDocumentRepository(db),
		Organizations:     NewOrganizationRepository(db),
		OrganizationUsers: NewOrganizationUserRepository(db),
		Users:             NewUserRepository(db),
	}
}
