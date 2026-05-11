# Payment Microservices

An educational Go project for learning the basic structure of a microservice application and message-based communication through RabbitMQ.

The project demonstrates a simple payment flow:

1. `api-gateway` receives an HTTP request to create a payment.
2. The gateway publishes a message to the RabbitMQ `payments` queue.
3. `payments-service` consumes the message from the `payments` queue, processes the payment, and publishes an event to the `payment_processed` queue.
4. `notification-service` consumes the event from the `payment_processed` queue and logs the notification.

## Project Structure

```text
docs/
  architecture.md         # source of truth for repository architecture

src/
  cmd/
    api-gateway/          # HTTP API process
    payments-service/     # payment consumer process
    notification-service/ # notification consumer process

  config/                 # process configuration

  internal/
    domain/               # models and interfaces, no infrastructure dependencies
      payment/
      notification/

    app/                  # use-cases, depends only on domain
      payment/
      notification/

    api/                  # transports: HTTP handlers and queue consumers
      http/
      consumer/

    adapters/             # implementations of domain interfaces
      rabbitmq/
      memory/
      logger/

requests/
  req.http                # example HTTP request
```

## Requirements

- Go 1.24+
- Docker and Docker Compose

## Configuration

By default, the application uses:

```env
API_GATEWAY_ADDR=:8080
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
PAYMENTS_QUEUE=payments
PAYMENT_PROCESSED_QUEUE=payment_processed
MESSAGE_LIMIT=100
```

An example is available in `.env.example`.

## Running the Project

### 1. Start RabbitMQ

```bash
docker compose up -d
```

The RabbitMQ Management UI will be available at:

```text
http://localhost:15672
```

Login and password:

```text
guest / guest
```

### 2. Start the API Gateway

In a separate terminal:

```bash
go run ./src/cmd/api-gateway
```

The service listens for HTTP requests on port `8080`.

### 3. Start the Payments Service

In a separate terminal:

```bash
go run ./src/cmd/payments-service
```

The service listens to the `payments` queue.

### 4. Start the Notification Service

In a separate terminal:

```bash
go run ./src/cmd/notification-service
```

The service listens to the `payment_processed` queue.

## Verification

Send a request:

```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "amount": 150.75,
    "currency": "USD",
    "idempotency_key": "txn-001"
  }'
```

Expected response:

```json
{
  "success": true,
  "message": "Payment created"
}
```

After that:

- `api-gateway` logs the received HTTP request;
- `payments-service` logs the payment processing step;
- `notification-service` logs the processed payment event.

## Build

```bash
go build ./src/cmd/api-gateway
go build ./src/cmd/payments-service
go build ./src/cmd/notification-service
```

## Tests

```bash
go test ./...
```

The project currently compiles, but there are no dedicated unit tests yet.

## Possible Improvements

- add unit tests for payment validation;
- replace the temporary repository with real storage;
- implement idempotency using `idempotency_key`;
- add graceful shutdown;
- add retry and dead-letter queue handling for message processing errors;
- move ports and queue names to configuration.
