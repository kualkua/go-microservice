package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"payment-microservices/src/config"
	"payment-microservices/src/internal/adapters/logger"
	"payment-microservices/src/internal/adapters/memory"
	"payment-microservices/src/internal/adapters/rabbitmq"
	consumernotification "payment-microservices/src/internal/api/consumer/notification"
	consumerpayment "payment-microservices/src/internal/api/consumer/payment"
	httppayment "payment-microservices/src/internal/api/http/payment"
	appnotification "payment-microservices/src/internal/app/notification"
	apppayment "payment-microservices/src/internal/app/payment"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application stopped with error: %v", err)
	}
}

func run() error {
	cfg := config.Load()

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	paymentPublisher := rabbitmq.NewPaymentPublisher(cfg, conn)
	createPayment := apppayment.NewCreatePaymentUseCase(paymentPublisher)
	httpPaymentHandler := httppayment.NewHandler(createPayment)

	mux := http.NewServeMux()
	mux.HandleFunc("/payments", httpPaymentHandler.CreatePayment)
	server := &http.Server{
		Addr:    cfg.APIAddress,
		Handler: mux,
	}

	paymentRepository := memory.NewPaymentRepository()
	processedEventPublisher := rabbitmq.NewProcessedEventPublisher(cfg, conn)
	processPayment := apppayment.NewProcessPaymentUseCase(paymentRepository, processedEventPublisher)
	paymentConsumer := consumerpayment.NewHandler(cfg, conn, processPayment)

	paymentNotifier := logger.NewPaymentProcessedNotifier()
	notifyPaymentProcessed := appnotification.NewNotifyPaymentProcessedUseCase(paymentNotifier)
	notificationConsumer := consumernotification.NewHandler(cfg, conn, notifyPaymentProcessed)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 3)
	go startHTTPServer(server, errCh)
	go startPaymentsConsumer(ctx, paymentConsumer, errCh)
	go startNotificationConsumer(ctx, notificationConsumer, errCh)

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-errCh:
		stop()
		shutdownHTTPServer(server)
		return err
	}

	shutdownHTTPServer(server)
	return nil
}

func startHTTPServer(server *http.Server, errCh chan<- error) {
	log.Printf("HTTP API is running on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- err
	}
}

func startPaymentsConsumer(ctx context.Context, handler *consumerpayment.Handler, errCh chan<- error) {
	if err := handler.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		errCh <- err
	}
}

func startNotificationConsumer(ctx context.Context, handler *consumernotification.Handler, errCh chan<- error) {
	if err := handler.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		errCh <- err
	}
}

func shutdownHTTPServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shut down HTTP server: %v", err)
	}
}
