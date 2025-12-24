# Инструкция по работе с платежами YooKassa

## Описание
Данный код создает платеж через API YooKassa, генерирует ссылку для оплаты и отслеживает статус платежа.

## Предварительные требования

### 1. Установка Go
- Установите Go версии 1.19 или выше
- Проверьте установку: `go version`

### 2. Настройка YooKassa
1. Зарегистрируйтесь в [YooKassa](https://yookassa.ru/)
2. Получите `shopId` и `secretKey` в личном кабинете
3. Для тестирования используйте тестовый режим

### 3. Установка зависимостей
```bash
go get github.com/stackload/yookassa
```

## Конфигурация
### 1. Настройка переменных окружения (рекомендуемый способ)
Создайте файл .env в корне проекта:

```
YOOKASSA_SHOP_ID=ваш_shop_id
YOOKASSA_SECRET_KEY=ваш_secret_key
YOOKASSA_URL=https://api.yookassa.ru/v3
YOOKASSA_RETURN_URL=https://example.com
```

### 2. Инициализация
```
import "github.com/stackload/yookassa"

func main() {
    ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := yookassa.NewPaymentService()
	if err != nil {
		log.Fatalf("Failed to create payment service: %v", err)
	}
}
```
### 3. Создание платежа
```
// Базовый пример
metadata := map[string]interface{}{
    "user_id": "123",
    "order_id": "order_456",
    "description": "Оплата заказа",
}

payment, err := service.SendInvoice(
    ctx,
    100.00,        // сумма (в рублях)
    "Оплата тестового заказа", // описание
    metadata,      // дополнительные данные
)
```
### 4. Получение ссылки для оплаты
```
if payment.Confirmation != nil {
    paymentURL := payment.Confirmation.ConfirmationURL
    fmt.Printf("Ссылка для оплаты: %s\n", paymentURL)
    // Можно отправить пользователю или открыть в браузере
}
```

### 5. Проверка статуса платежа
##### Способ 1: Однократная проверка
```
// Получить информацию о платеже
paymentInfo, err := service.GetPayment(ctx, paymentID)
if err != nil {
    log.Printf("Ошибка получения статуса: %v", err)
} else {
    fmt.Printf("Статус платежа: %s\n", paymentInfo.Status)
}
```
##### Способ 2: Ожидание платежа в горутине (рекомендуется)
```
// Ожидание платежа с таймаутом
err := service.CheckPayment(
    ctx,                     // контекст
    payment.ID,              // ID платежа
    300*time.Second,         // время ожидания (5 минут)
)

if err != nil {
    log.Printf("Ошибка ожидания платежа: %v", err)
} else {
    fmt.Println("Платеж успешно завершен!")
}
```
