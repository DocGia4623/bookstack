package repository

import (
	"bookstack/internal/constant"
	"bookstack/internal/models"

	"gorm.io/gorm"
)

type ShipperRepository interface {
	// Quản lý shipper
	GetAllShipper(string) ([]models.User, error)
	// Quản lý đơn hàng của shipper
	AssignOrderToShipper(orderID uint, shipperID uint) error
	GetOrderInRange(string) ([]models.Order, error)
	GetOrdersByShipper(shipperID uint) ([]models.Order, error)
	UpdateOrderStatus(orderID uint, status constant.OrderStatus) error
	GetPendingOrders() ([]models.Order, error)
	ReceiveOrder(orderId int, userId int) error
	GetReceivedOrders(userId int) ([]models.Order, error)
}

type shipperRepository struct {
	db *gorm.DB
}

func NewShipperRepository(db *gorm.DB) ShipperRepository {
	return &shipperRepository{db: db}
}

func (r *shipperRepository) GetReceivedOrders(userId int) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("shipper_id = ?", userId).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *shipperRepository) ReceiveOrder(orderId int, userId int) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderId).Update("shipper_id", userId).Error
}

func (r *shipperRepository) GetOrderInRange(place string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("address LIKE ?", "%"+place+"%").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *shipperRepository) GetAllShipper(roleName string) ([]models.User, error) {
	var shippers []models.User
	err := r.db.
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", roleName).
		Find(&shippers).Error
	if err != nil {
		return nil, err
	}
	return shippers, nil
}

func (r *shipperRepository) AssignOrderToShipper(orderID uint, shipperID uint) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("shipper_id", shipperID).Error
}

func (r *shipperRepository) GetOrdersByShipper(shipperID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("shipper_id = ?", shipperID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *shipperRepository) UpdateOrderStatus(orderID uint, status constant.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *shipperRepository) GetPendingOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("status = ?", constant.Pending).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
