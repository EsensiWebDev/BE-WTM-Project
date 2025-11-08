package bannerdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateOrderBannerRequest struct {
	Data []OrderBanner `json:"data"`
}

type OrderBanner struct {
	ID    uint `json:"id"`
	Order int  `json:"order"`
}

func (r *UpdateOrderBannerRequest) Validate() error {
	if len(r.Data) == 0 {
		return validation.Errors{
			"data": validation.NewInternalError(fmt.Errorf("data cannot be empty")),
		}
	}

	for _, item := range r.Data {
		if err := validation.ValidateStruct(&item,
			validation.Field(&item.ID, validation.Required.Error("Id is required"), validation.Min(0).Error("Order Id must be greater than or equal to 0")),
			validation.Field(&item.Order, validation.Required.Error("Order is required"), validation.Min(0).Error("Order must be greater than or equal to 0")),
		); err != nil {
			return err
		}
	}

	return nil
}
