package service

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"pvz-service/internal/repository"
)

type AuthService interface {
	Register(email, password, role string) (string, error)
	Login(email, password string) (string, string, error)
	DummyLogin(role string) (string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type AuthServiceImpl struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{authRepo: authRepo}
}

func (p *AuthServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (p *AuthServiceImpl) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *AuthServiceImpl) Register(email, password, role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", errors.New("invalid role")
	}

	hashedPassword, err := p.HashPassword(password)
	if err != nil {
		return "", errors.New("failed to process password")
	}

	return p.authRepo.CreateUser(email, hashedPassword, role)
}

func (p *AuthServiceImpl) Login(email, password string) (string, string, error) {
	userID, hashedPassword, role, err := p.authRepo.FindUserByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if err := p.ComparePassword(hashedPassword, password); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	return userID, role, nil
}

func (p *AuthServiceImpl) DummyLogin(role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", errors.New("invalid role")
	}

	userID, err := p.authRepo.FindUserByRole(role)
	if errors.Is(err, sql.ErrNoRows) {
		hashedPassword, err := p.HashPassword("password")
		if err != nil {
			return "", errors.New("failed to create dummy user")
		}

		return p.authRepo.CreateUser("dummy@example.com", hashedPassword, role)
	}

	return userID, err
}
