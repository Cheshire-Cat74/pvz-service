package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"pvz-service/internal/handler/models"
	"strconv"
	"testing"
	"time"

	"pvz-service/internal/domain"
	"pvz-service/internal/repository"
)

type MockPVZService struct {
	mock.Mock
}

func (m *MockPVZService) CreatePVZ(city string) (domain.PVZ, error) {
	args := m.Called(city)
	return args.Get(0).(domain.PVZ), args.Error(1)
}

func (m *MockPVZService) GetPVZByID(id string) (domain.PVZ, error) {
	args := m.Called(id)
	return args.Get(0).(domain.PVZ), args.Error(1)
}

func (m *MockPVZService) ListPVZsWithRelations(
	startDate, endDate string, page, limit int) ([]repository.PVZResponse, error) {
	args := m.Called(startDate, endDate, page, limit)
	return args.Get(0).([]repository.PVZResponse), args.Error(1)
}

func TestPVZHandlers_CreatePVZHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockPVZService)
	handler := NewPVZHandlers(mockService)

	t.Run("success", func(t *testing.T) {
		expectedPVZ := domain.PVZ{
			ID:   uuid.NewString(),
			City: "Москва",
		}

		mockService.On("CreatePVZ", "Москва").Return(expectedPVZ, nil)

		app.Post("/pvz", handler.CreatePVZHandler())

		req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(`{"city":"Москва"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		app.Post("/pvz", handler.CreatePVZHandler())

		req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(`invalid`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestPVZHandlers_GetPVZListHandler(t *testing.T) {
	app := fiber.New()
	mockProcessor := new(MockPVZService)
	handler := NewPVZHandlers(mockProcessor)

	t.Run("success with UTC timezone", func(t *testing.T) {
		expected := []repository.PVZResponse{
			{
				PVZ: domain.PVZ{
					ID:   "pvz1",
					City: "Москва",
				},
			},
		}

		startDate := "2025-04-01T00:00:00Z"
		endDate := "2025-04-20T23:59:59Z"
		page := 1
		limit := 10

		mockProcessor.On("ListPVZsWithRelations", startDate, endDate, page, limit).
			Return(expected, nil)

		app.Get("/pvz", handler.GetPVZListHandler())

		url := "/pvz?startDate=" + startDate + "&endDate=" + endDate +
			"&page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit)
		req := httptest.NewRequest("GET", url, nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("success with dynamic time", func(t *testing.T) {
		expected := []repository.PVZResponse{
			{
				PVZ: domain.PVZ{
					ID:   "pvz1",
					City: "Москва",
				},
			},
		}

		now := time.Now().UTC()
		startDate := now.Format(time.RFC3339)
		endDate := now.Add(24 * time.Hour).Format(time.RFC3339)
		page := 1
		limit := 10

		mockProcessor.On("ListPVZsWithRelations", startDate, endDate, page, limit).
			Return(expected, nil)

		app.Get("/pvz", handler.GetPVZListHandler())

		url := "/pvz?startDate=" + startDate + "&endDate=" + endDate +
			"&page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit)
		req := httptest.NewRequest("GET", url, nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("success without dates", func(t *testing.T) {
		expected := []repository.PVZResponse{
			{
				PVZ: domain.PVZ{
					ID:   "pvz1",
					City: "Москва",
				},
			},
		}

		page := 1
		limit := 10

		mockProcessor.On("ListPVZsWithRelations", "", "", page, limit).
			Return(expected, nil)

		app.Get("/pvz", handler.GetPVZListHandler())

		url := "/pvz?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit)
		req := httptest.NewRequest("GET", url, nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("invalid page number (0)", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?page=0&limit=10", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid page number (negative)", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?page=-1&limit=10", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid limit (0)", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?page=1&limit=0", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid limit (too large)", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())

		req := httptest.NewRequest("GET", "/pvz?page=1&limit=31", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "limit must be between 1 and 30", errorResp.Message)
	})

	t.Run("valid maximum limit", func(t *testing.T) {
		expected := []repository.PVZResponse{}

		mockProcessor.On("ListPVZsWithRelations", "", "", 1, 30).
			Return(expected, nil)

		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?page=1&limit=30", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("missing page parameter", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?limit=10", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing limit parameter", func(t *testing.T) {
		app.Get("/pvz", handler.GetPVZListHandler())
		req := httptest.NewRequest("GET", "/pvz?page=1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
