package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stackload/yookassa"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := yookassa.NewPaymentService()
	if err != nil {
		log.Fatalf("Failed to create payment service: %v", err)
	}

	metadata := map[string]interface{}{
		"user_id": strconv.FormatInt(1, 10),
	}

	pay, err := service.SendInvoice(ctx, 100, "test", metadata)
	if err != nil {
		log.Fatalf("Failed to send invoice: %v", err)
	}
	fmt.Printf("Создан платеж: %s\n", pay.ID)
	fmt.Printf("URL для оплаты: %s\n", pay.Confirmation.ConfirmationURL)

	// Ждем 30 секунд перед проверкой статуса
	fmt.Println("Ожидание 5 секунд для проверки статуса...")
	time.Sleep(5 * time.Second)

	// Ждем, пока пользователь совершит платеж (периодическая проверка)
	err = service.CheckPayment(ctx, pay.ID, 60*time.Second)
	if err != nil {
		log.Printf("Ошибка при ожидании платежа: %v", err)
	}

}
