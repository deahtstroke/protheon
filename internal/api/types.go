package api

import (
	"time"
)

type RegisterRequest struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

type RegisterResponse struct {
	ID                string `json:"id"`
	HeartbeatInterval int    `json:"heart_beat"`
	QueueName         string `json:"queue_name"`
}

type HeartbeatRequest struct {
	ID          string    `json:"id"`
	JobsDone    int       `json:"jobs_done"`
	LastJobTime time.Time `json:"last_job_time"`
	Uptime      string    `json:"uptime"`
}
