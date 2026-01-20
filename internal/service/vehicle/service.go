package vehicle

import (
	"context"
	"dealer_golang_api/utils"

	"dealer_golang_api/internal/service/brand"
	vtype "dealer_golang_api/internal/service/type"
)

type Service struct {
	Repo      VehicleRepository
	BrandRepo brand.Repository
	TypeRepo  vtype.Repository
}

func NewService(repo VehicleRepository, b brand.Repository, t vtype.Repository) *Service {
	return &Service{Repo: repo, BrandRepo: b, TypeRepo: t}
}

func (s *Service) Create(req CreateVehicleRequest) error {

	brandID, err := s.BrandRepo.Ensure(context.Background(), req.Brand)
	if err != nil {
		return err
	}

	typeID, err := s.TypeRepo.Ensure(context.Background(), req.Type)
	if err != nil {
		return err
	}

	v := Vehicle{
		Name:         req.Name,
		BrandID:      brandID,
		TypeID:       typeID,
		FuelType:     utils.NormalizeFuelType(req.FuelType),
		Transmission: utils.NormalizeTransmission(req.Transmission),
		Price:        req.Price,
		Stock:        req.Stock,
	}

	return s.Repo.Create(context.Background(), v)
}

func (s *Service) Update(id int, req UpdateVehicleRequest) error {

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.FuelType != nil {
		updates["fuel_type"] = utils.NormalizeFuelType(*req.FuelType)
	}
	if req.Transmission != nil {
		updates["transmission"] = utils.NormalizeTransmission(*req.Transmission)
	}
	if req.Brand != nil {
		id, _ := s.BrandRepo.Ensure(context.Background(), *req.Brand)
		updates["brand_id"] = id
	}
	if req.Type != nil {
		id, _ := s.TypeRepo.Ensure(context.Background(), *req.Type)
		updates["type_id"] = id
	}

	return s.Repo.UpdatePartial(context.Background(), id, updates)
}

func (s *Service) ImportCar(req ImportVehicleRequest) error {
	car, err := FetchCarSpec(req.Make, req.Model)
	if err != nil {
		return err
	}

	brandID, err := s.BrandRepo.Ensure(context.Background(), req.Brand)
	if err != nil {
		return err
	}

	typeID, err := s.TypeRepo.Ensure(context.Background(), req.Type)
	if err != nil {
		return err
	}

	v := Vehicle{
		Name:         car.Make + " " + car.Model,
		BrandID:      brandID,
		TypeID:       typeID,
		FuelType:     utils.NormalizeFuelType(car.FuelType),
		Transmission: utils.NormalizeTransmission(car.Transmission),
		Price:        0,
		Stock:        0,
	}

	return s.Repo.Create(context.Background(), v)
}

func (s *Service) GetAll() ([]VehicleResponse, error) {
	return s.Repo.GetAll(context.Background())
}

func (s *Service) GetByID(id int) (VehicleResponse, error) {
	return s.Repo.GetByID(context.Background(), id)
}

func (s *Service) LowStock() ([]VehicleLowStock, error) {
	return s.Repo.GetLowStock(context.Background())
}
