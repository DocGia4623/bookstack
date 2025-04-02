package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type OrderService interface {
	CreateOrder(request.OrderRequest, int) (models.Order, error)
	CancelOrder(int) error
	GetOrder(userID int) (models.Order, error)
	GetUserOrder(orderId int) ([]models.Order, error)
}

type OrderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderServiceImpl(repository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		repo: repository,
	}
}

func (o *OrderServiceImpl) GetUserOrder(userId int) ([]models.Order, error) {
	return o.repo.GetUserOrder(userId)
}

func (o *OrderServiceImpl) GetOrder(orderId int) (models.Order, error) {
	return o.repo.GetOrder(orderId)
}

func (o *OrderServiceImpl) CreateOrder(request request.OrderRequest, userId int) (models.Order, error) {
	return o.repo.CreateOrder(request, userId)
}

func (o *OrderServiceImpl) CancelOrder(orderId int) error {
	return o.repo.CancelOrder(orderId)
}
