package handler

import (
	"pvz-service/internal/handler/models"
	"pvz-service/internal/prometheus"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"pvz-service/internal/domain"
	"pvz-service/internal/service"
)

type PVZHandlers struct {
	pvzService service.PVZService
}

func NewPVZHandlers(pvzProcessor service.PVZService) *PVZHandlers {
	return &PVZHandlers{pvzService: pvzProcessor}
}

func (h *PVZHandlers) CreatePVZHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body domain.PVZ
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid request"})
		}

		pvz, err := h.pvzService.CreatePVZ(body.City)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: err.Error()})
		}

		prometheus.PickupPointsCreated.Inc()

		return c.Status(fiber.StatusCreated).JSON(pvz)

	}
}

func (h *PVZHandlers) GetPVZListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageStr := c.Query("page")
		if pageStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Message: "page parameter is required",
			})
		}

		limitStr := c.Query("limit")
		if limitStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Message: "limit parameter is required",
			})
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Message: "page must be a positive integer",
			})
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 30 {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Message: "limit must be between 1 and 30",
			})
		}

		startDate := c.Query("startDate")
		if startDate != "" {
			if _, err := time.Parse(time.RFC3339, startDate); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Message: "invalid startDate format, must be RFC3339",
				})
			}
		}

		endDate := c.Query("endDate")
		if endDate != "" {
			if _, err := time.Parse(time.RFC3339, endDate); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Message: "invalid endDate format, must be RFC3339",
				})
			}
		}

		result, err := h.pvzService.ListPVZsWithRelations(startDate, endDate, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(result)
	}
}
