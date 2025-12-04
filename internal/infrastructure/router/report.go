package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/report_handler"

	"github.com/gin-gonic/gin"
)

func ReportRouter(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	reportHandler := report_handler.NewReportHandler(app.Usecases.ReportUsecase)
	reportGroup := routerGroup.Group("/reports", mm.Auth, mm.TimeoutSlow)
	{
		reportGroup.GET("/agent", mm.RequirePermission("report:view"), reportHandler.ReportAgent)
		reportGroup.GET("/summary", mm.RequirePermission("report:view"), reportHandler.ReportSummary)
		reportGroup.GET("/agent/detail", mm.RequirePermission("report:view"), reportHandler.ReportAgentDetail)
	}
}
