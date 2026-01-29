package writer

import (
	"errors"
	"github.com/google/uuid"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"strings"
)

import (
	"fmt"
)

func (h *WriteHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
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

	err := h.ProductService.DeleteProduct(r.Context(), productID)
	if err != nil {
		fmt.Printf("[ERROR] - Error deleting product: %s\n", err.Error())
		if errors.Is(domain.ErrProductNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(fmt.Sprintf("error deleting product: %s", err.Error())))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("error deleting product")))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
