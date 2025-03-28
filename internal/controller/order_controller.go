package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(serv service.OrderService) *OrderController {
	return &OrderController{
		service: serv,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create an order based on the provided request data
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body request.OrderRequest true "Order request payload"
// @Success 200 {object} response.WebResponse "Order created successfully"
// @Failure 400 {object} response.WebResponse "Invalid request"
// @Failure 500 {object} response.WebResponse "Server error"
// @Router /orders [post]
func (controller *OrderController) CreateOrder(c *gin.Context) {
	var webResponse response.WebResponse
	var request request.OrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	order, err := controller.service.CreateOrder(request)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Server error",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Order Created",
		Data:    order,
	}
	c.JSON(http.StatusOK, webResponse)
}
