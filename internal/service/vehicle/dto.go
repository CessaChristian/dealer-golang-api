package vehicle

type CreateVehicleRequest struct {
	Name         string  `json:"name" validate:"required,min=2"`
	Brand        string  `json:"brand" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	FuelType     string  `json:"fuel_type" validate:"required"`
	Transmission string  `json:"transmission" validate:"required"`
	Price        float64 `json:"price" validate:"gt=0"`
	Stock        int     `json:"stock" validate:"gte=0"`
}

type UpdateVehicleRequest struct {
	Name         *string  `json:"name" validate:"omitempty,min=2"`
	Brand        *string  `json:"brand"`
	Type         *string  `json:"type"`
	FuelType     *string  `json:"fuel_type"`
	Transmission *string  `json:"transmission"`
	Price        *float64 `json:"price" validate:"omitempty,gt=0"`
	Stock        *int     `json:"stock" validate:"omitempty,gte=0"`
}

type VehicleResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Brand        string  `json:"brand"`
	Type         string  `json:"type"`
	FuelType     string  `json:"fuel_type"`
	Transmission string  `json:"transmission"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
}

type CarSpecResponse struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	FuelType     string `json:"fuel_type"`
	Transmission string `json:"transmission"`
	Class        string `json:"class"`
}

type ImportVehicleRequest struct {
	Brand string `json:"brand" validate:"required"`
	Type  string `json:"type" validate:"required"`
	Make  string `json:"make" validate:"required"`
	Model string `json:"model" validate:"required"`
}
