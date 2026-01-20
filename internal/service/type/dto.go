package vtype

type CreateTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type TypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
