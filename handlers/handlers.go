package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/MXLange/rinha-only-go/entities"
	"github.com/MXLange/rinha-only-go/repository"
	"github.com/MXLange/rinha-only-go/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	channel      chan *entities.Payment
	repository   *repository.MemoryRepository
	fetch        *services.Fetch
	apiInstances []string
}

func NewHandler(channel chan *entities.Payment, repo *repository.MemoryRepository, fetch *services.Fetch, apiInstances []string) (*Handler, error) {
	if channel == nil {
		return nil, errors.New("channel is required")
	}

	if repo == nil {
		return nil, errors.New("repository is required")
	}

	if fetch == nil {
		return nil, errors.New("fetch service is required")
	}

	return &Handler{
		channel:      channel,
		repository:   repo,
		fetch:        fetch,
		apiInstances: apiInstances,
	}, nil
}

func (h *Handler) NewPayment(c *fiber.Ctx) error {
	var payment *entities.Payment = new(entities.Payment)
	if err := c.BodyParser(payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if payment.ID == "" || payment.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	h.channel <- payment

	return c.Status(fiber.StatusAccepted).JSON(payment)
}

func (h *Handler) GetSummary(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	internal := c.Query("internal")

	fmt.Println("Fetching summary with parameters:", "from:", from, "to:", to)

	dateFormat := "2006-01-02T15:04:05.000Z"

	var fromTime, toTime *time.Time = nil, nil

	if from != "" {
		parsedFrom, err := time.Parse(dateFormat, from)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'from' date format"})
		}
		fromTime = &parsedFrom
	}

	if to != "" {
		parsedTo, err := time.Parse(dateFormat, to)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'to' date format"})
		}
		toTime = &parsedTo
	}

	summary := h.repository.GetSummary(fromTime, toTime)

	fmt.Println("Summary fetched:", summary)
	if internal == "" {
		for _, instance := range h.apiInstances {

			if instance == "" {
				continue
			}

			instanceSummary, err := h.fetch.GetInstanceSummary(instance, from, to)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch instance summary"})
			}

			summary.Default.TotalRequests += instanceSummary.Default.TotalRequests
			summary.Default.TotalAmount += instanceSummary.Default.TotalAmount
			summary.Fallback.TotalRequests += instanceSummary.Fallback.TotalRequests
			summary.Fallback.TotalAmount += instanceSummary.Fallback.TotalAmount
		}
	}

	summary.Default.TotalAmount = float64(int(summary.Default.TotalAmount*100)) / 100
	summary.Fallback.TotalAmount = float64(int(summary.Fallback.TotalAmount*100)) / 100

	return c.JSON(summary)
}
