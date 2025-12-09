package report_handler

import (
	"net/http"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ReportSummary godoc
// @Summary      Generate Report Summary
// @Description  Generate a summary report of bookings with optional filters
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param date_from query string false "Start date for the report in YYYY-MM-DD format"
// @Param date_to query string false "End date for the report in YYYY-MM-DD format"
// @Success 200 {object} response.ResponseWithData{data=reportdto.ReportSummaryResponse} "Successfully generated summary report"
// @Security BearerAuth
// @Router       /reports/summary [get]
func (rh *ReportHandler) ReportSummary(c *gin.Context) {
	ctx := c.Request.Context()

	var req reportdto.ReportSummaryRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding ReportSummary request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := rh.reportUsecase.ReportSummary(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error generating ReportSummary", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate report summary")
		return
	}

	message := "Successfully generated summary report"
	if resp == nil {
		message = "No data found for the given criteria"
		resp = &reportdto.ReportSummaryResponse{}
	}

	response.Success(c, resp, message)
}
