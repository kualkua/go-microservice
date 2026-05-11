package payment

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Received payment request: %+v", req)

	if err := validateCreatePaymentRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.createPayment.Execute(r.Context(), req.toDomain()); err != nil {
		log.Printf("Failed to create payment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(paymentResponse{
		Success: true,
		Message: "Payment created",
	}); err != nil {
		log.Printf("Failed to encode payment response: %v", err)
	}
}
