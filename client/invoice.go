package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/stackload/yookassa/models"
)

// CreateInvoice создает платеж
func (c *Client) CreateInvoice(ctx context.Context, amount float64, description string, metadata map[string]interface{}) (*models.Payment, error) {
	amountStr := fmt.Sprintf("%.2f", amount)
	_ = godotenv.Load()
	returnURL := os.Getenv("YOOKASSA_RETURN_URL")

	paymentRequest := models.PaymentRequest{
		Amount: models.Amount{
			Value:    amountStr,
			Currency: "RUB",
		},
		Confirmation: models.Confirmation{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Capture:     true,
		Description: description,
		Metadata:    metadata,
	}

	payment, err := c.createPayment(ctx, paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return payment, nil
}

// CreatePayment создает платеж в ЮKassa
func (c *Client) createPayment(ctx context.Context, request models.PaymentRequest) (*models.Payment, error) {
	paymentURL := fmt.Sprintf("%s/payments", c.baseURL)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", paymentURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Idempotence-Key", generateUUID())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", string(body))
	}

	var payment models.Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &payment, nil
}

// GetPayment получает информацию о платеже
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	paymentURL := fmt.Sprintf("%s/payments/%s", c.baseURL, paymentID)

	req, err := http.NewRequestWithContext(ctx, "GET", paymentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", string(body))
	}

	var payment models.Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &payment, nil
}

// Функция для ожидания завершения платежа с таймаутом
func (c *Client) WaitForPaymentComplete(ctx context.Context, paymentID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second) // Проверяем каждые 5 секунд
	defer ticker.Stop()

	fmt.Printf("Ожидание оплаты (таймаут: %v)...\n", timeout)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("таймаут ожидания платежа")

		case <-ticker.C:
			payment, err := c.GetPayment(ctx, paymentID)
			if err != nil {
				return fmt.Errorf("ошибка при проверке платежа: %w", err)
			}

			fmt.Printf("Текущий статус: %s\n", payment.Status)

			// Проверяем завершенные статусы
			switch payment.Status {
			case "succeeded":
				fmt.Println("✅ Платеж успешно завершен!")
				return nil
			case "canceled":
				return fmt.Errorf("платеж отменен")
			case "waiting_for_capture":
				fmt.Println("⚠️ Платеж ожидает подтверждения (capture)")
				return nil
				// Статус "pending" продолжаем ждать
			}
		}
	}
}

func generateUUID() string {
	return uuid.NewString()
}
