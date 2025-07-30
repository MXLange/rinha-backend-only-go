package repository

import (
	"sync"
	"time"

	"github.com/MXLange/rinha-only-go/entities"
)

type MemoryRepository struct {
	mu       sync.Mutex
	data     map[time.Time][]entities.Payment
	getMutex *sync.Mutex
}

func NewMemoryRepository() (*MemoryRepository, *sync.Mutex) {

	mu := &sync.Mutex{}

	return &MemoryRepository{
		mu:       sync.Mutex{},
		data:     make(map[time.Time][]entities.Payment),
		getMutex: mu,
	}, mu
}

func (r *MemoryRepository) Save(payment *entities.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, err := time.Parse(time.RFC3339, payment.RequestedAt)
	if err != nil {
		return err
	}

	if _, exists := r.data[t]; !exists {
		r.data[t] = []entities.Payment{
			*payment,
		}
		return nil
	}

	r.data[t] = append(r.data[t], *payment)
	return nil
}

func (r *MemoryRepository) GetSummary(from, to *time.Time) entities.PaymentSummary {
	r.mu.Lock()
	defer r.mu.Unlock()

	if from == nil && to == nil {
		return r.GetAllPaymentsSummary()
	}
	if from != nil && to != nil {
		return r.GetFromToPaymentsSummary(from, to)
	}
	if from != nil {
		return r.GetFromPaymentsSummary(from)
	}
	return r.GetToPaymentsSummary(to)
}

func (r *MemoryRepository) GetFromPaymentsSummary(from *time.Time) entities.PaymentSummary {
	defaultSummary := entities.Summary{}
	fallbackSummary := entities.Summary{}

	for t, payments := range r.data {
		if t.After(*from) {
			for _, payment := range payments {
				if payment.IsDefault {
					defaultSummary.TotalRequests++
					defaultSummary.TotalAmount += payment.Amount
				} else {
					fallbackSummary.TotalRequests++
					fallbackSummary.TotalAmount += payment.Amount
				}
			}
		}
	}

	return entities.PaymentSummary{
		Default:  defaultSummary,
		Fallback: fallbackSummary,
	}
}

func (r *MemoryRepository) GetToPaymentsSummary(to *time.Time) entities.PaymentSummary {
	defaultSummary := entities.Summary{}
	fallbackSummary := entities.Summary{}

	for t, payments := range r.data {
		if t.Before(*to) {
			for _, payment := range payments {
				if payment.IsDefault {
					defaultSummary.TotalRequests++
					defaultSummary.TotalAmount += payment.Amount
				} else {
					fallbackSummary.TotalRequests++
					fallbackSummary.TotalAmount += payment.Amount
				}
			}
		}
	}
	return entities.PaymentSummary{
		Default:  defaultSummary,
		Fallback: fallbackSummary,
	}
}

func (r *MemoryRepository) GetFromToPaymentsSummary(from, to *time.Time) entities.PaymentSummary {
	defaultSummary := entities.Summary{}
	fallbackSummary := entities.Summary{}

	for t, payments := range r.data {
		if t.After(*from) && t.Before(*to) {
			for _, payment := range payments {
				if payment.IsDefault {
					defaultSummary.TotalRequests++
					defaultSummary.TotalAmount += payment.Amount
				} else {
					fallbackSummary.TotalRequests++
					fallbackSummary.TotalAmount += payment.Amount
				}
			}
		}
	}

	return entities.PaymentSummary{
		Default:  defaultSummary,
		Fallback: fallbackSummary,
	}
}

func (r *MemoryRepository) GetAllPaymentsSummary() entities.PaymentSummary {
	defaultSummary := entities.Summary{}
	fallbackSummary := entities.Summary{}

	for _, payments := range r.data {
		for _, payment := range payments {
			if payment.IsDefault {
				defaultSummary.TotalRequests++
				defaultSummary.TotalAmount += payment.Amount
			} else {
				fallbackSummary.TotalRequests++
				fallbackSummary.TotalAmount += payment.Amount
			}
		}
	}

	return entities.PaymentSummary{
		Default:  defaultSummary,
		Fallback: fallbackSummary,
	}
}
