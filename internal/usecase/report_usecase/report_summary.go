package report_usecase

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (ru *ReportUsecase) ReportSummary(ctx context.Context, req *reportdto.ReportSummaryRequest) (*reportdto.ReportSummaryResponse, error) {
	// Parse dates
	dateFrom, dateTo, err := ru.parseDates(req)
	if err != nil {
		return nil, err
	}

	// Create base filter
	filterReq := filter.ReportSummaryFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	// Prepare response and synchronization
	var (
		resp     reportdto.ReportSummaryResponse
		mu       sync.Mutex
		wg       sync.WaitGroup
		firstErr error
	)

	// Goroutine 1: Booking Summary Data
	wg.Add(1)
	go func() {
		defer wg.Done()

		if ctx.Err() != nil {
			return // Context cancelled
		}

		bookingSummary, err := ru.reportRepo.ReportBookingSummary(ctx, filterReq)
		if err != nil {
			ru.setError(&mu, &firstErr, err)
			return
		}

		confirmed, cancelled, rejected := ru.processBookingSummary(bookingSummary)

		mu.Lock()
		resp.SummaryData.ConfirmedBooking = confirmed
		resp.SummaryData.CancelledBooking = cancelled
		resp.SummaryData.RejectedBooking = rejected
		defer mu.Unlock()
	}()

	// Goroutine 3: Graphic Data
	wg.Add(1)
	go func() {
		defer wg.Done()

		if ctx.Err() != nil {
			return // Context cancelled
		}

		graphicData, err := ru.reportRepo.ReportForGraph(ctx, filterReq)
		if err != nil {
			ru.setError(&mu, &firstErr, err)
			return
		}

		graphicData = ru.processGraphicData(ctx, graphicData, filterReq)

		mu.Lock()
		resp.GraphicData = graphicData
		defer mu.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for any errors
	if firstErr != nil {
		return nil, firstErr
	}

	return &resp, nil
}

// Helper functions

func (ru *ReportUsecase) processGraphicData(ctx context.Context, graphicData []entity.ReportForGraph, req filter.ReportSummaryFilter) []entity.ReportForGraph {

	if req.DateFrom == nil || req.DateTo == nil {
		return graphicData
	}

	mapData := make(map[string]int64)
	for _, datum := range graphicData {
		if datum.DateTime != nil && !datum.DateTime.IsZero() {
			normalizedDate := datum.DateTime.In(constant.AsiaJakarta).Format(time.DateOnly)
			mapData[normalizedDate] = datum.Count
		}
	}

	var current, dateTo time.Time
	current = time.Now().In(constant.AsiaJakarta)
	if req.DateFrom != nil {
		current = req.DateFrom.In(constant.AsiaJakarta)
	}

	if req.DateTo != nil {
		dateTo = req.DateTo.In(constant.AsiaJakarta).AddDate(0, 0, -1)
	}

	maxIterations := 366
	iterations := 0

	var result []entity.ReportForGraph
	for !current.After(dateTo) {
		iterations++
		if iterations > maxIterations {
			logger.Error(ctx, "Error: Maximum iteration reached")
			break
		}

		dateString := current.Format("2006-01-02")
		result = append(result, entity.ReportForGraph{
			Date:  current.Format("2006-01-02"),
			Count: mapData[dateString],
		})
		current = current.AddDate(0, 0, 1)
	}

	return result
}

func (ru *ReportUsecase) parseDates(req *reportdto.ReportSummaryRequest) (*time.Time, *time.Time, error) {
	var dateFrom, dateTo *time.Time
	// Parse DateFrom
	if req.DateFrom != "" {
		dateFromDt, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid date_from: %s", err.Error())
		}
		dateFrom = &dateFromDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, constant.AsiaJakarta)
		dateFrom = &startOfMonth
	}

	// Parse DateTo
	if req.DateTo != "" {
		dateToDt, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid date_to: %s", err.Error())
		}
		dateToDt = dateToDt.AddDate(0, 0, 1)
		dateTo = &dateToDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, constant.AsiaJakarta)
		dateTo = &endOfMonth
	}

	return dateFrom, dateTo, nil
}

func (ru *ReportUsecase) processBookingSummary(bookingSummary []entity.MonthlyBookingSummary) (reportdto.DataTotalWithPercentage, reportdto.DataTotalWithPercentage, reportdto.DataTotalWithPercentage) {
	summaryDataBooking := reportdto.DataTotalWithPercentage{}
	summaryDataCancel := reportdto.DataTotalWithPercentage{}
	summaryDataReject := reportdto.DataTotalWithPercentage{}

	if len(bookingSummary) == 2 {
		summaryDataBooking.Count = bookingSummary[0].ConfirmedBooking
		summaryDataCancel.Count = bookingSummary[0].CancelledBooking
		summaryDataReject.Count = bookingSummary[0].RejectedBooking

		summaryDataBooking = ru.calculateBookingPercentage(bookingSummary[0].ConfirmedBooking, bookingSummary[1].ConfirmedBooking, summaryDataBooking, "Contribution for this period")
		summaryDataCancel = ru.calculateBookingPercentage(bookingSummary[0].CancelledBooking, bookingSummary[1].CancelledBooking, summaryDataCancel, "Contribution for this period")
		summaryDataReject = ru.calculateBookingPercentage(bookingSummary[0].RejectedBooking, bookingSummary[1].RejectedBooking, summaryDataReject, "Contribution for this period")
	}

	return summaryDataBooking, summaryDataCancel, summaryDataReject
}

func (ru *ReportUsecase) calculateBookingPercentage(current, previous int64, data reportdto.DataTotalWithPercentage, message string) reportdto.DataTotalWithPercentage {
	if previous == 0 {
		data.Percent = 0
		data.Message = "No data for comparison period"
	} else {
		percent := (float64(current) / float64(previous)) * 100
		data.Percent = math.Round(percent*100) / 100
		data.Message = message
	}
	return data
}

func (ru *ReportUsecase) setError(mu *sync.Mutex, firstErr *error, err error) {
	mu.Lock()
	defer mu.Unlock()
	if *firstErr == nil {
		*firstErr = err
	}
}
