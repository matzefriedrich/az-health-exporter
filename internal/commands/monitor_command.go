package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/matzefriedrich/az-health-exporter/internal/monitor"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

type healthMonitorCommand struct {
	use  types.CommandName `flag:"monitor" short:"Runs the resource monitor server"`
	Port int               `flag:"p" short:"The listening port of the monitoring server"`

	healthMonitor monitor.HealthMonitor
	mu            sync.RWMutex
}

var _ types.TypedCommand = (*healthMonitorCommand)(nil)

func NewHealthMonitorCommand(
	healthMonitor monitor.HealthMonitor) *cobra.Command {
	instance := &healthMonitorCommand{
		healthMonitor: healthMonitor,
		mu:            sync.RWMutex{},
	}
	return commands.CreateTypedCommand(instance)
}

func (h *healthMonitorCommand) Execute(ctx context.Context) {

	go h.healthMonitor.StartMonitoring(ctx)

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/status", h.statusHandler)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", h.Port)
	log.Printf("Starting server on %s (listening on all interfaces)", addr)
	log.Printf("Endpoints:")
	log.Printf("  - Health Check: %s/health", addr)
	log.Printf("  - Status API:   %s/status", addr)
	log.Printf("  - Metrics:      %s/metrics", addr)
	log.Printf("Access locally at: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

func (m *healthMonitorCommand) statusHandler(w http.ResponseWriter, request *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]*monitor.ResourceHealth, 0)
	healthStatus, err := m.healthMonitor.GetHealthStatus(request.Context())
	if err != nil {
		log.Printf("Failed to get health status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, status := range healthStatus {
		statuses = append(statuses, status)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": time.Now(),
		"resources": statuses,
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
