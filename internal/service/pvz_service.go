package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"pvz-service/internal/domain"
	"pvz-service/internal/repository"
)

type PVZService interface {
	CreatePVZ(city string) (domain.PVZ, error)
	GetPVZByID(id string) (domain.PVZ, error)
	ListPVZsWithRelations(startDate, endDate string, page, limit int) ([]repository.PVZResponse, error)
}

type PVZServiceImpl struct {
	pvzRepo repository.PVZRepository
}

func NewPVZService(pvzRepo repository.PVZRepository) *PVZServiceImpl {
	return &PVZServiceImpl{pvzRepo: pvzRepo}
}

func (p *PVZServiceImpl) CreatePVZ(city string) (domain.PVZ, error) {
	allowedCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}

	if !allowedCities[city] {
		return domain.PVZ{}, errors.New("invalid city")
	}

	return p.pvzRepo.CreatePVZ(city, uuid.New)
}

func (p *PVZServiceImpl) GetPVZByID(id string) (domain.PVZ, error) {
	return p.pvzRepo.GetPVZByID(id)
}

func (p *PVZServiceImpl) ListPVZsWithRelations(startDate, endDate string, page, limit int) ([]repository.PVZResponse, error) {
	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return nil, errors.New("invalid start date format")
		}
	}

	if endDate != "" {
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
	}

	if page < 1 {
		return nil, errors.New("invalid page number")
	}

	if limit < 1 || limit > 30 {
		return nil, errors.New("invalid limit")
	}

	offset := (page - 1) * limit
	return p.pvzRepo.ListPVZsWithRelations(start, end, limit, offset)
}
