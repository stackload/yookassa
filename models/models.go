package models

import (
	"time"
)

type Payment struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	Paid          bool                   `json:"paid"`
	Amount        Amount                 `json:"amount"`
	Confirmation  Confirmation           `json:"confirmation"`
	CreatedAt     time.Time              `json:"created_at"`
	Description   string                 `json:"description,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	PaymentMethod *PaymentMethod         `json:"payment_method,omitempty"`
	Recipient     *Recipient             `json:"recipient,omitempty"`
	Refundable    bool                   `json:"refundable"`
	Test          bool                   `json:"test"`
}

type PaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	Capture      bool         `json:"capture"`
	Description  string       `json:"description,omitempty"`
	Metadata     interface{}  `json:"metadata,omitempty"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type            string `json:"type"`
	ReturnURL       string `json:"return_url,omitempty"`
	ConfirmationURL string `json:"confirmation_url,omitempty"`
}

// Дополнительные структуры из документации ответа
type PaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}
