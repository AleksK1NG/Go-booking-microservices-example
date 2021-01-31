package http

// MapUserRoutes
func (h *UserHandlers) MapUserRoutes() {
	h.group.POST("/register", h.Register())
	h.group.POST("/login", h.Login())
	h.group.PUT("/:id/avatar", h.UpdateAvatar(), h.mw.SessionMiddleware)
	h.group.GET("/:id", h.GetUserByID())
	h.group.PUT("/:id", h.Update(), h.mw.SessionMiddleware)
	h.group.GET("/me", h.GetMe(), h.mw.SessionMiddleware)
	h.group.GET("/csrf", h.GetCSRFToken(), h.mw.SessionMiddleware)
}
