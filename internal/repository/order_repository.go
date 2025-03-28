package repository

import (
	"bookstack/internal/constant"
	"bookstack/internal/dto/request"
	"bookstack/internal/models"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(request.OrderRequest) (models.Order, error)
	GetOrder(int) ([]models.Order, error)
	CancelOrder(int) error
}

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{
		DB: db,
	}
}

func (o *OrderRepositoryImpl) GetOrder(userId int) ([]models.Order, error) {
	var orders []models.Order

	err := o.DB.Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
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

func (o *OrderRepositoryImpl) CreateOrder(request request.OrderRequest) (models.Order, error) {
	var order models.Order
	var totalPrice float64

	// Copy thông tin từ request sang order (chỉ copy được các field cùng kiểu)
	err := copier.Copy(&order, request)
	if err != nil {
		return models.Order{}, err
	}

	// Chuyển đổi OrderDetailRequest thành OrderDetail
	var orderDetails []models.OrderDetail
	for _, detail := range request.OrderDetails {
		price, err := o.GetBookPrice(int(detail.BookID))
		if err != nil {
			return models.Order{}, err
		}
		orderDetails = append(orderDetails, models.OrderDetail{
			BookID:   detail.BookID,
			Quantity: detail.Quantity,
			Price:    price, // Lấy giá sách từ DB thay vì hardcode
		})
		totalPrice += price * float64(detail.Quantity)
	}

	// Gán giá trị còn thiếu vào Order
	order.OrderDetail = orderDetails
	order.TotalPrice = totalPrice
	order.Status = constant.Pending

	// Lưu vào DB (dùng &order để tránh lỗi reflect)
	err = o.DB.Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (o *OrderRepositoryImpl) CancelOrder(orderId int) error {
	var order models.Order
	err := o.DB.Where("order_id = ?", orderId).Find(&order).Error
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
