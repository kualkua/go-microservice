package gateway

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"payment-microservices/internal/contracts"
	"payment-microservices/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	conn *amqp.Connection
}

func NewHandler(conn *amqp.Connection) *Handler {
	return &Handler{conn: conn}
}

func (h *Handler) CreatePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req contracts.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Received payment request: %+v", req)

	err = validateCreateHandler(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch, err := h.conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	_, err = rabbitmq.DeclarePaymentsQueue(ch)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ch.Publish(
		"",
		"payments",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := contracts.PaymentResponse{
		Success: true,
		Message: "Payment created",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func validateCreateHandler(req *contracts.PaymentRequest) error {
	if req.UserId == "" {
		return errors.New("user_id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if req.Currency == "" {
		return errors.New("currency is required")
	}
	if req.IdempotencyKey == "" {
		return errors.New("idempotency_key is required")
	}
	return nil
}
