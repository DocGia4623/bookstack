package repository

import (
	"bookstack/internal/constant"
	"bookstack/internal/dto/request"
	"bookstack/internal/models"

	"fmt"
	"strconv"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(request.OrderRequest, int) (models.Order, error)
	GetOrder(int) (models.Order, error)
	GetUserOrder(int) ([]models.Order, error)
	CancelOrder(int) error
	UpdateOrderStatus(webhookPayload map[string]interface{}) error
}

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{
		DB: db,
	}
}

func (o *OrderRepositoryImpl) GetUserOrder(userId int) ([]models.Order, error) {
	var orders []models.Order
	err := o.DB.Preload("OrderDetail").Preload("OrderDetail.Book").Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepositoryImpl) GetOrder(id int) (models.Order, error) {
	var orders models.Order

	err := o.DB.Where("id = ?", id).Find(&orders).Error
	if err != nil {
		return models.Order{}, err
	}
	return orders, err
}

func (o *OrderRepositoryImpl) GetBookPrice(bookId int) (float64, error) {
	var book models.Book
	err := o.DB.Where("id = ?", bookId).First(&book).Error
	if err != nil {
		return 0, err
	}
	return book.Price, err
}

func (o *OrderRepositoryImpl) CreateOrder(request request.OrderRequest, userId int) (models.Order, error) {
	var order models.Order
	var totalPrice float64

	// Copy thông tin từ request sang order (chỉ copy được các field cùng kiểu)
	err := copier.Copy(&order, request)
	if err != nil {
		return models.Order{}, err
	}

	// Copy các trường không được copy tự động
	order.Address = request.Address
	order.Phone = request.Phone
	order.UserID = uint(userId)

	// Chuyển đổi OrderDetailRequest thành OrderDetail
	var orderDetails []models.OrderDetail
	for _, detail := range request.OrderDetails {
		// Lấy thông tin sách từ DB
		var book models.Book
		if err := o.DB.First(&book, detail.BookID).Error; err != nil {
			return models.Order{}, fmt.Errorf("book not found: %w", err)
		}

		orderDetails = append(orderDetails, models.OrderDetail{
			BookID:   detail.BookID,
			Quantity: detail.Quantity,
			Price:    book.Price,
			Book:     book,
		})
		totalPrice += book.Price * float64(detail.Quantity)
	}

	// Gán giá trị còn thiếu vào Order
	order.OrderDetail = orderDetails
	order.TotalPrice = totalPrice
	order.Status = constant.Pending

	// Lưu vào DB với preload Book
	err = o.DB.Preload("OrderDetail.Book").Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (o *OrderRepositoryImpl) CancelOrder(orderId int) error {
	var order models.Order
	err := o.DB.Where("id = ?", orderId).Find(&order).Error
	if err != nil {
		return err
	}
	order.Status = constant.Cancelled
	err = o.DB.Save(order).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepositoryImpl) UpdateOrderStatus(webhookPayload map[string]interface{}) error {
	// Extract order ID from webhook payload
	orderID, ok := webhookPayload["order_id"].(string)
	if !ok {
		return fmt.Errorf("invalid order ID in webhook payload")
	}

	// Convert order ID to uint
	id, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %v", err)
	}

	// Update order status to Confirmed
	return o.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", constant.Confirmed).Error
}
