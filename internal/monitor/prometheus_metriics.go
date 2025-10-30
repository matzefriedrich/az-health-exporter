package monitor

import "github.com/prometheus/client_golang/prometheus"

type PrometheusMetrics struct {
	healthyGauge   PrometheusGaugeVec
	lastCheckGauge PrometheusGaugeVec
}

type prometheusGaugeVecWrapper struct {
	PrometheusGaugeVecSetup
	PrometheusGaugeVec

	gaugeVec    *prometheus.GaugeVec
	labels      []string
	labelValues map[ResourceLabel]string
}

func (p *prometheusGaugeVecWrapper) Register() PrometheusGaugeVec {
	prometheus.MustRegister(p.gaugeVec)
	return p
}

type PrometheusGaugeVecSetup interface {
	Register() PrometheusGaugeVec
}

type PrometheusGaugeVec interface {
	WithLabelValues(labels ...func() (ResourceLabel, string)) PrometheusGaugeVec
	Set(value float64)
}

func newPrometheusGaugeVecWrapper(name PrometheusMetricName, help string, labels ...string) PrometheusGaugeVecSetup {
	opts := prometheus.GaugeOpts{Name: string(name), Help: help}
	return &prometheusGaugeVecWrapper{
		gaugeVec: prometheus.NewGaugeVec(opts, labels),
		labels:   labels,
	}
}

func (p *prometheusGaugeVecWrapper) Set(value float64) {
	values := make([]string, len(p.labels))
	for _, label := range p.labels {
		resourceLabel := ResourceLabel(label)
		v, found := p.labelValues[resourceLabel]
		if !found {
			v = ""
		}
		values = append(values, v)
	}
	p.gaugeVec.WithLabelValues(values...).Set(value)
}

type ResourceLabel string

const (
	ResourceID        ResourceLabel = "resource_id"
	ResourceName      ResourceLabel = "resource_name"
	ResourceType      ResourceLabel = "resource_type"
	AvailabilityState ResourceLabel = "availability_state"
)

func withResourceID(resourceID string) func() (ResourceLabel, string) {
	return func() (ResourceLabel, string) {
		return ResourceID, resourceID
	}
}

func withResourceName(resourceName string) func() (ResourceLabel, string) {
	return func() (ResourceLabel, string) {
		return ResourceName, resourceName
	}
}

func withResourceType(resourceType string) func() (ResourceLabel, string) {
	return func() (ResourceLabel, string) {
		return ResourceType, resourceType
	}
}

func withAvailabilityState(availabilityState string) func() (ResourceLabel, string) {
	return func() (ResourceLabel, string) {
		return AvailabilityState, availabilityState
	}
}

func (p *prometheusGaugeVecWrapper) WithLabelValues(labels ...func() (ResourceLabel, string)) PrometheusGaugeVec {
	values := make(map[ResourceLabel]string)
	for _, labelFunc := range labels {
		name, value := labelFunc()
		values[name] = value
	}
	p.labelValues = values
	return p
}

type PrometheusMetricName string

const (
	AzureResourceHealthStatus             PrometheusMetricName = "azure_resource_health_status"
	AzureResourceHealthLastCheckTimestamp PrometheusMetricName = "azure_resource_health_last_check_timestamp"
)

func registerPrometheusMetrics() (*PrometheusMetrics, error) {
	healthyGauge := newPrometheusGaugeVecWrapper(AzureResourceHealthStatus, "Azure resource health status (1 = healthy, 0 = unhealthy)", "resource_id", "resource_name", "resource_type", "availability_state").Register()
	lastCheckGauge := newPrometheusGaugeVecWrapper(AzureResourceHealthLastCheckTimestamp, "Timestamp of last health check", "resource_id", "resource_name").Register()
	return &PrometheusMetrics{
		healthyGauge:   healthyGauge,
		lastCheckGauge: lastCheckGauge,
	}, nil
}
