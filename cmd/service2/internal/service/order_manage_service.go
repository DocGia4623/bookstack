package service

import "bookstack/internal/repository"

type OrderManageService interface {
	ManageShiper(orderID uint) error
}

type orderManageService struct {
	orderRepository repository.OrderRepository
}

func NewOrderManageService(orderRepository repository.OrderRepository) OrderManageService {
	return &orderManageService{
		orderRepository: orderRepository,
	}
}

func (s *orderManageService) ManageShiper(orderID uint) error {
	return nil
}
