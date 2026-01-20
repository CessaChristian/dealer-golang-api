package report

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"dealer_golang_api/internal/service/favorite"
	"dealer_golang_api/internal/service/vehicle"
)

type Service struct {
	VehicleSvc  *vehicle.Service
	FavoriteSvc *favorite.Service
}

func NewService(vs *vehicle.Service, fs *favorite.Service) *Service {
	return &Service{
		VehicleSvc:  vs,
		FavoriteSvc: fs,
	}
}

// ========== LOW STOCK REPORT ==========

func (s *Service) LowStockJSON(ctx context.Context) ([]vehicle.VehicleLowStock, error) {
	return s.VehicleSvc.LowStock()
}

func (s *Service) LowStockCSV(ctx context.Context) (string, error) {
	list, err := s.VehicleSvc.LowStock()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	w := csv.NewWriter(&b)

	// Header
	w.Write([]string{"vehicle_id", "name", "brand", "type", "stock", "price"})

	// Rows
	for _, v := range list {
		w.Write([]string{
			fmt.Sprint(v.ID),
			v.Name,
			v.Brand,
			v.Type,
			fmt.Sprint(v.Stock),
			fmt.Sprintf("%.2f", v.Price),
		})
	}

	w.Flush()
	return b.String(), nil
}

// ========== FAVORITE REPORT ==========

func (s *Service) FavoriteJSON(ctx context.Context) ([]favorite.FavoriteReport, error) {
	return s.FavoriteSvc.GetAllFavoritesAdmin()
}

func (s *Service) FavoriteCSV(ctx context.Context) (string, error) {
	list, err := s.FavoriteSvc.GetAllFavoritesAdmin()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	w := csv.NewWriter(&b)

	// Header
	w.Write([]string{"user", "vehicle", "brand", "type", "price"})

	// Rows
	for _, f := range list {
		w.Write([]string{
			f.Name,
			f.VehicleName,
			f.Brand,
			f.Type,
			fmt.Sprintf("%.2f", f.Price),
		})
	}

	w.Flush()
	return b.String(), nil
}
