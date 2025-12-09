package report_handler

import (
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ReportAgentDetail godoc
// @Summary      Generate Agent Report Detail
// @Description  Generate a detailed report for agent bookings with pagination
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param hotel_id query int false "Filter by Hotel Id"
// @Param agent_id query int false "Filter by Agent Id"
// @Param date_from query string false "Start date for the report in YYYY-MM-DD format"
// @Param date_to query string false "End date for the report in YYYY-MM-DD format"
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.ReportAgentDetail} "Successfully generated detailed report for agent bookings"
// @Security BearerAuth
// @Router       /reports/agent/detail [get]
func (rh *ReportHandler) ReportAgentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	var req reportdto.ReportAgentDetailRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding ReportAgentDetail request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := rh.reportUsecase.ReportAgentDetail(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error generating ReportAgentDetail", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate report detail")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully generated detailed report for agent bookings"
	var datas []entity.ReportAgentDetail
	if resp != nil {
		datas = resp.ReportAgentDetailData
		if len(datas) == 0 {
			message = "No data found for the given criteria"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, datas, message, pagination)
}
