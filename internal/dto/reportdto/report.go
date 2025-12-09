package reportdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ReportRequest struct {
	DateFrom              string `json:"date_from" form:"date_from"`
	DateTo                string `json:"date_to" form:"date_to"`
	HotelID               []uint `json:"hotel_id" form:"hotel_id"`
	AgentCompanyID        []uint `json:"agent_company_id" form:"agent_company_id"`
	dto.PaginationRequest `json:",inline"`
}

type ReportSummaryRequest struct {
	DateFrom              string `json:"date_from" form:"date_from"`
	DateTo                string `json:"date_to" form:"date_to"`
	dto.PaginationRequest `json:",inline"`
}

type ReportSummaryResponse struct {
	SummaryData SummaryData             `json:"summary_data"`
	GraphicData []entity.ReportForGraph `json:"graphic_data"`
}

type ReportAgentResponse struct {
	Data  []entity.ReportAgentBooking `json:"data"`
	Total int64                       `json:"total"`
}

type SummaryData struct {
	ConfirmedBooking DataTotalWithPercentage `json:"confirmed_booking"`
	CancelledBooking DataTotalWithPercentage `json:"cancelled_booking"`
	RejectedBooking  DataTotalWithPercentage `json:"rejected_booking"`
}

type DataTotalWithPercentage struct {
	Count   int64   `json:"count"`
	Percent float64 `json:"percent"`
	Message string  `json:"message"`
}

type ReportAgentDetailRequest struct {
	HotelID               uint   `json:"hotel_id" form:"hotel_id"`
	AgentID               uint   `json:"agent_id" form:"agent_id"`
	DateFrom              string `json:"date_from" form:"date_from"`
	DateTo                string `json:"date_to" form:"date_to"`
	dto.PaginationRequest `json:",inline"`
}

type ReportAgentDetailResponse struct {
	ReportAgentDetailData []entity.ReportAgentDetail `json:"report_agent_detail_data"`
	Total                 int64                      `json:"total"`
}
