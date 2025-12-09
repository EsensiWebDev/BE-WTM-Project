package hotel_handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sync"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

var csvBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024) // 32KB buffer
	},
}

// DownloadFormat godoc
// @Summary      Download CSV format
// @Description  Download CSV format for hotel import
// @Tags         Hotel
// @Accept       json
// @Produce      octet-stream
// @Success      200  {file}  binary  "Successfully download CSV format"
// @Router       /hotels/download-format [get]
// @Security BearerAuth
func (hh *HotelHandler) DownloadFormat(c *gin.Context) {
	ctx := c.Request.Context()

	// Create the CSV content
	var content bytes.Buffer

	// Write notes as comments
	notes := []string{
		"# NOTES:",
		"# - Column separator: ;",
		"# - Facilities: comma-separated list (e.g., WiFi,Pool,Gym)",
		"# - Nearby places: format PlaceName,Distance|PlaceName,Distance (e.g., Mall Ambassador,3|Kuningan City,2)",
		"",
	}

	for _, note := range notes {
		if _, err := content.WriteString(note + "\n"); err != nil {
			logger.Error(ctx, "Failed to write CSV notes", err.Error())
			response.Error(c, http.StatusInternalServerError, "Failed to generate CSV format")
			return
		}
	}

	// Create CSV writer
	csvWriter := csv.NewWriter(&content)
	csvWriter.Comma = ';'

	// Write header
	header := []string{
		"name", "sub_district", "district", "email", "province",
		"description", "rating", "nearby_places", "facilities",
		"tiktok", "website", "instagram",
	}

	if err := csvWriter.Write(header); err != nil {
		logger.Error(ctx, "Failed to write CSV header", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate CSV format")
		return
	}

	// Write sample data
	sampleData := []string{
		"Rose Hotel",
		"Setiabudi",
		"Jakarta",
		"rose@example.com",
		"Jakarta Province",
		"Comfortable hotel located in the city center",
		"5",
		"Mall Ambassador,3|Kuningan City,2",
		"WiFi,Pool,Gym",
		"https://tiktok.com/@rosehotel",
		"https://rosehotel.com",
		"https://instagram.com/rosehotel",
	}

	if err := csvWriter.Write(sampleData); err != nil {
		logger.Error(ctx, "Failed to write CSV sample data", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate CSV format")
		return
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		logger.Error(ctx, "Failed to flush CSV writer", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to generate CSV format")
		return
	}

	// Set response headers
	filename := "hotel_import_format.csv"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Length", fmt.Sprintf("%d", content.Len()))

	// Stream dengan buffer pool
	buf := csvBufferPool.Get().([]byte)
	defer csvBufferPool.Put(buf)

	reader := bytes.NewReader(content.Bytes())
	if _, err := io.CopyBuffer(c.Writer, reader, buf); err != nil {
		logger.Error(ctx, "Error streaming CSV file:", err.Error())
		// Jangan kirim error response setelah sebagian data terkirim
	}
}
