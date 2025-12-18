package user_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateUserByAdmin godoc
// @Summary Update a user by admin
// @Description Update a user with the provided details
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData int true "User ID"
// @Param full_name formData string true "Full Name"
// @Param role formData string true "Role (e.g., admin, agent, customer)"
// @Param email formData string true "Email"
// @Param username formData string true "Username"
// @Param phone formData string true "Phone"
// @Param currency formData string false "Currency (e.g., IDR, USD)"
// @Param kakao_talk_id formData string false "Kakao Talk Id"
// @Param promo_group_id formData int false "Promo Group ID (required if role is agent)"
// @Param agent_company formData string false "Agent Company (required if role is agent)"
// @Param certificate formData file false "Certificate (optional)"
// @Param photo_selfie formData file false "File Selfie"
// @Param photo_id_card formData file false "File Id Card"
// @Param name_card formData file false "Name Card"
// @Param is_active formData bool true "Is Active"
// @Router /users [put]
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated user"
func (uh *UserHandler) UpdateUserByAdmin(c *gin.Context) {
	var req userdto.UpdateUserByAdminRequest
	ctx := c.Request.Context()

	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Error validating request:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}

		// fallback: unknown validation error
		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := uh.userUsecase.UpdateUserByAdmin(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating user", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully updated user")
}
