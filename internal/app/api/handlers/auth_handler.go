package handlers

import (
	"net/http"

	"dot-be-go/internal/service"

	"github.com/labstack/echo/v4"
)

// Register registers a new user
func (h *Handler) Register(c echo.Context) error {
	req := new(service.AuthRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	resp, err := h.AuthService.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

// Login authenticates a user
func (h *Handler) Login(c echo.Context) error {
	req := new(service.AuthRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.AuthService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// GetProfile returns the user profile
func (h *Handler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
