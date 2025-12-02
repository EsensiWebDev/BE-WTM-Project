package report_usecase

import (
	"context"
	"fmt"
	"sync"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
)

func (ru *ReportUsecase) ReportSummary(ctx context.Context, req *reportdto.ReportRequest) (*reportdto.ReportSummaryResponse, error) {
	// Parse dates
	dateFrom, dateTo, isRangeDate, err := ru.parseDates(req)
	if err != nil {
		return nil, err
	}

	// Create base filter
	filterReq := ru.createBaseFilter(dateFrom, dateTo, req, isRangeDate)

	// Create independent filters for each goroutine
	filterBooking := ru.createFilterCopy(filterReq)
	filterNewCustomer := ru.createFilterCopy(filterReq)
	filterGraphic := ru.createFilterCopy(filterReq)

	// Adjust dates for comparison periods
	if !isRangeDate {
		ru.adjustComparisonDates(&filterBooking, &filterNewCustomer, dateFrom)
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

		bookingSummary, err := ru.reportRepo.ReportBookingSummary(ctx, filterBooking)
		if err != nil {
			ru.setError(&mu, &firstErr, err)
			return
		}

		confirmed, cancelled := ru.processBookingSummary(bookingSummary, isRangeDate)

		mu.Lock()
		defer mu.Unlock()
		resp.SummaryData.ConfirmedBooking = confirmed
		resp.SummaryData.CancellationBooking = cancelled
	}()

	// Goroutine 2: New Customer Summary Data
	wg.Add(1)
	go func() {
		defer wg.Done()

		if ctx.Err() != nil {
			return // Context cancelled
		}

		newCustomer, err := ru.reportRepo.ReportNewAgentSummary(ctx, filterNewCustomer)
		if err != nil {
			ru.setError(&mu, &firstErr, err)
			return
		}

		newUser := ru.processNewCustomerSummary(newCustomer, isRangeDate)

		mu.Lock()
		defer mu.Unlock()
		resp.SummaryData.NewCustomer = newUser
	}()

	// Goroutine 3: Graphic Data
	wg.Add(1)
	go func() {
		defer wg.Done()

		if ctx.Err() != nil {
			return // Context cancelled
		}

		graphicData, err := ru.reportRepo.ReportForGraph(ctx, filterGraphic)
		if err != nil {
			ru.setError(&mu, &firstErr, err)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resp.GraphicData = graphicData
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

func (ru *ReportUsecase) parseDates(req *reportdto.ReportRequest) (*time.Time, *time.Time, bool, error) {
	var dateFrom, dateTo *time.Time
	var isDateFromSet, isDateToSet bool

	// Parse DateFrom
	if req.DateFrom != "" {
		dateFromDt, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, nil, false, fmt.Errorf("invalid date_from: %s", err.Error())
		}
		isDateFromSet = true
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
			return nil, nil, false, fmt.Errorf("invalid date_to: %s", err.Error())
		}
		isDateToSet = true
		dateTo = &dateToDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, constant.AsiaJakarta)
		dateTo = &endOfMonth
	}

	return dateFrom, dateTo, (isDateFromSet || isDateToSet), nil
}

func (ru *ReportUsecase) createBaseFilter(dateFrom, dateTo *time.Time, req *reportdto.ReportRequest, isRangeDate bool) filter.ReportFilter {
	filterReq := filter.ReportFilter{
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		IsRangeDate: isRangeDate,
	}

	return filterReq
}

func (ru *ReportUsecase) adjustComparisonDates(filterBooking, filterNewCustomer *filter.ReportFilter, dateFrom *time.Time) {
	dateFromBooking := dateFrom.AddDate(0, 0, -30)
	filterBooking.DateFrom = &dateFromBooking

	dateFromNewCustomer := dateFrom.AddDate(0, 0, -30)
	filterNewCustomer.DateFrom = &dateFromNewCustomer
}

