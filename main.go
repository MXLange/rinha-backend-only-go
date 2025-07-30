package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/MXLange/rinha-only-go/entities"
	"github.com/MXLange/rinha-only-go/handlers"
	"github.com/MXLange/rinha-only-go/repository"
	"github.com/MXLange/rinha-only-go/services"
	"github.com/gofiber/fiber/v2"
)

func main() {

	instancesStr := os.Getenv("API_INSTANCES")
	instances := strings.Split(instancesStr, ",")

	baseUrlDefault := os.Getenv("BASE_URL_DEFAULT")
	baseUrlFallback := os.Getenv("BASE_URL_FALLBACK")
	workersStr := os.Getenv("WORKERS")
	workers, err := strconv.Atoi(workersStr)
	if err != nil || workers <= 0 {
		workers = 10 // Default to 10 workers if parsing fails
	}

	fetch, err := services.NewFetch(baseUrlDefault, baseUrlFallback)
	if err != nil {
		panic(err)
	}

	paymentChannel := make(chan *entities.Payment, 100000)
	repo, getMutex := repository.NewMemoryRepository()
	handler, err := handlers.NewHandler(paymentChannel, repo, fetch, instances)
	if err != nil {
		panic(err)
	}

	worker, err := services.NewWorker(paymentChannel, repo, fetch, 10, getMutex)
	if err != nil {
		panic(err)
	}

	worker.Start()

	app := fiber.New()
	app.Post("/payments", handler.NewPayment)
	app.Get("/payments-summary", handler.GetSummary)

	app.Listen(":8080")
}
