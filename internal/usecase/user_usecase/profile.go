package user_usecase

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"wtm-backend/internal/domain/entity"
	dtouser "wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) Profile(ctx context.Context) (*dtouser.ProfileResponse, error) {

	userCtx, err := uu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error generating user from context", err.Error())
		return nil, fmt.Errorf("failed to generate user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "User context is nil")
		return nil, fmt.Errorf("user not found in context")
	}

	user, err := uu.userRepo.GetUserByID(ctx, userCtx.ID)
	if err != nil {
		logger.Error(ctx, "Error getting user by Id", err.Error())
		return nil, err
	}

	if user == nil {
		logger.Error(ctx, "User not found in database")
		return nil, fmt.Errorf("user not found in database")
	}

	var photoProfile string
	if strings.TrimSpace(user.PhotoSelfie) != "" {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPublic)
		photoProfile, err = uu.fileStorage.GetFile(ctx, bucketName, user.PhotoSelfie)
		if err != nil {
			logger.Error(ctx, "Error getting user profile photo", err.Error())
			return nil, fmt.Errorf("failed to get user profile photo: %s", err.Error())
		}
	}

	var certificateURL string
	if strings.TrimSpace(user.Certificate) != "" {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPrivate)
		certificateURL, err = uu.fileStorage.GetFile(ctx, bucketName, user.Certificate)
		if err != nil {
			logger.Error(ctx, "Error getting user certificate", err.Error())
			return nil, fmt.Errorf("failed to get user certificate: %s", err.Error())
		}
	}

	var nameCardURL string
	if strings.TrimSpace(user.NameCard) != "" {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPrivate)
		nameCardURL, err = uu.fileStorage.GetFile(ctx, bucketName, user.NameCard)
		if err != nil {
			logger.Error(ctx, "Error getting user name card", err.Error())
			return nil, fmt.Errorf("failed to get user name card: %s", err.Error())
		}
	}

	profileResponse := &dtouser.ProfileResponse{
		ID:                  user.ID,
		FullName:            user.FullName,
		Username:            user.Username,
		Email:               user.Email,
		Phone:               user.Phone,
		Password:            strings.Repeat("*", 8),
		PhotoProfile:        photoProfile,
		Certificate:         certificateURL,
		NameCard:            nameCardURL,
		KakaoTalkID:         user.KakaoTalkID,
		Status:              user.StatusName,
		NotificationSetting: summarizeNotificationSettings(user.UserNotificationSettings),
	}

	return profileResponse, nil
}

func summarizeNotificationSettings(settings []entity.UserNotificationSetting) []dtouser.NotificationSetting {
	allTypes := []string{constant.ConstBooking, constant.ConstReject}
	grouped := map[string][]string{}
	channelSet := map[string]bool{}

	for _, s := range settings {
		channelSet[s.Channel] = true
		if s.IsEnabled {
			grouped[s.Channel] = append(grouped[s.Channel], s.Type)
		}
	}

	var result []dtouser.NotificationSetting
	for channel := range channelSet {
		enabledTypes := grouped[channel]
		isEnable := len(enabledTypes) > 0

		var typeStr string
		sort.Strings(enabledTypes)
		if !isEnable {
			typeStr = ""
		} else if reflect.DeepEqual(enabledTypes, allTypes) {
			typeStr = constant.ConstAll
		} else {
			typeStr = strings.Join(enabledTypes, ",")
		}

		result = append(result, dtouser.NotificationSetting{
			Channel:  channel,
			Type:     typeStr,
			IsEnable: isEnable,
		})
	}

	return result
}