func (ru *ReportUsecase) processBookingSummary(bookingSummary []entity.MonthlyBookingSummary, isRangeDate bool) (reportdto.DataTotalWithPercentage, reportdto.DataTotalWithPercentage) {
	summaryDataBooking := reportdto.DataTotalWithPercentage{}
	summaryDataCancel := reportdto.DataTotalWithPercentage{}

	if isRangeDate {
		// Range date logic
		if len(bookingSummary) == 2 {
			summaryDataBooking.Count = bookingSummary[1].ConfirmedBooking
			summaryDataCancel.Count = bookingSummary[1].CancellationBooking

			summaryDataBooking = ru.calculateBookingPercentage(bookingSummary[0].ConfirmedBooking, bookingSummary[1].ConfirmedBooking, summaryDataBooking, "Contribution for this period")
			summaryDataCancel = ru.calculateBookingPercentage(bookingSummary[0].CancellationBooking, bookingSummary[1].CancellationBooking, summaryDataCancel, "Contribution for this period")
		}
	} else {
		// Non-range date logic
		switch len(bookingSummary) {
		case 1:
			summaryDataBooking.Count = bookingSummary[0].ConfirmedBooking
			summaryDataCancel.Count = bookingSummary[0].CancellationBooking
			summaryDataBooking.Percent = 0
			summaryDataBooking.Message = "No comparison data available"
			summaryDataCancel.Percent = 0
			summaryDataCancel.Message = "No comparison data available"

		case 2:
			summaryDataBooking.Count = bookingSummary[1].ConfirmedBooking
			summaryDataCancel.Count = bookingSummary[1].CancellationBooking

			summaryDataBooking = ru.calculateBookingPercentage(bookingSummary[0].ConfirmedBooking, bookingSummary[1].ConfirmedBooking, summaryDataBooking, "Comparison with last period")
			summaryDataCancel = ru.calculateBookingPercentage(bookingSummary[0].CancellationBooking, bookingSummary[1].CancellationBooking, summaryDataCancel, "Comparison with last period")
		}
	}

	return summaryDataBooking, summaryDataCancel
}

func (ru *ReportUsecase) calculateBookingPercentage(previous, current int64, data reportdto.DataTotalWithPercentage, message string) reportdto.DataTotalWithPercentage {
	if previous == 0 {
		data.Percent = 0
		data.Message = "No data for comparison period"
	} else {
		data.Percent = (float64(current) / float64(previous)) * 100
		data.Message = message
	}
	return data
}

func (ru *ReportUsecase) processNewCustomerSummary(newCustomer []entity.MonthlyNewAgentSummary, isRangeDate bool) reportdto.DataTotalWithPercentage {
	summaryDataNewUser := reportdto.DataTotalWithPercentage{}

	if isRangeDate {
		// Range date logic
		if len(newCustomer) == 2 {
			summaryDataNewUser.Count = newCustomer[0].NewAgent
			summaryDataNewUser = ru.calculateNewUserPercentage(newCustomer[0].NewAgent, newCustomer[1].NewAgent, summaryDataNewUser, "Contribution for this period")
		}
	} else {
		// Non-range date logic
		switch len(newCustomer) {
		case 1:
			summaryDataNewUser.Count = newCustomer[0].NewAgent
			summaryDataNewUser.Percent = 0
			summaryDataNewUser.Message = "No comparison data available"

		case 2:
			summaryDataNewUser.Count = newCustomer[1].NewAgent
			summaryDataNewUser = ru.calculateNewUserPercentage(newCustomer[0].NewAgent, newCustomer[1].NewAgent, summaryDataNewUser, "Comparison with last period")
		}
	}

	return summaryDataNewUser
}

func (ru *ReportUsecase) calculateNewUserPercentage(previous, current int64, data reportdto.DataTotalWithPercentage, message string) reportdto.DataTotalWithPercentage {
	if previous == 0 {
		data.Percent = 0
		data.Message = "No data for comparison period"
	} else {
		data.Percent = (float64(current) / float64(previous)) * 100
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

func (ru *ReportUsecase) createFilterCopy(filter filter.ReportFilter) filter.ReportFilter {
	copyFilter := filter

	// Deep copy pointers
	if filter.DateFrom != nil {
		dateFromCopy := *filter.DateFrom
		copyFilter.DateFrom = &dateFromCopy
	}
	if filter.DateTo != nil {
		dateToCopy := *filter.DateTo
		copyFilter.DateTo = &dateToCopy
	}

	return copyFilter
}
