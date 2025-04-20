package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"

	"pvz-service/internal/domain"
	"pvz-service/internal/repository"
)

type ReceptionService interface {
	CreateReception(pvzID string) (domain.Reception, error)
	CloseLastReception(pvzID string) (domain.Reception, error)
}

type ReceptionServiceImpl struct {
	receptionRepo repository.ReceptionRepository
}

func NewReceptionService(receptionRepo repository.ReceptionRepository) *ReceptionServiceImpl {
	return &ReceptionServiceImpl{receptionRepo: receptionRepo}
}

func (p *ReceptionServiceImpl) CreateReception(pvzID string) (domain.Reception, error) {
	hasOpen, err := p.receptionRepo.HasOpenReception(pvzID)
	if err != nil {
		return domain.Reception{}, errors.New("database error")
	}
	if hasOpen {
		return domain.Reception{}, errors.New("open reception already exists for this PVZ")
	}

	receptionID, err := p.receptionRepo.CreateReception(pvzID, uuid.New)
	if err != nil {
		return domain.Reception{}, errors.New("failed to create reception")
	}

	return p.receptionRepo.GetReceptionByID(receptionID)
}

func (p *ReceptionServiceImpl) CloseLastReception(pvzID string) (domain.Reception, error) {
	reception, err := p.receptionRepo.GetOpenReception(pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Reception{}, errors.New("no open reception found for this PVZ")
		}
		return domain.Reception{}, errors.New("database error")
	}

	now := time.Now()
	if err := p.receptionRepo.CloseReception(reception.ID, now); err != nil {
		return domain.Reception{}, errors.New("failed to close reception")
	}

	reception.Status = "close"
	reception.ClosedAt = &now
	return reception, nil
}
