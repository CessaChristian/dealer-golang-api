package favorite

type AddFavoriteRequest struct {
	VehicleID int `json:"vehicle_id" validate:"required"`
}
