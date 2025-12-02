package seed

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
)

func (s *Seed) SeedUser() {

	var countS, countR, countP, countU int64

	s.db.Model(&model.StatusUser{}).Count(&countS)
	s.db.Model(&model.Role{}).Count(&countR)
	s.db.Model(&model.Permission{}).Count(&countP)
	s.db.Model(&model.User{}).Count(&countU)

	statuses := []model.StatusUser{
		{ID: constant.StatusUserWaitingApprovalID, Status: constant.StatusUserWaitingApproval},
		{Status: constant.StatusUserActive},
		{Status: constant.StatusUserReject},
		{Status: constant.StatusUserInactive},
	}

	if countS == 0 {
		if err := s.db.Create(&statuses).Error; err != nil {
			log.Fatalf("Failed to seed users: %s", err.Error())
		}
		log.Println("Seeding users completed")
	}

	if countR == 0 || countP == 0 || countU == 0 {

		s.db.Unscoped().Where("1 = 1").Delete(&model.User{})
		s.db.Unscoped().Where("1 = 1").Delete(&model.Permission{})
		s.db.Unscoped().Where("1 = 1").Delete(&model.Role{})

		perms := []model.Permission{
			{Permission: "account:view", Page: "account", Action: "view"},
			{Permission: "account:create", Page: "account", Action: "create"},
			{Permission: "account:edit", Page: "account", Action: "edit"},
			{Permission: "account:delete", Page: "account", Action: "delete"},
			{Permission: "hotel:view", Page: "hotel", Action: "view"},
			{Permission: "hotel:create", Page: "hotel", Action: "create"},
			{Permission: "hotel:edit", Page: "hotel", Action: "edit"},
			{Permission: "hotel:delete", Page: "hotel", Action: "delete"},
			{Permission: "promo:view", Page: "promo", Action: "view"},
			{Permission: "promo:create", Page: "promo", Action: "create"},
			{Permission: "promo:edit", Page: "promo", Action: "edit"},
			{Permission: "promo:delete", Page: "promo", Action: "delete"},
			{Permission: "promo-group:view", Page: "promo-group", Action: "view"},
			{Permission: "promo-group:create", Page: "promo-group", Action: "create"},
			{Permission: "promo-group:edit", Page: "promo-group", Action: "edit"},
			{Permission: "promo-group:delete", Page: "promo-group", Action: "delete"},
			{Permission: "report:view", Page: "report", Action: "view"},
			{Permission: "report:create", Page: "report", Action: "create"},
			{Permission: "report:edit", Page: "report", Action: "edit"},
			{Permission: "report:delete", Page: "report", Action: "delete"},
			{Permission: "booking:view", Page: "booking", Action: "view"},
			{Permission: "booking:create", Page: "booking", Action: "create"},
			{Permission: "booking:edit", Page: "booking", Action: "edit"},
			{Permission: "booking:delete", Page: "booking", Action: "delete"},
		}
		if err := s.db.Create(&perms).Error; err != nil {
			log.Fatalf("Failed to seed roles: %s", err.Error())
		}

		permMap := map[string][]model.Permission{}
		for _, p := range perms {
			permMap[p.Action] = append(permMap[p.Action], p)
		}

		collectPerms := func(actions ...string) []model.Permission {
			var result []model.Permission
			for _, act := range actions {
				result = append(result, permMap[act]...)
			}
			return result
		}

		superAdmin := model.Role{Role: "Super Admin"}
		admin := model.Role{Role: "Admin", Permissions: collectPerms("view", "create", "edit")}
		agent := model.Role{Role: "Agent", Permissions: collectPerms("view")}
		support := model.Role{Role: "Support", Permissions: collectPerms("view", "edit")}

		roles := []model.Role{superAdmin, admin, agent, support}
		if err := s.db.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed roles: %s", err.Error())
		}
		countR = int64(len(roles))

		defaultPassword, err := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %s", err.Error())
		}

		var statusesMap = map[string]model.StatusUser{}
		for _, status := range statuses {
			statusesMap[status.Status] = status
		}

		var rolesMap = map[string]model.Role{}
		for _, role := range roles {
			rolesMap[role.Role] = role
		}

		//Mapping user -> role
		superAdminUser := model.User{
			FullName: "Super Admin",
			Username: "superadmin",
			Password: string(defaultPassword),
			StatusID: statusesMap[constant.StatusUserActive].ID,
			RoleID:   rolesMap["Super Admin"].ID,
		}

		users := []model.User{superAdminUser}
		if err := s.db.Create(&users).Error; err != nil {
			log.Fatalf("Failed to seed users: %s", err.Error())
		}

		log.Println("Seeding completed")
	}

	if countR < 4 {
		s.db.Unscoped().Where("1 = 1").Delete(&model.Role{})

		roles := []model.Role{
			{Role: "Super Admin"},
			{Role: "Admin"},
			{Role: "Agent"},
			{Role: "Support"},
		}

		if err := s.db.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed roles: %s", err.Error())
		}

		log.Println("Seeding roles completed")
	}
}
