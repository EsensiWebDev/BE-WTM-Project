package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/booking_handler"
)

func BookingRoute(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	bookingHandler := booking_handler.NewBookingHandler(app.Usecases.BookingUsecase)

	bookingRouter := routerGroup.Group("/bookings", mm.Auth)
	{
		cart := bookingRouter.Group("/cart")
		{
			cart.POST("", mm.TimeoutSlow, bookingHandler.AddToCart)
			cart.GET("", bookingHandler.ListCart)
			cart.DELETE("/:id", bookingHandler.RemoveFromCart)
			cart.POST("/guests", bookingHandler.AddGuestsToCart)
			cart.POST("/sub-guest", bookingHandler.AddGuestToSubCart)
			cart.DELETE("/guests", bookingHandler.RemoveGuestsFromCart)
		}
		bookingRouter.GET("/ids", bookingHandler.ListBookingIDs)
		bookingRouter.GET("/:booking_id/sub-ids", bookingHandler.ListSubBookingIDs)
		bookingRouter.POST("/checkout", mm.TimeoutSlow, bookingHandler.CheckOutCart)
		bookingRouter.GET("", mm.RequirePermission("booking:view"), bookingHandler.ListBookings)
		bookingRouter.GET("/booking-status", mm.RequirePermission("promo:view"), bookingHandler.ListStatusBooking)
		bookingRouter.POST("/booking-status", mm.RequirePermission("booking:edit"), bookingHandler.UpdateStatusBooking)
		bookingRouter.POST("/payment-status", mm.RequirePermission("booking:edit"), bookingHandler.UpdateStatusPayment)
		bookingRouter.GET("/payment-status", mm.RequirePermission("promo:view"), bookingHandler.ListStatusPayment)
		bookingRouter.GET("/history", mm.RequirePermission("promo:view"), bookingHandler.ListBookingHistory)
		bookingRouter.GET("/logs", mm.RequirePermission("promo:view"), bookingHandler.ListBookingLog)
		bookingRouter.POST("/receipt", bookingHandler.UploadReceipt)
		bookingRouter.POST("/:sub_booking_id/cancel", bookingHandler.CancelBooking)
	}
}
