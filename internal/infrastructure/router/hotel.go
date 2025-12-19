package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/hotel_handler"

	"github.com/gin-gonic/gin"
)

func HotelRoute(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	hotelHandler := hotel_handler.NewHotelHandler(app.Usecases.HotelUsecase, app.Config)

	group := routerGroup.Group("", mm.TimeoutFast)
	{
		//group.GET("/example_hotel_data", hotelHandler.HotelDataDummy)
		//group.GET("/example_room_type_data", hotelHandler.RoomTypeDataDummy)
		group.GET("/example_nearby_places", hotelHandler.ExampleDataNearbyPlace)
		group.GET("/example_social_medias", hotelHandler.ExampleDataSocialMedias)
		group.GET("/example_additional_features", hotelHandler.ExampleDataAdditionalFeatures)
		group.GET("/example_room_options", hotelHandler.ExampleRoomOptions)

		hotels := group.Group("/hotels")
		{
			hotels.GET("", mm.Auth, mm.RequirePermission("hotel:view"), hotelHandler.ListHotels)
			hotels.POST("", mm.Auth, mm.RequirePermission("hotel:create"), mm.TimeoutFile, hotelHandler.CreateHotel)
			hotels.GET("/download-format", mm.Auth, hotelHandler.DownloadFormat)
			hotels.POST("/upload", mm.Auth, mm.RequirePermission("hotel:create"), mm.TimeoutFile, hotelHandler.UploadHotel)
			hotels.PUT("/:id", mm.Auth, mm.RequirePermission("hotel:edit"), mm.TimeoutFile, hotelHandler.UpdateHotel)
			hotels.GET("/:id", mm.Auth, mm.RequirePermission("hotel:view"), hotelHandler.DetailHotel)
			hotels.DELETE("/:id", mm.Auth, mm.RequirePermission("hotel:delete"), hotelHandler.RemoveHotel)

			agents := hotels.Group("/agent")
			{
				agents.GET("", mm.Auth, hotelHandler.ListHotelsForAgent)
				agents.GET("/:id", mm.Auth, hotelHandler.DetailHotelForAgent)
			}

			roomTypes := hotels.Group("/room-types", mm.Auth)
			{
				roomTypes.GET("", mm.RequirePermission("hotel:view"), hotelHandler.ListRoomTypes)
				roomTypes.POST("", mm.RequirePermission("hotel:edit"), hotelHandler.AddRoomType, mm.TimeoutFile)
				roomTypes.PUT("/:id", mm.RequirePermission("hotel:edit"), hotelHandler.UpdateRoomType, mm.TimeoutFile)
				roomTypes.DELETE("/:id", mm.RequirePermission("hotel:edit"), hotelHandler.RemoveRoomType)
			}

			hotels.GET("/bed-types", mm.Auth, hotelHandler.ListAllBedTypes)
			hotels.GET("/facilities", mm.Auth, hotelHandler.ListFacilities)
			hotels.GET("/additional-rooms", mm.Auth, hotelHandler.ListAdditionalRooms)
			hotels.GET("/room-available", mm.Auth, hotelHandler.ListRoomAvailable)
			hotels.PUT("/room-available", mm.Auth, hotelHandler.UpdateRoomAvailable)
			hotels.GET("/statuses", mm.Auth, hotelHandler.ListStatusHotel)
			hotels.PUT("/status", mm.Auth, mm.RequirePermission("hotel:edit"), hotelHandler.UpdateStatus)
			hotels.GET("/provinces", hotelHandler.ListProvinces)
		}
	}
}
