package client

import "errors"

var (
	ErrorEnvParametr = errors.New("YOOKASSA_SHOP_ID or YOOKASSA_SECRET_KEY not set")
	ErrorNewClient   = errors.New("new client error")
)
