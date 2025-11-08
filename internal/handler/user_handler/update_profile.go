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

// UpdateProfile updates user profile information
// @Summary Update user profile
// @Description Update user profile information
// @Tags Profile
// @Produce json
// @Router /profile [put]
// @Param request body userdto.UpdateProfileRequest true "Update Profile Request"
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated profile"
func (uh *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Bind the request data to the UpdateProfileRequest struct
	var req userdto.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the request data
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

	// Call the use case to update the profile
	if err := uh.userUsecase.UpdateProfile(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating profile:", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed updating profile: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully updated profile")
}
