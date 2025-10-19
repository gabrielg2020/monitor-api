package entities

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp int64             `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}
