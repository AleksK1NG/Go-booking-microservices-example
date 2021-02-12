package v1

// MapRoutes
func (h *hotelsHandlers) MapRoutes() {
	h.group.GET("", h.GetHotels())
	h.group.GET("/:hotel_id", h.GetHotelByID())
	h.group.POST("", h.CreateHotel(), h.mw.SessionMiddleware)
	h.group.PUT("/:hotel_id", h.UpdateHotel(), h.mw.SessionMiddleware)
	h.group.PUT("/:hotel_id/image", h.UploadImage(), h.mw.SessionMiddleware)
}
