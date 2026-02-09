package vtype

type CreateTypeRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type TypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
