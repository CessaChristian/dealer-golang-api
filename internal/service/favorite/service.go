package favorite

import "context"

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) AddFavorite(userID, vehicleID int) error {
	f := Favorite{
		UserID:    userID,
		VehicleID: vehicleID,
	}

	return s.Repo.AddFavorite(context.Background(), f)
}

func (s *Service) GetFavorites(userID int) ([]FavoriteReport, error) {
	return s.Repo.GetFavoritesByUser(context.Background(), userID)
}

func (s *Service) GetAllFavoritesAdmin() ([]FavoriteReport, error) {
	return s.Repo.GetAllFavoritesAdmin(context.Background())
}
