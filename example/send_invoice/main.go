package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/stackload/yookassa/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := client.NewClient()
	if err != nil {
		log.Println("error", err)
	}

	//подготовка запроса и данных для создания платежа
	var (
		amount      float64 = 100
		description         = "test payment"
		user_id     int64   = 1
		metadata            = map[string]interface{}{
			"user_id": strconv.FormatInt(user_id, 10),
		}
	)
	//
	invoice, err := client.CreateInvoice(ctx, amount, description, metadata)
	//
	payment, err := client.GetPayment(ctx, invoice.ID)
	if err != nil {
		log.Println("error", err)
	}
	fmt.Printf("ConfirmationURL: %s\n", payment.Confirmation.ConfirmationURL)

}
