package api

type RegisterRequest struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

type RegisterResponse struct {
	ID                string `json:"id"`
	HeartbeatInterval int    `json:"heart_beat"`
	QueueName         string `json:"queue_name"`
}
