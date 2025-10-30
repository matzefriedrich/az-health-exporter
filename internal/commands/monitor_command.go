package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
}

var _ types.TypedCommand = (*healthMonitorCommand)(nil)

func NewHealthMonitorCommand(
	healthMonitor monitor.HealthMonitor) *cobra.Command {
	instance := &healthMonitorCommand{
		healthMonitor: healthMonitor,
	}
	return commands.CreateTypedCommand(instance)
}

func (h *healthMonitorCommand) Execute(ctx context.Context) {

	go h.healthMonitor.StartMonitoring(ctx)

	http.HandleFunc("/health", healthCheckHandler)
	// http.HandleFunc("/status", h.healthMonitor.statusHandler)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", h.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  - Health Check: http://localhost%s/health", addr)
	log.Printf("  - Status API:   http://localhost%s/status", addr)
	log.Printf("  - Metrics:      http://localhost%s/metrics", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

// healthCheckHandler returns basic health check
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}
