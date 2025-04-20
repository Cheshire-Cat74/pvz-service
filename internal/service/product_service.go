package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"

	"pvz-service/internal/domain"
)

type ProductService interface {
	AddProduct(receptionID string, productType string, idGenerator func() uuid.UUID) (string, error)
	GetProductByID(id string) (domain.Product, error)
	GetLastProduct(receptionID string) (domain.Product, error)
	DeleteProduct(id string) error
}

type ReceptionRepository interface {
	GetOpenReception(pvzID string) (domain.Reception, error)
}

type ProductServiceImpl struct {
	productRepo   ProductService
	receptionRepo ReceptionRepository
}

func NewProductService(
	productRepo ProductService,
	receptionRepo ReceptionRepository,
) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
	}
}

func (p *ProductServiceImpl) AddProduct(pvzID, productType string) (domain.Product, error) {
	allowedTypes := map[string]bool{
		"электроника": true,
		"одежда":      true,
		"обувь":       true,
	}

	if !allowedTypes[productType] {
		return domain.Product{}, errors.New("invalid product type")
	}

	reception, err := p.receptionRepo.GetOpenReception(pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, errors.New("no open reception for this PVZ")
		}
		return domain.Product{}, errors.New("database error")
	}

	productID, err := p.productRepo.AddProduct(reception.ID, productType, uuid.New)
	if err != nil {
		return domain.Product{}, errors.New("failed to add product")
	}

	return p.productRepo.GetProductByID(productID)
}

func (p *ProductServiceImpl) DeleteLastProduct(pvzID string) error {
	reception, err := p.receptionRepo.GetOpenReception(pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no open reception for this PVZ")
		}
		return errors.New("database error")
	}

	product, err := p.productRepo.GetLastProduct(reception.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no products to delete in this reception")
		}
		return errors.New("database error")
	}

	return p.productRepo.DeleteProduct(product.ID)
}
