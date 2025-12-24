package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	authHeader string
	shopID     string
}

func NewClient() (*Client, error) {
	_ = godotenv.Load()
	shopID := os.Getenv("YOOKASSA_SHOP_ID")
	secretKey := os.Getenv("YOOKASSA_SECRET_KEY")
	baseURL := os.Getenv("YOOKASSA_URL")

	if baseURL == "" {
		baseURL = "https://api.yookassa.ru/v3"
	}

	if shopID == "" || secretKey == "" {
		return nil, ErrorEnvParametr
	}

	auth := fmt.Sprintf("%s:%s", shopID, secretKey)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		authHeader: fmt.Sprintf("Basic %s", encodedAuth),
		shopID:     shopID,
	}, nil
}
