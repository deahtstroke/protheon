package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"sync"
)

type Worker struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	IP       string `json:"ip"`
}

var (
	workers   = make(map[string]Worker)
	workersMu sync.Mutex
)

func RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	id := uuid.New().String()
	worker := Worker{
		ID:       id,
		Hostname: req.Hostname,
		OS:       req.OS,
		IP:       host,
	}

	workersMu.Lock()
	workers[id] = worker
	workersMu.Unlock()

	resp := RegisterResponse{
		ID:                id,
		HeartbeatInterval: 30,
		QueueName:         "jobs",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Printf("[Conductor] Worker registered: %+v\n", worker)
}

func ReceiveHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	workersMu.Lock()
	worker, ok := workers[req.ID]
	workersMu.Unlock()

	if !ok {
		log.Printf("Unable to find worker with Id [%s]", req.ID)
	} else {
		log.Printf("Received heartbet for worker: %+v", worker)
	}
}
