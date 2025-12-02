package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UploadReceipt godoc
// @Summary      Upload Receipt
// @Description  Upload a receipt for a booking
// @Tags         Booking
// @Accept       multipart/form-data
// @Produce      json
// @Param        receipt formData file true "Receipt File"
// @Param        booking_id formData string false "Booking ID"
// @Param        sub_booking_id formData string false "Sub Booking ID"
// @Success      200 {object} response.Response "Successfully uploaded receipt"
// @Security     BearerAuth
// @Router       /bookings/receipt [post]
func (bh *BookingHandler) UploadReceipt(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.UploadReceiptRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind form data:", err.Error())
		response.Error(c, http.StatusBadRequest, "Failed to bind form data")
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

	if err := bh.bookingUsecase.UploadReceipt(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to upload receipt", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to upload receipt")
		return
	}

	response.Success(c, nil, "Successfully uploaded receipt")
}
