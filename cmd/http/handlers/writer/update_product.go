package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"microservice-products-catalog/cmd/http/dto"
	"microservice-products-catalog/internal/domain"
	"net/http"
	"strings"
)

func (h *WriteHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	var body dto.UpdateProductRequest
	if err := json.Unmarshal(bytes, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("DTO validation error: %s", err)))
		return
	}

	// TODO [technical debate] Create a mapper to parse data
	product := &domain.Product{}
	product.ID = productID
	if body.Name != nil {
		product.Name = *body.Name
	}

	if body.Description != nil {
		product.Description = *body.Description
	}

	if body.Price != nil {
		product.Price = *body.Price
	}

	if body.Stock != nil {
		product.Stock = *body.Stock
	}

	err = h.ProductService.UpdateProduct(r.Context(), product)
	if err != nil {
		fmt.Printf("[ERROR] - Error updating product: %s\n", err.Error())
		if errors.Is(domain.ErrProductNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(fmt.Sprintf("error updating product: %s", err.Error())))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("error updating product")))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
