package main

import (
	"errors"
	"fmt"
)

type Order struct {
	ID       int
	Customer string
	Products []string
	Total    float64
	Status   string
}

// OrderService содержит только бизнес-логику создания заказа.
type OrderService struct {
	repository RepositoryWriter
	notifier   Notifier
}

func NewOrderService(
	repository RepositoryWriter,
	notifier Notifier,
) *OrderService {
	return &OrderService{
		repository: repository,
		notifier:   notifier,
	}
}

func (s *OrderService) CreateOrder(
	customer string,
	products []string,
	total float64,
) error {
	if customer == "" {
		return errors.New("customer is required")
	}

	if len(products) == 0 {
		return errors.New("at least one product is required")
	}

	if total <= 0 {
		return errors.New("total must be greater than zero")
	}

	order := Order{
		Customer: customer,
		Products: products,
		Total:    total,
		Status:   "pending",
	}

	if err := s.repository.Save(order); err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	if err := s.notifier.Send(customer, order); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	return nil
}
