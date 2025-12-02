package booking_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
)

func (br *BookingRepository) CreateInvoice(ctx context.Context, invoices []entity.Invoice) error {
	db := br.db.GetTx(ctx)

	var invoicesModel []model.Invoice
	for _, invoice := range invoices {
		detailJSON, err := json.Marshal(invoice.DetailInvoice)
		if err != nil {
			return fmt.Errorf("failed to marshal detail invoice: %w", err)
		}
		invoiceModel := model.Invoice{
			BookingDetailID: invoice.BookingDetailID,
			InvoiceCode:     invoice.InvoiceCode,
			Detail:          detailJSON,
		}
		invoicesModel = append(invoicesModel, invoiceModel)
	}
	if err := db.WithContext(ctx).Create(&invoicesModel).Error; err != nil {
		return fmt.Errorf("failed to create invoices: %w", err)
	}
	return nil
}
