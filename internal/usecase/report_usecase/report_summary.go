package report_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
)

func (ru *ReportUsecase) ReportSummary(ctx context.Context, req *reportdto.ReportRequest) (*reportdto.ReportSummaryResponse, error) {
	var dateFrom, dateTo *time.Time
	var isDateFromSet, isDateToSet bool
	// Parse dates
	if req.DateFrom != "" {
		dateFromDt, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from: %s", err.Error())
		}
		isDateFromSet = true
		dateFrom = &dateFromDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, constant.AsiaJakarta)
		dateFrom = &startOfMonth
	}

	if req.DateTo != "" {
		dateToDt, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to: %s", err.Error())
		}
		isDateToSet = true
		dateTo = &dateToDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, constant.AsiaJakarta)
		dateTo = &endOfMonth
	}

	filterReq := filter.ReportFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	if req.HotelID > 0 {
		filterReq.HotelID = &req.HotelID
	}
	if req.AgentCompanyID > 0 {
		filterReq.AgentCompanyID = &req.AgentCompanyID
	}

	if isDateFromSet || isDateToSet {
		filterReq.IsRangeDate = true
	}

	// Channel sink
	type result struct {
		summaryDataBooking reportdto.DataTotalWithPercentage
		summaryDataCancel  reportdto.DataTotalWithPercentage
		summaryDataNewUser reportdto.DataTotalWithPercentage
		graphicData        []entity.ReportForGraph
		err                error
	}
	ch := make(chan result, 3)

	// Goroutine 1: SummaryData
	go func() {

		if !filterReq.IsRangeDate {
			dateFromDt := dateFrom.AddDate(0, 0, -30)
			filterReq.DateFrom = &dateFromDt
		}

		bookingSummary, err := ru.reportRepo.ReportBookingSummary(ctx, filterReq)
		if err != nil {
			ch <- result{err: err}
			return
		}

		summaryDataBooking := reportdto.DataTotalWithPercentage{}
		summaryDataCancel := reportdto.DataTotalWithPercentage{}

		if filterReq.IsRangeDate {
			// this month and total
			if len(bookingSummary) == 2 {

				summaryDataBooking.Count = bookingSummary[1].ConfirmedBooking
				summaryDataCancel.Count = bookingSummary[1].CancellationBooking

				if bookingSummary[0].ConfirmedBooking == 0 {
					summaryDataBooking.Percent = 0
					summaryDataBooking.Message = "No confirmed booking for this period"
				} else {
					summaryDataBooking.Percent = (float64(bookingSummary[0].ConfirmedBooking) / float64(bookingSummary[1].ConfirmedBooking)) * 100
					summaryDataBooking.Message = "Contribution for this period"
				}

				if bookingSummary[0].CancellationBooking == 0 {
					summaryDataCancel.Percent = 0
					summaryDataCancel.Message = "No cancellation booking for this period"
				} else {
					summaryDataCancel.Percent = (float64(bookingSummary[0].CancellationBooking) / float64(bookingSummary[1].CancellationBooking)) * 100
					summaryDataCancel.Message = "Contribution for this period"
				}
			}
		} else {
			// this period only
			if len(bookingSummary) == 1 {

				summaryDataBooking.Count = bookingSummary[0].ConfirmedBooking
				summaryDataCancel.Count = bookingSummary[0].CancellationBooking
				summaryDataBooking.Percent = 0
				summaryDataBooking.Message = "No comparison data available"
				summaryDataCancel.Percent = 0
				summaryDataCancel.Message = "No comparison data available"
			}

			// this period and last period
			if len(bookingSummary) == 2 {

				summaryDataBooking.Count = bookingSummary[1].ConfirmedBooking
				summaryDataCancel.Count = bookingSummary[1].CancellationBooking

				if bookingSummary[0].ConfirmedBooking == 0 {
					summaryDataBooking.Percent = 0
					summaryDataBooking.Message = "No confirmed booking for this period"
				} else {
					summaryDataBooking.Percent = (float64(bookingSummary[1].ConfirmedBooking) / float64(bookingSummary[0].ConfirmedBooking)) * 100
					summaryDataBooking.Message = "Comparison with last period"
				}

				if bookingSummary[0].CancellationBooking == 0 {
					summaryDataCancel.Percent = 0
					summaryDataCancel.Message = "No cancellation booking for last period"
				} else {
					summaryDataCancel.Percent = (float64(bookingSummary[1].CancellationBooking) / float64(bookingSummary[0].CancellationBooking)) * 100
					summaryDataCancel.Message = "Comparison with last period"
				}
			}
		}

		ch <- result{summaryDataBooking: summaryDataBooking, summaryDataCancel: summaryDataCancel}
	}()

	// Goroutine 2: New Customer Summary
	go func() {
		if !filterReq.IsRangeDate {
			dateFromDt := dateFrom.AddDate(0, 0, -30)
			filterReq.DateFrom = &dateFromDt
		}

		newCustomer, err := ru.reportRepo.ReportNewAgentSummary(ctx, filterReq)
		if err != nil {
			ch <- result{err: err}
			return
		}

		summaryDataNewUser := reportdto.DataTotalWithPercentage{}

		if filterReq.IsRangeDate {
			// this month and total
			if len(newCustomer) == 2 {
				summaryDataNewUser.Count = newCustomer[0].NewAgent

				if newCustomer[0].NewAgent == 0 {
					summaryDataNewUser.Percent = 0
					summaryDataNewUser.Message = "No new customer for this period"
				} else {
					summaryDataNewUser.Percent = (float64(newCustomer[0].NewAgent) / float64(newCustomer[1].NewAgent)) * 100
					summaryDataNewUser.Message = "Contribution for this period"
				}

			}
		} else {
			// this period only
			if len(newCustomer) == 1 {
				summaryDataNewUser.Count = newCustomer[0].NewAgent
				summaryDataNewUser.Percent = 0
				summaryDataNewUser.Message = "No comparison data available"
			}

			// this period and last period
			if len(newCustomer) == 2 {
				summaryDataNewUser.Count = newCustomer[1].NewAgent

				if newCustomer[0].NewAgent == 0 {
					summaryDataNewUser.Percent = 0
					summaryDataNewUser.Message = "No new customer for this period"
				} else {
					summaryDataNewUser.Percent = (float64(newCustomer[1].NewAgent) / float64(newCustomer[0].NewAgent)) * 100
					summaryDataNewUser.Message = "Comparison with last period"
				}

			}
		}

		ch <- result{summaryDataNewUser: summaryDataNewUser}

	}()

	// Goroutine 3: DataGrafik
	go func() {
		graphicData, err := ru.reportRepo.ReportForGraph(ctx, filterReq)
		if err != nil {
			ch <- result{err: err}
			return
		}
		ch <- result{graphicData: graphicData}
	}()

	// Collect results
	var resp reportdto.ReportSummaryResponse
	for i := 0; i < 3; i++ {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		if (r.summaryDataBooking != reportdto.DataTotalWithPercentage{}) || (r.summaryDataCancel != reportdto.DataTotalWithPercentage{}) || (r.summaryDataNewUser != reportdto.DataTotalWithPercentage{}) {
			resp.SummaryData = reportdto.SummaryData{
				ConfirmedBooking:    r.summaryDataBooking,
				CancellationBooking: r.summaryDataCancel,
				NewCustomer:         r.summaryDataNewUser,
			}
		}
		if len(r.graphicData) > 0 {
			resp.GraphicData = r.graphicData
		}
	}

	return &resp, nil
}
