# Payment Microservices

An educational Go project for learning layered Go application structure and message-based communication through RabbitMQ.

The project demonstrates a simple payment flow:

1. The HTTP API receives a request to create a payment.
2. The payment API use-case publishes a message to the RabbitMQ `payments` queue.
3. The payment consumer reads from `payments`, processes the payment, and publishes a processed event to `payment_processed`.
4. The notification consumer reads from `payment_processed` and logs the notification.

The app is intentionally small, but the repository is organized around clean layer boundaries.

## Project Structure

```text
src/
  cmd/
    main-app/
      main.go                 # dependency wiring and application startup

  config/
    config.go                 # process configuration: env and settings loading

  internal/
    adapters/                 # implementations of domain interfaces
      rabbitmq/               # queue publishers and RabbitMQ connection
      memory/                 # in-memory payment repository
      logger/                 # notification logger adapter

    api/                      # entrypoints: HTTP handlers and queue consumers
      consumer/
        payment/
          handler.go
          transform.go
        notification/
          handler.go
          transform.go

      http/
        payment/
          builder.go
          handler.go
          transform.go        # api <-> domain mapping

      grpc/                   # place for future gRPC transports

    app/                      # use-cases, depends only on domain
      payment/
        usecase.go
      notification/
        usecase.go

    domain/                   # business models and interfaces
      payment/
        model.go
        repository.go
        port.go
      notification/
        port.go

  pkg/                        # independent shared helpers, when needed

requests/
  req.http                    # example HTTP request
```

## Layer Rules

- `domain` contains business entities and interfaces only.
- `app` contains use-cases and depends only on `domain`.
- `api` contains transport-specific code: HTTP handlers, consumers, validation, and mapping.
- `adapters` implement `domain` interfaces using concrete tools such as RabbitMQ, databases, loggers, or external clients.
- `cmd/main-app` wires dependencies together and starts the application.

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

RabbitMQ Management UI:

```text
http://localhost:15672
```

Login and password:

```text
guest / guest
```

### 2. Start the Application

```bash
go run ./src/cmd/main-app
```

The single `main-app` process starts:

- HTTP API on `:8080`;
- payment consumer for the `payments` queue;
- notification consumer for the `payment_processed` queue.

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

After that, the application logs the payment request, payment processing, and notification event.

## Build

```bash
go build ./src/cmd/main-app
```

## Tests

```bash
go test ./...
```

The project currently compiles, but there are no dedicated unit tests yet.

## Possible Improvements

- add unit tests for payment validation and use-cases;
- replace the temporary memory repository with real storage;
- implement idempotency using `idempotency_key`;
- add retry and dead-letter queue handling for message processing errors;
- add gRPC transport under `internal/api/grpc`.
