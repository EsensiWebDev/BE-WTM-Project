package user_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/userdto"
)

func (uu *UserUsecase) UpdateRoleAccess(ctx context.Context, req *userdto.UpdateRoleAccessRequest) error {
	roleID := getRoleID(req.Role)

	perm, err := uu.userRepo.GetPermissionByPageAction(ctx, req.Page, req.Action)
	if err != nil {
		return fmt.Errorf("permission not found: %s", err.Error())
	}

	hasAccess, err := uu.userRepo.HasRolePermission(ctx, roleID, perm.ID)
	if err != nil {
		return err
	}

	if req.Allowed && !hasAccess {
		return uu.userRepo.AddRolePermission(ctx, roleID, perm.ID)
	}

	if !req.Allowed && hasAccess {
		return uu.userRepo.RemoveRolePermission(ctx, roleID, perm.ID)
	}

	// Tidak ada perubahan
	return nil
}
