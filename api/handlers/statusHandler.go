package handlers

import (
	"Projeect/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

func StatusHandler(c echo.Context) error {
	req := new(model.StatusReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "JSON parse failed",
			"desc":  err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"desc":  err.(validator.ValidationErrors),
		})
	}

	user, err := psql.FetchUser(req.NationalCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
			"desc":  err.Error(),
		})
	}

	if user.IPAddress != c.RealIP() {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied.",
			"desc":  "Unauthorized IP address",
		})
	}

	return c.JSON(http.StatusOK, user)
}
