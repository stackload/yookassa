package yookassa

import (
	"context"
	"time"

	"github.com/stackload/yookassa/client"
	"github.com/stackload/yookassa/models"
)

type PaymentService struct {
	yookassa *client.Client
}

func NewPaymentService() (*PaymentService, error) {
	client, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &PaymentService{
		yookassa: client,
	}, nil
}

// CreatePayment создает платеж и возвращает URL для оплаты
func (p *PaymentService) SendInvoice(ctx context.Context, amount float64, description string, metadata map[string]interface{}) (*models.Payment, error) {
	payment, err := p.yookassa.CreateInvoice(ctx, amount, description, metadata)
	return payment, err
}

func (p *PaymentService) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	payment, err := p.yookassa.GetPayment(ctx, paymentID)
	return payment, err
}

func (p *PaymentService) CheckPayment(ctx context.Context, paymentID string, timeout time.Duration) error {
	err := p.yookassa.WaitForPaymentComplete(ctx, paymentID, timeout)
	return err
}
