package handlers

import (
	"taobin-service/internal/core/domain"
	"taobin-service/internal/core/ports"
	"taobin-service/utils/validator"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	srv       ports.Service
	db        *gorm.DB
	validator validator.Validator
}

func New(srv ports.Service, db *gorm.DB) *HTTPHandler {
	return &HTTPHandler{
		srv:       srv,
		db:        db,
		validator: validator.New(),
	}
}

func (hdl *HTTPHandler) TestCheck(c *fiber.Ctx) error {
	sqlDB, err := hdl.db.DB()
	if err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ResponseBody{Status: domain.InternalServerError})
	}

	err = sqlDB.Ping()
	if err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ResponseBody{Status: domain.InternalServerError})
	}
	return c.Status(fiber.StatusOK).JSON(domain.ResponseBody{Status: domain.Success, Data: ""})
}

func (hdl *HTTPHandler) UpdateStock(c *fiber.Ctx) error {
	var request domain.StockRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
	}
	if err := hdl.validator.ValidateStruct(request); err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
	}
	if request.ID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
	}
	response, err := hdl.srv.UpdateStock(request)
	if err != nil {
		msg := domain.ResponseBody{
			Status: domain.InternalServerError,
		}
		msg.Status.Message = []string{
			err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(msg)
	}
	return c.Status(fiber.StatusOK).JSON(domain.ResponseBody{Status: domain.Success, Data: response})
}

func (hdl *HTTPHandler) GetStocks(c *fiber.Ctx) error {
	var err error
	var data []domain.StockResponse
	condition := domain.QueryStockRequest{}
	err = c.QueryParser(&condition)
	if err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
	}

	err = hdl.validator.ValidateStruct(condition)
	if err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
	}

	idStr := c.Params("id")
	if idStr != "" {
		id, err := c.ParamsInt("id")
		if err != nil {
			logrus.Errorln(err)
			return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseBody{Status: domain.BadRequest})
		}
		condition.ID = &id
	}
	result, err := hdl.srv.GetStock(condition)
	if err != nil {
		logrus.Errorln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ResponseBody{Status: domain.InternalServerError})
	}
	if result.Stocks == nil {
		data = make([]domain.StockResponse, 0)
	} else {
		data = result.Stocks
	}

	return c.Status(fiber.StatusOK).JSON(domain.ResponseBody{
		Status:      domain.Success,
		Data:        data,
		CurrentPage: result.CurrentPage,
		PerPage:     result.PerPage,
		TotalItem:   result.TotalItem,
	})
}
