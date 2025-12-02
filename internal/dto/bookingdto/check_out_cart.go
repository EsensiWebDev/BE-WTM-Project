package bookingdto

import "wtm-backend/internal/domain/entity"

type CheckOutCartResponse struct {
	Invoice []DataInvoice `json:"invoice"`
}

type DataInvoice struct {
	entity.DetailInvoice `json:",inline"`
	InvoiceNumber        string `json:"invoice_number"`
	InvoiceDate          string `json:"invoice_date"`
	Receipt              string `json:"receipt"`
}
