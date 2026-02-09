package brand

type CreateBrandRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type BrandResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
