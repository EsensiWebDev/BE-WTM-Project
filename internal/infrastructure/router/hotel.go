package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/hotel_handler"
)

func HotelRoute(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	hotelHandler := hotel_handler.NewHotelHandler(app.Usecases.HotelUsecase, app.Config)

	group := routerGroup.Group("", middlewareMap.TimeoutFast)
	{
		//group.GET("/example_hotel_data", hotelHandler.HotelDataDummy)
		//group.GET("/example_room_type_data", hotelHandler.RoomTypeDataDummy)
		group.GET("/example_nearby_places", hotelHandler.ExampleDataNearbyPlace)
		group.GET("/example_social_medias", hotelHandler.ExampleDataSocialMedias)
		group.GET("/example_additional_features", hotelHandler.ExampleDataAdditionalFeatures)
		group.GET("/example_room_options", hotelHandler.ExampleRoomOptions)

		hotels := group.Group("/hotels", middlewareMap.Auth)
		{
			hotels.GET("", hotelHandler.ListHotels)
			hotels.POST("", hotelHandler.CreateHotel, middlewareMap.TimeoutFile)
			hotels.PUT("/:id", hotelHandler.UpdateHotel, middlewareMap.TimeoutFile)
			hotels.GET("/:id", hotelHandler.DetailHotel)
			hotels.DELETE("/:id", hotelHandler.RemoveHotel)

			agents := hotels.Group("/agent")
			{
				agents.GET("", hotelHandler.ListHotelsForAgent)
				agents.GET("/:id", hotelHandler.DetailHotelForAgent)
			}

			roomTypes := hotels.Group("/room-types")
			{
				roomTypes.GET("", hotelHandler.ListRoomTypes)
				roomTypes.POST("", hotelHandler.AddRoomType, middlewareMap.TimeoutFile)
				roomTypes.PUT("/:id", hotelHandler.UpdateRoomType, middlewareMap.TimeoutFile)
				roomTypes.DELETE("/:id", hotelHandler.RemoveRoomType)
			}

			hotels.GET("/bed-types", hotelHandler.ListAllBedTypes)
			hotels.GET("/facilities", hotelHandler.ListFacilities)
			hotels.GET("/additional-rooms", hotelHandler.ListAdditionalRooms)
			hotels.GET("/room-available", hotelHandler.ListRoomAvailable)
			hotels.PUT("/room-available", hotelHandler.UpdateRoomAvailable)
			hotels.GET("/statuses", hotelHandler.ListStatusHotel)
			hotels.PUT("/status", hotelHandler.UpdateStatus)
			hotels.GET("/provinces", hotelHandler.ListProvinces)
		}
	}
}
