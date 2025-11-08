package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// Register handles user registration requests.
// @Summary User Registration
// @Description Registers a new user with username, email, and etc.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param full_name formData string true "Full Name"
// @Param agent_company formData string false "Agent Company"
// @Param email formData string true "Email"
// @Param phone formData string true "Phone"
// @Param username formData string true "Username"
// @Param kakao_talk_id formData string true "Kakao Talk Id"
// @Param password formData string true "Password"
// @Param photo_selfie formData file true "File Selfie""
// @Param photo_id_card formData file true "File Id Card"
// @Param certificate formData file false "Certificate"
// @Param name_card formData file true "Name Card"
// @Success 200 {object} response.Response "Successfully registered"
// @Router /register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.RegisterRequest
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

	if err := uh.userUsecase.Register(ctx, &req); err != nil {
		logger.Error(ctx, "Error registering user:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Registration failed")
		return
	}

	response.Success(c, nil, "Successfully registered")
}
