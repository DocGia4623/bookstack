package service

import (
	"bookstack/internal/constant"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type ShipperOrderManageService interface {
	// Quản lý shipper
	GetAllShipper(string) ([]models.User, error)

	// Quản lý đơn hàng của shipper
	AssignOrderToShipper(orderID uint, shipperID uint) error
	GetOrdersByShipper(shipperID uint) ([]models.Order, error)
	UpdateOrderStatus(orderID uint, status constant.OrderStatus) error
	GetPendingOrders() ([]models.Order, error)
	GetOrderInRange(string) ([]models.Order, error)
	// Shipper nhận đơn hàng
	ReceiveOrder(orderId int, userId int) error
	GetReceivedOrders(userId int) ([]models.Order, error)
}

type shipperOrderManageService struct {
	ShipperRepository repository.ShipperRepository
}

func NewOrderManageService(shipperRepository repository.ShipperRepository) ShipperOrderManageService {
	return &shipperOrderManageService{
		ShipperRepository: shipperRepository,
	}
}

func (s *shipperOrderManageService) GetReceivedOrders(userId int) ([]models.Order, error) {
	return s.ShipperRepository.GetReceivedOrders(userId)
}

func (s *shipperOrderManageService) ReceiveOrder(orderId int, userId int) error {
	return s.ShipperRepository.ReceiveOrder(orderId, userId)
}

func (s *shipperOrderManageService) GetOrderInRange(place string) ([]models.Order, error) {
	return s.ShipperRepository.GetOrderInRange(place)
}

// Quản lý shipper
func (s *shipperOrderManageService) GetAllShipper(roleName string) ([]models.User, error) {
	return s.ShipperRepository.GetAllShipper(roleName)
}

// Quản lý đơn hàng của shipper
func (s *shipperOrderManageService) AssignOrderToShipper(orderID uint, shipperID uint) error {
	return s.ShipperRepository.AssignOrderToShipper(orderID, shipperID)
}

func (s *shipperOrderManageService) GetOrdersByShipper(shipperID uint) ([]models.Order, error) {
	return s.ShipperRepository.GetOrdersByShipper(shipperID)
}

func (s *shipperOrderManageService) UpdateOrderStatus(orderID uint, status constant.OrderStatus) error {
	return s.ShipperRepository.UpdateOrderStatus(orderID, status)
}

func (s *shipperOrderManageService) GetPendingOrders() ([]models.Order, error) {
	return s.ShipperRepository.GetPendingOrders()
}
