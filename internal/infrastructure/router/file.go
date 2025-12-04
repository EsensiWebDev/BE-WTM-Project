package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/file_handler"

	"github.com/gin-gonic/gin"
)

func FileRouter(app *bootstrap.Application, routerGroup *gin.RouterGroup) {
	fileHandler := file_handler.NewFileHandler(app.Usecases.FileUsecase)

	fileRouterGroup := routerGroup.Group("/files")
	{
		fileRouterGroup.GET("/:bucket/*object", fileHandler.GetFile)
	}
}
