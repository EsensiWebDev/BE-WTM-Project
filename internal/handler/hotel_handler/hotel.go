package hotel_handler

import (
	"net/http"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/usecase/hotel_usecase"
	"wtm-backend/pkg/constant"

	"github.com/gin-gonic/gin"
)

type HotelHandler struct {
	hotelUsecase domain.HotelUsecase
	config       *config.Config
}

func NewHotelHandler(hotelUsecase *hotel_usecase.HotelUsecase, config *config.Config) *HotelHandler {
	return &HotelHandler{
		hotelUsecase: hotelUsecase,
		config:       config,
	}
}

// ExampleDataNearbyPlace godoc
// @Summary Example Data Nearby Places
// @Description Example payload JSON for nearby places
// @Tags Hotel
// @Accept json
// @Produce json
// @Success 200 {object} []hoteldto.NearbyPlace "Payload nearby places in JSON format"
// @Router /example_nearby_places [get]
func (hh *HotelHandler) ExampleDataNearbyPlace(c *gin.Context) {
	nearbyPlaces := []hoteldto.NearbyPlace{
		{
			Name:     "Pantai Indah",
			Distance: 1.2,
		},
		{
			Name:     "Mall Central",
			Distance: 0.5,
		},
	}

	c.JSON(http.StatusOK, nearbyPlaces)
}

// ExampleDataSocialMedias godoc
// @Summary Example Data Social Media
// @Description Example payload JSON for social medias
// @Tags Hotel
// @Accept json
// @Produce json
// @Success 200 {object} []hoteldto.SocialMedia "Payload social medias in JSON format"
// @Router /example_social_medias [get]
func (hh *HotelHandler) ExampleDataSocialMedias(c *gin.Context) {
	socialMedias := []hoteldto.SocialMedia{
		{
			Platform: "instagram",
			Link:     "https://www.instagram.com/hotel-box",
		},
		{
			Platform: "tiktok",
			Link:     "https://www.tiktok.com/hotel-box",
		},
	}

	c.JSON(http.StatusOK, socialMedias)
}

// ExampleDataAdditionalFeatures godoc
// @Summary Example Data Additional Features
// @Description Example payload JSON for additional room features
// @Tags Hotel
// @Accept json
// @Produce json
// @Success 200 {object} []hoteldto.RoomAdditional "Payload additional room features in JSON format"
// @Router /example_additional_features [get]
func (hh *HotelHandler) ExampleDataAdditionalFeatures(c *gin.Context) {
	priceValue1 := 50000.0
	priceValue2 := 75000.0
	paxValue := 2

	additionalFeatures := []hoteldto.RoomAdditional{
		{
			Name:       "Extra Bed",
			Category:   constant.AdditionalServiceCategoryPrice,
			Price:      &priceValue1,
			IsRequired: false,
		},
		{
			Name:       "Extra Sofa",
			Category:   constant.AdditionalServiceCategoryPrice,
			Price:      &priceValue2,
			IsRequired: false,
		},
		{
			Name:       "Extra Guest",
			Category:   constant.AdditionalServiceCategoryPax,
			Pax:        &paxValue,
			IsRequired: true,
		},
	}

	c.JSON(http.StatusOK, additionalFeatures)
}

// ExampleRoomOptions godoc
// @Summary Example Room Options
// @Description Example payload for room options based on breakfast type
// @Tags Hotel
// @Accept json
// @Produce json
// @Param room_options query string true "Type of room options" Enums(without_breakfast, with_breakfast)
// @Success 200 {object} hoteldto.BreakfastBase "Payload room options without breakfast in JSON format"
// @Success 200 {object} hoteldto.BreakfastWith "Payload room options with breakfast in JSON format"
// @Router /example_room_options [get]
func (hh *HotelHandler) ExampleRoomOptions(c *gin.Context) {

	roomOptions := hoteldto.BreakfastBase{
		Price:  150000,
		IsShow: true,
	}

	typeRoomOptions := c.Query("room_options")
	if typeRoomOptions == "with_breakfast" {
		roomOptionsWith := hoteldto.BreakfastWith{
			BreakfastBase: hoteldto.BreakfastBase{
				Price:  200000,
				IsShow: true,
			},
			Pax: 2,
		}
		c.JSON(http.StatusOK, roomOptionsWith)
		return
	}

	c.JSON(http.StatusOK, roomOptions)
}
