package user_usecase

import (
	"context"
	"sort"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
)

func (uu *UserUsecase) ListRoleAccess(ctx context.Context) ([]userdto.ListRoleAccessResponse, error) {
	// Ambil semua permission yang ada
	allPerms, err := uu.userRepo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	// Ambil semua role beserta permissions yang dimiliki
	roles, err := uu.userRepo.GetAllRolesWithPermissions(ctx)
	if err != nil {
		return nil, err
	}

	// Buat lookup: page -> actions[]
	allPages := buildAllPageActions(allPerms)

	var result []userdto.ListRoleAccessResponse
	for _, role := range roles {
		matrix := userdto.ListRoleAccessResponse{
			Role:   role.Role,
			Access: initAccessMatrix(allPages),
		}

		// Tandai permission yang dimiliki role sebagai true
		for _, perm := range role.Permissions {
			if perm.Page == "" || perm.Action == "" {
				continue // Skip data invalid
			}
			if actions, ok := matrix.Access[perm.Page]; ok {
				actions[perm.Action] = true
			}
		}

		result = append(result, matrix)
	}

	return result, nil
}

// Membuat daftar page -> actions unik dari seluruh permission
func buildAllPageActions(perms []entity.Permission) map[string][]string {
	pageActions := make(map[string]map[string]struct{})

	for _, p := range perms {
		if p.Page == "" || p.Action == "" {
			continue
		}
		if _, ok := pageActions[p.Page]; !ok {
			pageActions[p.Page] = make(map[string]struct{})
		}
		pageActions[p.Page][p.Action] = struct{}{}
	}

	// Ubah ke map[string][]string
	result := make(map[string][]string)
	for page, actionsMap := range pageActions {
		var actions []string
		for action := range actionsMap {
			actions = append(actions, action)
		}
		sort.Strings(actions) // optional: sort alphabetically
		result[page] = actions
	}
	return result
}

// Inisialisasi matriks akses default: semua false
func initAccessMatrix(allPages map[string][]string) map[string]map[string]bool {
	access := make(map[string]map[string]bool)
	for page, actions := range allPages {
		access[page] = make(map[string]bool)
		for _, action := range actions {
			access[page][action] = false
		}
	}
	return access
}
