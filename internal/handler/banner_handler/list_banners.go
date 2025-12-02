package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListBanners godoc
// @Summary      List Banners
// @Description  List all banners with pagination and search functionality
// @Tags         Banner
// @Accept       json
// @Produce      json
// @Param        page  query  int  false  "Page number for pagination"
// @Param        limit query  int  false  "Number of items per page"
// @Param        search query  string  false  "Search term for banner title"
// @Param        sort    query  string  false  "Sort order for banners"
// @Param        dir      query  string  false  "Sort direction for banners"
// @Param        is_active query  bool  false  "Filter by active status of banners"
// @Success      200  {object}  response.ResponseWithPagination{data=[]bannerdto.BannerData} "Successfully retrieved banners"
// @Security      BearerAuth
// @Router       /banners [get]
func (bh *BannerHandler) ListBanners(c *gin.Context) {
	ctx := c.Request.Context()

	var req bannerdto.ListBannerRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding query parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	resp, err := bh.bannerUsecase.ListBanners(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error listing banners:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error retrieving banners")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved banners"

	if resp == nil || resp.Banners == nil || len(resp.Banners) == 0 {
		message = "No banners found"
		response.EmptyList(c, message, pagination)
		return
	}

	banners := resp.Banners
	pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))

	response.SuccessWithPagination(c, banners, message, pagination)
}
