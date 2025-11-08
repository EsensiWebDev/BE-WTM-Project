package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/report_handler"
)

func ReportRouter(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	reportHandler := report_handler.NewReportHandler(app.Usecases.ReportUsecase)
	reportGroup := routerGroup.Group("/reports", middlewareMap.Auth, middlewareMap.TimeoutSlow)
	{
		reportGroup.GET("/agent", reportHandler.ReportAgent)
		reportGroup.GET("/summary", reportHandler.ReportSummary)
		reportGroup.GET("/agent/detail", reportHandler.ReportAgentDetail)
	}
}
