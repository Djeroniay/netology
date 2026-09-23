package main

import "sync"

type MemoryRepository struct {
	mu     sync.Mutex
	orders []Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: make([]Order, 0),
	}
}

func (r *MemoryRepository) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders = make([]Order, 0)
	return nil
}

func (r *MemoryRepository) Save(order Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = len(r.orders) + 1
	r.orders = append(r.orders, order)

	return nil
}

func (r *MemoryRepository) Orders() []Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]Order, len(r.orders))
	copy(result, r.orders)

	return result
}
