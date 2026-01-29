package writer

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"microservice-products-catalog/cmd/http/dto"
	"microservice-products-catalog/internal/domain"
	"net/http"
)

func (h *WriteHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}*/

	//JWToken := r.Header.Get("token")
	//if JWToken == "" {
	//} // handle token logic

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	var body dto.CreateProductRequest
	if err := json.Unmarshal(bytes, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	// TODO [tech debate] create a custom validate for DTO data coming on request
	/*validate := validator.New()
	if err := validate.Struct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("DTO validation error: %s", err)))
		return
	}*/

	// TODO [technical debate] Create a mapper to parse data
	err = h.ProductService.CreateProduct(r.Context(), domain.Product{
		ID:          uuid.New().String(),
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Stock:       body.Stock,
	})

	if err != nil {
		fmt.Printf("[ERROR] - Error creating product: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error creating product")))
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
