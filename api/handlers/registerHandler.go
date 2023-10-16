package handlers

import (
	"Projeect/internal/model"
	"Projeect/utils"
	"Projeect/utils/datasource"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	QUEUE = "reqs"
)

func RegisterHandler(c echo.Context) error {
	var req model.RegisterReq
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "JSON parse failed",
			"desc":  err.Error(),
		})
	}

	req.IPAddress = c.RealIP()
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"desc":  err.(validator.ValidationErrors),
		})
	}
	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		NationalCode: utils.EncodeBase64(req.NationalCode),
		IPAddress:    req.IPAddress,
	}

	_, err = psql.SaveUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database record push failed",
			"desc":  err.Error(),
		})
	}

	uuid1 := utils.EncodeBase64(req.NationalCode) + "@1"

	err = datasource.UploadPic([]byte(req.Pic1), uuid1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed tp push image S3",
			"desc":  err.Error(),
		})
	}

	uuid2 := utils.EncodeBase64(req.NationalCode) + "@2"

	err = datasource.UploadPic([]byte(req.Pic1), uuid2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to push image S3",
			"desc":  err.Error(),
		})
	}

	err = enqueueNationalCode(user.NationalCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to publish request on RabbitMQ",
			"desc":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.RegisterRes{
		Message: "Successfully registered",
	})
}

func enqueueNationalCode(code string) error {
	ctx := context.Background()
	err := rabbitMQ.Publish(ctx, QUEUE, code)
	if err != nil {
		return err
	}

	return nil
}
