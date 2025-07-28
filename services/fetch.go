package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MXLange/rinha-only-go/entities"
)

type Fetch struct {
	principal     string
	principalHost string
	fallback      string
	fallbackHost  string
	client        *http.Client
}

func NewFetch(baseDefault, baseFallback string) (*Fetch, error) {

	if baseDefault == "" || baseFallback == "" {
		return nil, fmt.Errorf("base URLs cannot be empty")
	}

	return &Fetch{
		principal: baseDefault,
		fallback:  baseFallback,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (f *Fetch) SendPayment(payment *entities.Payment) (error, bool) {

	if payment == nil {
		return fmt.Errorf("payment cannot be nil"), false
	}

	var err error
	if err = f.send(f.principal, payment); err != nil {
		if err = f.send(f.fallback, payment); err != nil {
			return fmt.Errorf("failed to send payment to both services: %v", err), false
		}
		return nil, false
	}

	return nil, true

}

func (f *Fetch) send(url string, payment *entities.Payment) error {

	body, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("error marshalling payment: %v", err)
	}

	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	url = fmt.Sprintf("%s/payments", url)

	fmt.Println("Sending payment to:", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payment to principal service: %v", err)
	}
	defer res.Body.Close()

	fmt.Println("Response status code:", res.StatusCode)

	if res.StatusCode == http.StatusUnprocessableEntity {
		return nil
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send payment to principal service, status code: %d", res.StatusCode)
	}

	return nil
}

func (f *Fetch) GetInstanceSummary(instance string, from, to string) (entities.PaymentSummary, error) {

	if instance == "" {
		return entities.PaymentSummary{}, fmt.Errorf("instance cannot be empty")
	}

	url := fmt.Sprintf("%s/payments-summary?internal=true", instance)
	fmt.Println("Fetching summary from:", url)
	if from != "" {
		url += fmt.Sprintf("&from=%s", from)
	}

	if to != "" {
		url += fmt.Sprintf("&to=%s", to)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return entities.PaymentSummary{}, fmt.Errorf("error creating request: %v", err)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return entities.PaymentSummary{}, fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return entities.PaymentSummary{}, fmt.Errorf("failed to get summary from %s, status code: %d", instance, res.StatusCode)
	}

	var summary entities.PaymentSummary
	if err := json.NewDecoder(res.Body).Decode(&summary); err != nil {
		return entities.PaymentSummary{}, fmt.Errorf("error decoding response: %v", err)
	}

	return summary, nil
}
