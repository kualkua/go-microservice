# Architecture and Repository Structure

## Layer Principles

- **`domain` depends on nothing**: models and interfaces only, with no concrete technologies.
- **`app` depends only on `domain`**: use-cases, no transport code and no concrete adapters; input models are transport-agnostic.
- **`api` depends on `app` and `domain`**: transport code (HTTP/gRPC/Telegram/consumer), validation/mapping, and use-case calls.
- **`adapters` implement `domain` interfaces**: persistence, queues, external clients; no use-case logic.

## Target Directory Structure

```text
src/
  cmd/
    main-app/
      main.go                 # dependency wiring and application startup

  config/
    config.go                 # process configuration: env and settings loading

  internal/
    adapters/
      <tool>/                 # e.g. filestore, postgres, nats: domain interface implementations

    api/
      http/
        <module>/
          builder.go
          handler.go
          transform.go        # api <-> domain mapping, and api <-> app DTO mapping if needed

      grpc/
        <module>/
          service.go
          transform.go

      consumer/
        <module>/
          handler.go

      telegram/
        builder.go            # Telegram adapter wiring
        handler.go            # main loop and update parsing
        client.go             # Telegram Bot API calls
        transform.go          # Telegram User/Message -> app input mapping
        keyboards.go          # menu representation, may use domain data
        types.go              # API update/response DTOs

    app/
      <feature>/
        usecase.go            # scenario; only domain types at the boundary

    domain/
      <bounded-context>/
        model.go
        repository.go
        port.go

  pkg/
    ...                       # shared utilities without domain logic, when needed
```

Transports are added as **parallel** trees under `internal/api/`, without mixing use-case logic into adapters.
