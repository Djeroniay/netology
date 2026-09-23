package main

import "fmt"

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (s *EmailSender) Send(customer string, order Order) error {
	fmt.Printf(
		"Email: заказ клиента %s сохранён. Сумма: %.2f\n",
		customer,
		order.Total,
	)

	return nil
}

type SMSSender struct{}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

func (s *SMSSender) Send(customer string, order Order) error {
	fmt.Printf(
		"SMS: заказ клиента %s сохранён. Сумма: %.2f\n",
		customer,
		order.Total,
	)

	return nil
}
