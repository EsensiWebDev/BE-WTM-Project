package report_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ReportAgent godoc
// @Summary      Generate Agent Report
// @Description  Generate a report for agent bookings within a specified date range and filters
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search term"
// @Param hotel_id query []int false "Filter by Hotel Id" collectionFormat(multi)
// @Param agent_company_id query []int false "Filter by Agent Company Id" collectionFormat(multi)
// @Param date_from query string false "Start date for the report in YYYY-MM-DD format"
// @Param date_to query string false "End date for the report in YYYY-MM-DD format"
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.ReportAgentBooking} "Successfully generated report for agent bookings"
// @Security BearerAuth
// @Router       /reports/agent [get]
func (rh *ReportHandler) ReportAgent(c *gin.Context) {
	ctx := c.Request.Context()

	var req reportdto.ReportRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding ReportAgent request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := rh.reportUsecase.ReportAgent(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error generating ReportAgent", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate report")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully generated report for agent bookings"

	var datas []entity.ReportAgentBooking
	if resp != nil {
		datas = resp.Data
		if len(datas) == 0 {
			message = "No data found for the given criteria"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, datas, message, pagination)
}
