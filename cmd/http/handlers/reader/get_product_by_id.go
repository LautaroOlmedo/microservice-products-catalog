package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"strings"
)

func (h *ReaderHandler) HandleGetProductByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		//http.Error(w, "invalid product id", http.StatusBadRequest)
		//return
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(fmt.Sprintf("error reading body")))
		if err != nil {
			return
		}
		return
	}

	productID := parts[len(parts)-1]
	if productID == "" {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(productID); err != nil {
		http.Error(w, "invalid product id format, must be UUID", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(domain.ErrProductNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(fmt.Sprintf("error fetching product: %s", err.Error())))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("error fetching product: %s", err.Error())))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	productsResponse, err := json.Marshal(product)
	if err != nil {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}

	_, err = w.Write(productsResponse)
	if err != nil {
		return
	}
}
