package file_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024) // 32KB buffer, bisa disesuaikan
	},
}

// GetFile godoc
// @Summary      Get File
// @Description  Retrieve a file from the specified bucket and object key.
// @Tags         File
// @Accept       json
// @Produce      octet-stream
// @Param        bucket  path  string  true  "Bucket name"
// @Param        object  path  string  true  "Object key"
// @Success      200  {file}  binary  "Successfully retrieved file"
// @Router       /files/{bucket}/{object} [get]
func (fh *FileHandler) GetFile(c *gin.Context) {
	ctx := c.Request.Context()

	bucket := c.Param("bucket")
	if bucket == "" {
		logger.Error(ctx, "Bucket is required")
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	object := c.Param("object")
	if object == "" {
		logger.Error(ctx, "Object is required")
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	file, err := fh.fileUsecase.GetFiles(ctx, bucket, object)
	if err != nil {
		logger.Error(ctx, "Error getting file:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error retrieving file")
		return
	}
	defer func(file domain.StreamableObject) {
		if err := file.Close(); err != nil {
			logger.Error(ctx, "Error closing file:", err.Error())
		}
	}(file)

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.GetFilename()))
	c.Header("Content-Type", file.GetContentType())
	c.Header("Content-Length", fmt.Sprintf("%d", file.GetContentLength()))

	c.Stream(func(w io.Writer) bool {
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf)

		_, err := io.CopyBuffer(w, file, buf)
		return err == nil
	})
}
