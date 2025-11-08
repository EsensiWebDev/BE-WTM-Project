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

// UpdateFile updates user's profile file
// @Summary Update user profile file
// @Description Update user profile file
// @Tags Profile
// @Produce json
// @Accept multipart/form-data
// @Param file_type formData string true "File Type" Enums(photo, certificate, name_card)
// @Param file formData file true "File"
// @Router /profile/file [put]
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated profile file"
func (uh *UserHandler) UpdateFile(c *gin.Context) {
	ctx := c.Request.Context()

	// Bind the request data to the UpdateFileRequest struct
	var req userdto.UpdateFileRequest
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

	// Call the usecase to update the user's profile photo
	if err := uh.userUsecase.UpdateFile(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating user photo:", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Error updating user photo: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully updated profile file")
}
