package user_usecase

import (
	"context"
	"errors"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (uu *UserUsecase) CreateUserByAdmin(ctx context.Context, userReq *userdto.CreateUserByAdminRequest) error {

	return uu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {

		// Generate random password
		randomString := utils.GenerateSafeRandomString(8)
		passRandom, err := utils.GeneratePassword(txCtx, randomString)
		if err != nil {
			logger.Error(txCtx, "Error to generate password", err.Error())
			return err
		}

		newUser := &entity.User{
			FullName:    userReq.FullName,
			Username:    userReq.Email,
			Password:    passRandom,
			Email:       userReq.Email,
			Phone:       userReq.Phone,
			StatusID:    constant.StatusUserActiveID,
			RoleID:      getRoleID(userReq.Role),
			KakaoTalkID: userReq.KakaoTalkID,
		}

		if newUser.RoleID == constant.RoleAgentID {
			if userReq.PromoGroupID > 0 {
				promoGroup, err := uu.promoGroupRepo.GetPromoGroupByID(txCtx, userReq.PromoGroupID)
				if err != nil {
					logger.Error(ctx, "Error to get promo group by Id", err.Error())
					return err
				}

				if promoGroup != nil && promoGroup.ID > 0 {
					newUser.PromoGroupID = &userReq.PromoGroupID
				}
			}

			if strings.TrimSpace(userReq.AgentCompany) != "" {

				agentCompany, err := uu.userRepo.CreateAgentCompany(txCtx, userReq.AgentCompany)
				if err != nil {
					logger.Error(txCtx, "Error to add agent company", err.Error())
					return err
				}

				newUser.AgentCompanyID = &agentCompany.ID
			}
		}

		userDB, err := uu.userRepo.CreateUser(txCtx, newUser)
		if err != nil {
			logger.Error(txCtx, "Error to add user", err.Error())
			return err
		}

		if userReq.PhotoSelfie != nil {
			if err := uu.uploadAndAssign(txCtx, userDB, userReq.PhotoSelfie, "selfie", &userDB.PhotoSelfie, constant.ConstPublic); err != nil {
				logger.Error(txCtx, "Error uploading selfie photo", err.Error())
				return errors.New("upload selfie photo error")
			}
		}

		if userReq.PhotoIDCard != nil {
			if err := uu.uploadAndAssign(txCtx, userDB, userReq.PhotoIDCard, "id_card", &userDB.PhotoIDCard, constant.ConstPrivate); err != nil {
				logger.Error(txCtx, "Error uploading Id card photo", err.Error())
				return err
			}
		}

		if userReq.Certificate != nil {
			if err := uu.uploadAndAssign(txCtx, userDB, userReq.Certificate, "certificate", &userDB.Certificate, constant.ConstPrivate); err != nil {
				logger.Error(txCtx, "Error uploading certificate", err.Error())
				return err
			}
		}

		if userReq.NameCard != nil {
			if err := uu.uploadAndAssign(txCtx, userDB, userReq.NameCard, "name_card", &userDB.NameCard, constant.ConstPrivate); err != nil {
				logger.Error(txCtx, "Error uploading name card photo", err.Error())
				return err
			}
		}

		_, err = uu.userRepo.UpdateUser(txCtx, userDB)
		if err != nil {
			logger.Error(txCtx, "Error updating user after upload", err.Error())
			return err
		}

		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), uu.config.DurationCtxTOSlow)
			defer cancel()
			uu.sendEmail(newCtx, userDB.FullName, userDB.Email, randomString)
		}()

		return nil
	})
}

func (uu *UserUsecase) sendEmail(ctx context.Context, name, email, tempPassword string) {
	var statusEmail = constant.EmailAccountActivated

	emailTemplate, err := uu.emailRepo.GetEmailTemplateByName(ctx, statusEmail)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for status:", statusEmail)
		return
	}

	// Inject data
	data := AccountActivatedEmailData{
		FullName:          name,
		TemporaryPassword: tempPassword,
	}

	bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
	if err != nil {
		logger.Error(ctx, "Error parsing body HTML:", err.Error())
		return
	}

	bodyText := "Please view this email in HTML format." // Optional fallback

	err = uu.emailSender.Send(ctx, email, emailTemplate.Subject, bodyHTML, bodyText)
	if err != nil {
		logger.Error(ctx, "Error sending email:", err.Error())
	}

}

type AccountActivatedEmailData struct {
	FullName          string
	TemporaryPassword string
}
