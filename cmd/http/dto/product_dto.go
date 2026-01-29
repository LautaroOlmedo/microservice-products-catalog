package dto

// CreateProductRequest Ensure to add the necessaries validations to DTO.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0.1"`
	Stock       int     `json:"stock" validate:"required,min=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0.1"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
}
