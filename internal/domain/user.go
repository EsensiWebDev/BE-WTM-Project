package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/repository/filter"
)

type UserUsecase interface {
	Register(ctx context.Context, user *userdto.RegisterRequest) error
	CreateUserByAdmin(ctx context.Context, user *userdto.CreateUserByAdminRequest) error
	Profile(ctx context.Context) (*userdto.ProfileResponse, error)
	UpdateSetting(ctx context.Context, user *userdto.UpdateSettingRequest) error
	UpdateProfile(ctx context.Context, user *userdto.UpdateProfileRequest) error
	UpdateFile(ctx context.Context, user *userdto.UpdateFileRequest) error
	ListAgentCompanies(ctx context.Context, request *userdto.ListAgentCompaniesRequest) (*userdto.ListAgentCompaniesResponse, int64, error)
	ListUsersByAgentCompany(ctx context.Context, request *userdto.ListUsersByAgentCompanyRequest) (*userdto.ListUsersByAgentCompanyResponse, int64, error)
	ListUsersByRole(ctx context.Context, request *userdto.ListUsersByRoleRequest) (*userdto.ListUsersByRoleResponse, int64, error)
	ListUsers(ctx context.Context, req *userdto.ListUsersRequest) (*userdto.ListUsersResponse, error)
	UpdateUserByAdmin(ctx context.Context, req *userdto.UpdateUserByAdminRequest) error
	ListRoleAccess(ctx context.Context) ([]userdto.ListRoleAccessResponse, error)
	UpdateRoleAccess(ctx context.Context, req *userdto.UpdateRoleAccessRequest) error
	ListStatusUsers(ctx context.Context, req *userdto.ListStatusUsersRequest) (*userdto.ListStatusUsersResponse, int64, error)
	UpdateStatusUser(ctx context.Context, req *userdto.UpdateStatusUserRequest) error
}

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uint) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	CreateAgentCompany(ctx context.Context, agentCompany string) (*entity.AgentCompany, error)
	GetAgentCompanies(ctx context.Context, search string, limit, page int) ([]entity.AgentCompany, int64, error)
	GetUsersByAgentCompany(ctx context.Context, agentCompanyID uint, search string, limit, page int) ([]entity.User, int64, error)
	GetUserByRole(ctx context.Context, roleID uint, search string, limit, page int) ([]entity.User, int64, error)
	GetUsers(ctx context.Context, filter filter.UserFilter) ([]entity.User, int64, error)
	BulkUpdatePromoGroupMember(ctx context.Context, memberIDs []uint, promoGroupID uint) error
	GetAllRolesWithPermissions(ctx context.Context) ([]entity.Role, error)
	GetAllPermissions(ctx context.Context) ([]entity.Permission, error)
	GetPermissionByPageAction(ctx context.Context, page, action string) (*entity.Permission, error)
	AddRolePermission(ctx context.Context, roleID, permissionID uint) error
	RemoveRolePermission(ctx context.Context, roleID, permissionID uint) error
	HasRolePermission(ctx context.Context, roleID, permissionID uint) (bool, error)
	GetStatusUsers(ctx context.Context, filter *filter.DefaultFilter) ([]entity.StatusUser, int64, error)
	UpdateStatusUser(ctx context.Context, id uint, status uint) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
}
