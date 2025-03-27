package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type OrderService interface {
	CreateOrder(request.OrderRequest) (models.Order, error)
	CancelOrder(int) error
}

type OrderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderServiceImpl(repository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		repo: repository,
	}
}

func (o *OrderServiceImpl) CreateOrder(request request.OrderRequest) (models.Order, error) {
	return o.repo.CreateOrder(request)
}

func (o *OrderServiceImpl) CancelOrder(orderId int) error {
	return o.repo.CancelOrder(orderId)
}
