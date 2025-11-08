package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/booking_handler"
)

func BookingRoute(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	bookingHandler := booking_handler.NewBookingHandler(app.Usecases.BookingUsecase)

	bookingRouter := routerGroup.Group("/bookings", middlewareMap.Auth)
	{
		cart := bookingRouter.Group("/cart")
		{
			cart.POST("", bookingHandler.AddToCart, middlewareMap.TimeoutSlow)
			cart.GET("", bookingHandler.ListCart)
			cart.DELETE("/:id", bookingHandler.RemoveFromCart)
		}
		bookingRouter.POST("/checkout", bookingHandler.CheckOutCart)
		bookingRouter.POST("/status", bookingHandler.UpdateStatusBooking)
		bookingRouter.GET("", bookingHandler.ListBookings)
		bookingRouter.GET("/booking-status", bookingHandler.ListStatusBooking)
		bookingRouter.GET("/booking-payment", bookingHandler.ListStatusPayment)
		bookingRouter.GET("/history", bookingHandler.ListBookingHistory)
		bookingRouter.GET("/logs", bookingHandler.ListBookingLog)
		bookingRouter.POST("/receipt", bookingHandler.UploadReceipt)
	}
}
