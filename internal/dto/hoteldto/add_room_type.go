package hoteldto

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"wtm-backend/pkg/constant"

	validation "github.com/go-ozzo/ozzo-validation"
)

// AddRoomTypeRequest represents the request structure for adding a new room type to a hotel.
type AddRoomTypeRequest struct {
	HotelID          uint                    `json:"hotel_id" form:"hotel_id"`
	Name             string                  `json:"name" form:"name"`
	Photos           []*multipart.FileHeader `json:"photos" form:"photos"`
	WithoutBreakfast string                  `json:"without_breakfast" form:"without_breakfast"`
	WithBreakfast    string                  `json:"with_breakfast" form:"with_breakfast"`
	RoomSize         float64                 `json:"room_size" form:"room_size"`
	MaxOccupancy     int                     `json:"max_occupancy" form:"max_occupancy"`
	BedTypes         []string                `json:"bed_types" form:"bed_types"`
	IsSmokingRoom    bool                    `json:"is_smoking_room" form:"is_smoking_room"`
	Additional       string                  `json:"additional" form:"additional"`
	Description      string                  `json:"description" form:"description"`
	//TotalUnit        int                     `json:"total_unit" form:"total_unit"`
}

func (r *AddRoomTypeRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if err := validation.ValidateStruct(r,
		validation.Field(&r.HotelID, validation.Required.Error("Hotel ID is required")),
		validation.Field(&r.Name, validation.Required.Error("Name is required")),
		validation.Field(&r.RoomSize, validation.Required.Error("Room Size is required")),
		validation.Field(&r.MaxOccupancy, validation.Required.Error("Max Occupancy is required")),
		validation.Field(&r.BedTypes, validation.Required.Error("Bed Types is required")),
	); err != nil {
		return err
	}

	if len(r.Photos) == 0 {
		errs["photos"] = validation.NewInternalError(fmt.Errorf("at least one photo is required"))
	}

	if r.WithoutBreakfast == "" && r.WithBreakfast == "" {
		errs["without_breakfast"] = validation.NewInternalError(fmt.Errorf("without_breakfast and with_breakfast cannot be both empty"))
		errs["with_breakfast"] = validation.NewInternalError(fmt.Errorf("without_breakfast and with_breakfast cannot be both empty"))
	}

	if !isJSON(r.WithoutBreakfast) && !isJSON(r.WithBreakfast) {
		errs["without_breakfast"] = validation.NewInternalError(fmt.Errorf("without_breakfast and with_breakfast must be filled"))
		errs["with_breakfast"] = validation.NewInternalError(fmt.Errorf("without_breakfast and with_breakfast must be filled"))
	}

	// Validate Additional field if provided
	if len(r.Additional) > 0 {
		var additionalFeatures []RoomAdditional
		if err := json.Unmarshal([]byte(r.Additional), &additionalFeatures); err != nil {
			errs["additional"] = validation.NewInternalError(fmt.Errorf("additional must be a valid JSON array"))
		} else {
			for i, additional := range additionalFeatures {
				if err := additional.Validate(); err != nil {
					errs[fmt.Sprintf("additional[%d]", i)] = err
				}
			}
		}
	}

	if errs != nil && len(errs) > 0 {
		return errs
	}

	return nil
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

type BreakfastBase struct {
	Price  float64 `json:"price" form:"price"`
	IsShow bool    `json:"is_show" form:"is_show"`
}

type BreakfastWith struct {
	BreakfastBase
	Pax int `json:"pax" form:"pax"`
}

type RoomAdditional struct {
	Name       string   `json:"name" form:"name"`
	Category   string   `json:"category" form:"category"`     // "price" or "pax"
	Price      *float64 `json:"price,omitempty" form:"price"` // nullable, used when category="price"
	Pax        *int     `json:"pax,omitempty" form:"pax"`     // nullable, used when category="pax"
	IsRequired bool     `json:"is_required" form:"is_required"`
}

func (r *RoomAdditional) Validate() error {
	var errs validation.Errors = make(map[string]error)

	// Validate Name
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Name is required")),
		validation.Field(&r.Category, validation.Required.Error("Category is required"), validation.In(constant.AdditionalServiceCategoryPrice, constant.AdditionalServiceCategoryPax).Error("Category must be either 'price' or 'pax'")),
	); err != nil {
		return err
	}

	// Validate based on category
	if r.Category == constant.AdditionalServiceCategoryPrice {
		if r.Price == nil {
			errs["price"] = validation.NewInternalError(fmt.Errorf("price is required when category is 'price'"))
		}
		if r.Pax != nil {
			errs["pax"] = validation.NewInternalError(fmt.Errorf("pax must not be set when category is 'price'"))
		}
	} else if r.Category == constant.AdditionalServiceCategoryPax {
		if r.Pax == nil {
			errs["pax"] = validation.NewInternalError(fmt.Errorf("pax is required when category is 'pax'"))
		}
		if r.Price != nil {
			errs["price"] = validation.NewInternalError(fmt.Errorf("price must not be set when category is 'pax'"))
		}
	}

	if errs != nil && len(errs) > 0 {
		return errs
	}

	return nil
}
