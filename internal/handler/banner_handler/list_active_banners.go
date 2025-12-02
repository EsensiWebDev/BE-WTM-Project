package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListActiveBanners godoc
// @Summary      List Active Banners
// @Description  List all active banners
// @Tags         Banner
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]bannerdto.ActiveBanner} "Successfully retrieved active banners"
// @Router       /banners/active [get]
func (bh *BannerHandler) ListActiveBanners(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := bh.bannerUsecase.ListActiveBanners(ctx)
	if err != nil {
		logger.Error(ctx, "Error listing banners:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error retrieving banners")
		return
	}

	message := "Successfully retrieved banners"

	if resp == nil || resp.Banners == nil || len(resp.Banners) == 0 {
		message = "No active banners found"
		response.EmptyList(c, message, nil)
		return
	}

	banners := resp.Banners

	response.Success(c, banners, message)
}
