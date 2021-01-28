package http

// MapUserRoutes
func (h *UserHandlers) MapUserRoutes() {
	h.group.POST("/register", h.Register())
	h.group.POST("/login", h.Login())
}
