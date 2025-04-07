package controller

import (
	"bookstack/config"
	"bookstack/internal/constant"
	"bookstack/internal/dto/response"
	"bookstack/internal/messaging"
	"bookstack/internal/models"
	"bookstack/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipperController struct {
	ShipperOrderManageService service.ShipperOrderManageService
	UserService               service.UserService
}

func NewShipperController(shipperOrderManageService service.ShipperOrderManageService, userService service.UserService) *ShipperController {
	return &ShipperController{ShipperOrderManageService: shipperOrderManageService, UserService: userService}
}

func (c *ShipperController) GetReceivedOrders(ctx *gin.Context) {
	var webResponse response.WebResponse
	header := ctx.Request.Header
	token := header.Get("Authorization")
	userId, err := c.UserService.GetUserIdByToken(token)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		}
		ctx.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	orders, err := c.ShipperOrderManageService.GetReceivedOrders(userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var orderResponse []response.OrderResponse
	for _, order := range orders {
		orderResponse = append(orderResponse, c.CoppyToOrderResponse(order))
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Success",
		Data:    orderResponse,
	}
	ctx.JSON(http.StatusOK, webResponse)
}
func (c *ShipperController) ReceiveOrder(ctx *gin.Context) {
	webResponse := response.WebResponse{}
	orderIdStr := ctx.Param("orderId")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid order ID format",
			Data:    nil,
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}
	header := ctx.Request.Header
	token := header.Get("Authorization")
	userId, err := c.UserService.GetUserIdByToken(token)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		}
		ctx.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	err = c.ShipperOrderManageService.ReceiveOrder(orderId, userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Order received successfully",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

func (c *ShipperController) GetAllShipper(ctx *gin.Context) {
	var webResponse response.WebResponse

	shippers, err := c.ShipperOrderManageService.GetAllShipper("shipper")
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var userResponse []response.UserResponse
	for _, shipper := range shippers {
		userResponse = append(userResponse, c.CoppyToUserResponse(shipper))
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    userResponse,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

func (c *ShipperController) GetOrderInRange(ctx *gin.Context) {
	var webResponse response.WebResponse

	// Lấy giá trị từ form-data
	place := ctx.PostForm("place")
	if place == "" {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Place is required",
			Data:    nil,
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	orders, err := c.ShipperOrderManageService.GetOrderInRange(place)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var orderResponse []response.OrderResponse
	for _, order := range orders {
		orderResponse = append(orderResponse, c.CoppyToOrderResponse(order))
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Success",
		Data:    orderResponse,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *ShipperController) CoppyToUserResponse(user models.User) response.UserResponse {
	var userResponse response.UserResponse
	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	return userResponse
}

func (controller *ShipperController) CoppyToOrderResponse(order models.Order) response.OrderResponse {
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

func (controller *ShipperController) StartListeningForNewOrders() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	rabbitmq, err := messaging.NewRabbitMQ(conf)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}

	err = rabbitmq.ConsumeNewOrders(func(orderID uint, address string) {
		log.Printf("Processing new order: ID=%d, Address=%s", orderID, address)

		// Get all shippers
		shippers, err := controller.ShipperOrderManageService.GetAllShipper("shipper")
		if err != nil {
			log.Printf("Failed to get shippers: %v", err)
			return
		}

		// Find shippers in the same area
		for _, shipper := range shippers {
			if shipper.WorkingArea == address {
				log.Printf("Found matching shipper: ID=%d, Area=%s", shipper.ID, shipper.WorkingArea)
				// Here you can add logic to notify the shipper
				// For example, send a notification or update a database
			}
		}
	})
	if err != nil {
		log.Printf("Failed to start consuming new orders: %v", err)
	}
}
