package controller

import (
	"bookstack/cmd/service2/internal/service"
	"bookstack/internal/constant"
	"bookstack/internal/dto/response"
	"bookstack/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipperController struct {
	ShipperOrderManageService service.ShipperOrderManageService
}

func NewShipperController(shipperOrderManageService service.ShipperOrderManageService) *ShipperController {
	return &ShipperController{ShipperOrderManageService: shipperOrderManageService}
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

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    shippers,
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
