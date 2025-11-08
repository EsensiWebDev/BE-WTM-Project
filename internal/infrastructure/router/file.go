package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/file_handler"
)

func FileRouter(app *bootstrap.Application, routerGroup *gin.RouterGroup) {
	fileHandler := file_handler.NewFileHandler(app.Usecases.FileUsecase)

	fileRouterGroup := routerGroup.Group("/files")
	{
		fileRouterGroup.GET("/:bucket/:object", fileHandler.GetFile)
	}
}
