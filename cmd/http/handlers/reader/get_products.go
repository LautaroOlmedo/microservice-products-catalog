package reader

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"microservice-products-catalog/cmd/http/auth"
	"net/http"
	"strconv"
)

const defaultPaginationLimit = 10

func (h *ReaderHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}*/

	/*token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Header Authorization is required"))
		if err != nil {
			return
		}
		return
	}*/

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error parsing limit: %s", err)))
		if err != nil {
			return
		}
		return
	}

	// nextCursor := r.URL.Query().Get("next_cursor") // TODO implement pagination

	products, err := h.ProductService.GetProducts(r.Context(), limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error getting products: %s", err)))
		if err != nil {
			return
		}
		return
	}

	tok, err := h.TokenGenerator.Generate(
		r.Context(),
		auth.TokenClaims{
			Scope:     "products:read",
			RequestID: uuid.NewString(),
		},
	)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tok)
	w.WriteHeader(http.StatusOK)

	productsResponse, err := json.Marshal(products)
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

func parseLimit(request *http.Request) (int, error) {
	limitStr := request.URL.Query().Get("limit")

	if limitStr == "" {
		return defaultPaginationLimit, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("limit must be an integer: %w", err)
	}

	if limit <= 0 {
		return 0, fmt.Errorf("limit must be a positive number")
	}

	return limit, nil
}
