package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
	"context"
	"strconv"

	"github.com/plutov/paypal/v4"
)

type OrderService interface {
	CreateOrder(request.OrderRequest, int) (models.Order, error)
	CancelOrder(int) error
	GetOrder(userID int) (models.Order, error)
	GetUserOrder(orderId int) ([]models.Order, error)
	CreatePaypalOrder(*paypal.Client, int) (*paypal.Order, error)
	UpdateOrderStatus(webhookPayload map[string]interface{}) error
}

type OrderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderServiceImpl(repository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		repo: repository,
	}
}

func (o *OrderServiceImpl) UpdateOrderStatus(webhookPayload map[string]interface{}) error {
	return o.repo.UpdateOrderStatus(webhookPayload)
}

func (o *OrderServiceImpl) CreatePaypalOrder(c *paypal.Client, orderId int) (*paypal.Order, error) {
	order, err := o.repo.GetOrder(orderId)
	if err != nil {
		return nil, err
	}

	ord, err := c.CreateOrder(context.Background(), paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "USD",
				Value:    strconv.FormatFloat(order.TotalPrice, 'f', -1, 64),
			},
		},
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	return ord, nil
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
