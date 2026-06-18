package handler

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"
	"starter-api-golang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RoleLookup is a minimal role representation for dropdowns.
type RoleLookup struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LookupHandler struct {
	roleRepo       domainRepo.RoleRepository
	permissionRepo domainRepo.PermissionRepository
}

func NewLookupHandler(roleRepo domainRepo.RoleRepository, permissionRepo domainRepo.PermissionRepository) *LookupHandler {
	return &LookupHandler{roleRepo: roleRepo, permissionRepo: permissionRepo}
}

// Roles returns all roles as id+name pairs — no permission check, auth only.
// Used by any page that needs a role dropdown (e.g. Add User form).
func (h *LookupHandler) Roles(c *gin.Context) {
	roles, _, err := h.roleRepo.FindAll(domainRepo.RoleFilter{Page: 1, PerPage: 9999})
	if err != nil {
		response.InternalServerError(c, "failed to retrieve roles")
		return
	}

	lookup := make([]RoleLookup, len(roles))
	for i, r := range roles {
		lookup[i] = RoleLookup{ID: r.ID, Name: r.Name}
	}

	response.OK(c, "roles retrieved", lookup)
}

// Permissions returns all permissions (flat list) — no permission check, auth only.
// Used by the Role form to render the permission checkbox tree.
func (h *LookupHandler) Permissions(c *gin.Context) {
	perms, err := h.permissionRepo.FindAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve permissions")
		return
	}

	if perms == nil {
		perms = []entity.Permission{}
	}

	response.OK(c, "permissions retrieved", perms)
}
