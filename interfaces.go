package main

// RepositoryInitializer отвечает только за создание структуры хранилища.
type RepositoryInitializer interface {
	Init() error
}

// RepositoryWriter отвечает только за сохранение заказа.
type RepositoryWriter interface {
	Save(order Order) error
}

// Notifier отвечает за отправку уведомлений.
type Notifier interface {
	Send(customer string, order Order) error
}
