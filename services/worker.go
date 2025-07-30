package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/MXLange/rinha-only-go/entities"
	"github.com/MXLange/rinha-only-go/repository"
)

type Worker struct {
	channel     chan *entities.Payment
	concurrency int
	repository  *repository.MemoryRepository
	fetch       *Fetch
	getMutex    *sync.Mutex
}

func NewWorker(channel chan *entities.Payment, repo *repository.MemoryRepository, fetch *Fetch, concurrency int, getMutex *sync.Mutex) (*Worker, error) {

	if channel == nil {
		return nil, errors.New("channel is required")
	}

	if repo == nil {
		return nil, errors.New("repository is required")
	}

	if fetch == nil {
		return nil, errors.New("fetch service is required")
	}

	if concurrency <= 0 {
		return nil, errors.New("concurrency must be greater than zero")
	}

	if getMutex == nil {
		return nil, errors.New("getMutex is required")
	}

	return &Worker{
		channel:     channel,
		repository:  repo,
		fetch:       fetch,
		concurrency: concurrency,
		getMutex:    getMutex,
	}, nil
}

func (w *Worker) Start() {
	for i := 0; i < w.concurrency; i++ {
		go w.worker()
	}
}

func (w *Worker) worker() {
	for p := range w.channel {
		func() {
			w.getMutex.Lock()
			w.getMutex.Unlock()

			var err error

			p.RequestedAt = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

			if p.Err != "SAVE" {
				err, p.IsDefault = w.fetch.SendPayment(p)
				if err != nil {
					w.channel <- p
					return
				}
			}

			err = w.repository.Save(p)
			if err != nil {
				fmt.Println("Error saving payment:", err)
				p.Err = "SAVE"
				w.channel <- p
			}
		}()
	}
}
