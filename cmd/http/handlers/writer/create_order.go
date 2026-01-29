package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"microservice-products-catalog/cmd/http/dto"
	"microservice-products-catalog/internal/domain"
	"net/http"
)

func (h *WriteHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error rading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	var body dto.CreateOrderRequest
	if err := json.Unmarshal(bytes, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error rading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	// TODO [technical debate] Create a mapper to parse data
	err = h.OrderService.CreateOrder(r.Context(), body.ProductID, body.Quantity)

	if err != nil {
		fmt.Printf("[ERROR] - Error creating order: %s\n", err.Error())
		if errors.Is(domain.ErrProductNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte(fmt.Sprintf("error creating order: %s", err)))
			if err != nil {
				return
			}
			return
		}
		if errors.Is(domain.ErrInsufficientStock, err) {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(fmt.Sprintf("error creating order: %s", err)))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error creating order")))
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
