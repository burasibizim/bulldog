package auth

import (
	"api/pkg/utils"
	"context"
	"net/http"
)

type OrganizationUsersRepository interface {
	GetOrganizationUserRole(organizationID, userID string) (Role, error)
}

type AuthService struct {
	organizationUsersRepo OrganizationUsersRepository
}

func NewAuthService(repo OrganizationUsersRepository) *AuthService {
	return &AuthService{
		organizationUsersRepo: repo,
	}
}

func (s *AuthService) GetUserRole(organizationID, userID string) (Role, error) {
	return s.organizationUsersRepo.GetOrganizationUserRole(organizationID, userID)
}

func (s *AuthService) HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// RequirePermission creates a middleware that checks if the user has the required permission
func RequirePermission(authService *AuthService, requiredPermission Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
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

			organizationID := r.Header.Get("organization_id")
			if organizationID == "" {
				utils.JSONError(w, http.StatusBadRequest, "Organization ID is required")
				return
			}

			role, err := authService.GetUserRole(organizationID, userID)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, "Error checking user role")
				return
			}

			if !authService.HasPermission(Role(role), requiredPermission) {
				utils.JSONError(w, http.StatusForbidden, "Permission denied")
				return
			}

			// Add role to context for potential use in handlers
			ctx := context.WithValue(r.Context(), "userRole", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
