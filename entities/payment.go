package entities

type Payment struct {
	ID          string  `json:"correlationId"`
	Amount      float64 `json:"amount"`
	RequestedAt string  `json:"requestedAt"` //2025-07-15T12:34:56.000Z
	IsDefault   bool    `json:"-"`
	Err         string  `json:"-"`
	Attempts    int     `json:"-"`
}

type PaymentSummary struct {
	Default  Summary `json:"default"`
	Fallback Summary `json:"fallback"`
}

type Summary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}
