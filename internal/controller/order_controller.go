package controller

import (
	"bookstack/config"
	"bookstack/internal/constant"
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/models"
	"bookstack/internal/service"
	"bookstack/utils"
	"net/http"
	"strconv"

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

	var orderResponse response.OrderResponse
	// Copy basic order information
	orderResponse.OrderID = order.ID
	orderResponse.UserID = order.UserID
	orderResponse.TotalPrice = order.TotalPrice
	orderResponse.CreatedAt = order.CreatedAt.Format("2006-01-02 15:04:05")
	orderResponse.UpdatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")

	// Convert status from constant to string
	switch order.Status {
	case constant.Pending:
		orderResponse.Status = "Pending"
	case constant.Confirmed:
		orderResponse.Status = "Confirmed"
	case constant.Processing:
		orderResponse.Status = "Processing"
	case constant.Shipped:
		orderResponse.Status = "Shipped"
	case constant.Delivered:
		orderResponse.Status = "Delivered"
	case constant.Cancelled:
		orderResponse.Status = "Cancelled"
	case constant.Returned:
		orderResponse.Status = "Returned"
	case constant.Failed:
		orderResponse.Status = "Failed"
	}

	// Copy order details
	for _, item := range order.OrderDetail {
		var orderDetailResponse response.OrderDetailResponse
		orderDetailResponse.Quantity = item.Quantity
		orderDetailResponse.Price = item.Price

		// Copy book information
		var bookResponse response.BookOrderResponse
		bookResponse.ID = item.Book.ID
		bookResponse.Title = item.Book.Title
		bookResponse.Price = item.Book.Price
		bookResponse.Slug = item.Book.Slug
		bookResponse.CreatedBy = strconv.FormatUint(uint64(item.Book.CreatedBy), 10)

		orderDetailResponse.Book = bookResponse
		orderResponse.OrderDetail = append(orderResponse.OrderDetail, orderDetailResponse)
	}

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Order Created",
		Data:    orderResponse,
	}
	c.JSON(http.StatusOK, webResponse)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get an order by ID
// @Tags Orders
// @Produce json
// @Param orderId path int true "Order ID"
// @Success 200 {object} response.WebResponse "Order retrieved successfully"
// @Failure 400 {object} response.WebResponse "Invalid request"
// @Failure 500 {object} response.WebResponse "Server error"
// @Router /orders/{orderId} [get]
func (controller *OrderController) GetUserOrder(c *gin.Context) {
	var webResponse response.WebResponse
	token := c.GetHeader("Authorization")
	conf, err := config.LoadConfig()
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to load config",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	sub, err := utils.ValidateAccessToken(token, conf.AccessTokenSecret)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	userIdFloat, ok := sub.(float64)
	if !ok {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Invalid user ID format in token",
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	userId := int(userIdFloat)
	orders, err := controller.service.GetUserOrder(userId)
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
	var orderResponse []response.OrderResponse
	for _, order := range orders {
		orderResponse = append(orderResponse, controller.CoppyToOrderResponse(order))
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Order retrieved successfully",
		Data:    orderResponse,
	}
	c.JSON(http.StatusOK, webResponse)
}

// CancelOrder godoc
// @Summary Cancel an order by ID
// @Description Cancel an order by ID
// @Tags Orders
// @Produce json
// @Param orderId path int true "Order ID"
// @Success 200 {object} response.WebResponse "Order cancelled successfully"
// @Failure 400 {object} response.WebResponse "Invalid request"
// @Failure 500 {object} response.WebResponse "Server error"
// @Router /orders/{orderId}/cancel [post]
func (controller *OrderController) CancelOrder(c *gin.Context) {
	var webResponse response.WebResponse
	orderIdStr := c.Param("orderId")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.service.CancelOrder(orderId)
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
		Message: "Order cancelled successfully",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *OrderController) CoppyToOrderResponse(order models.Order) response.OrderResponse {
	var orderResponse response.OrderResponse

	// Copy basic order information
	orderResponse.OrderID = order.ID
	orderResponse.UserID = order.UserID
	orderResponse.TotalPrice = order.TotalPrice
	orderResponse.Address = order.Address
	orderResponse.Phone = order.Phone
	orderResponse.CreatedAt = order.CreatedAt.Format("2006-01-02 15:04:05")
	orderResponse.UpdatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")

	// Convert status from constant to string
	switch order.Status {
	case constant.Pending:
		orderResponse.Status = "Pending"
	case constant.Confirmed:
		orderResponse.Status = "Confirmed"
	case constant.Processing:
		orderResponse.Status = "Processing"
	case constant.Shipped:
		orderResponse.Status = "Shipped"
	case constant.Delivered:
		orderResponse.Status = "Delivered"
	case constant.Cancelled:
		orderResponse.Status = "Cancelled"
	case constant.Returned:
		orderResponse.Status = "Returned"
	case constant.Failed:
		orderResponse.Status = "Failed"
	}

	// Copy order details
	for _, item := range order.OrderDetail {
		var orderDetailResponse response.OrderDetailResponse
		orderDetailResponse.Quantity = item.Quantity
		orderDetailResponse.Price = item.Price

		// Copy book information
		var bookResponse response.BookOrderResponse
		bookResponse.ID = item.Book.ID
		bookResponse.Title = item.Book.Title
		bookResponse.Price = item.Book.Price
		bookResponse.Slug = item.Book.Slug
		bookResponse.CreatedBy = strconv.FormatUint(uint64(item.Book.CreatedBy), 10)

		orderDetailResponse.Book = bookResponse
		orderResponse.OrderDetail = append(orderResponse.OrderDetail, orderDetailResponse)
	}

	return orderResponse
}
