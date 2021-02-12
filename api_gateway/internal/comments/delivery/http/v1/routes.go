package v1

// MapRoutes
func (c *commentsHandlers) MapRoutes() {
	c.group.GET("/:comment_id", c.GetCommByID())
	c.group.POST("", c.CreateComment(), c.mw.SessionMiddleware)
	c.group.PUT("/:comment_id", c.UpdateComment(), c.mw.SessionMiddleware)
	c.group.PUT("/comments/hotel/:hotel_id", c.GetByHotelID())
}
