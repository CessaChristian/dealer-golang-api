package favorite

type Favorite struct {
	UserID    int `db:"user_id"`
	VehicleID int `db:"vehicle_id"`
}

type FavoriteReport struct {
	Name    string  `json:"name"`
	VehicleName string  `json:"vehicle_name"`
	Brand       string  `json:"brand"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
}
