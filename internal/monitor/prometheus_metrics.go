package monitor

import "github.com/prometheus/client_golang/prometheus"

type PrometheusMetrics struct {
	healthyGauge   PrometheusGaugeVec
	lastCheckGauge PrometheusGaugeVec
}

type prometheusGaugeVecWrapper struct {
	gaugeVec    *prometheus.GaugeVec
	labels      []ResourceLabel
	labelValues map[ResourceLabel]string
}

var _ PrometheusGaugeVecSetup = (*prometheusGaugeVecWrapper)(nil)
var _ PrometheusGaugeVec = (*prometheusGaugeVecWrapper)(nil)

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

func newPrometheusGaugeVecWrapper(name PrometheusMetricName, help string, labels ...ResourceLabel) PrometheusGaugeVecSetup {
	opts := prometheus.GaugeOpts{Name: string(name), Help: help}
	labelsAsStrings := make([]string, 0)
	for _, label := range labels {
		labelsAsStrings = append(labelsAsStrings, string(label))
	}
	return &prometheusGaugeVecWrapper{
		gaugeVec: prometheus.NewGaugeVec(opts, labelsAsStrings),
		labels:   labels,
	}
}

func (p *prometheusGaugeVecWrapper) Set(value float64) {
	values := make([]string, 0)
	for _, label := range p.labels {
		resourceLabel := label
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
	AvailabilityState ResourceLabel = "availability_state"
	ResourceGroup     ResourceLabel = "resource_group"
	ResourceID        ResourceLabel = "resource_id"
	ResourceName      ResourceLabel = "resource_name"
	ResourceType      ResourceLabel = "resource_type"
)

func withResourceGroup(resourceGroup string) func() (ResourceLabel, string) {
	return with(ResourceGroup, resourceGroup)
}

func withResourceID(resourceID string) func() (ResourceLabel, string) {
	return with(ResourceID, resourceID)
}

func withResourceName(resourceName string) func() (ResourceLabel, string) {
	return with(ResourceName, resourceName)
}

func withResourceType(resourceType string) func() (ResourceLabel, string) {
	return with(ResourceType, resourceType)
}

func withAvailabilityState(availabilityState string) func() (ResourceLabel, string) {
	return with(AvailabilityState, availabilityState)
}

func with(resourceLabel ResourceLabel, value string) func() (ResourceLabel, string) {
	return func() (ResourceLabel, string) {
		return resourceLabel, value
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
	healthyGauge := newPrometheusGaugeVecWrapper(AzureResourceHealthStatus, "Azure resource health status (1 = healthy, 0 = unhealthy)", ResourceID, ResourceGroup, ResourceName, ResourceType, AvailabilityState).Register()
	lastCheckGauge := newPrometheusGaugeVecWrapper(AzureResourceHealthLastCheckTimestamp, "Timestamp of last health check", ResourceID, ResourceGroup, ResourceName).Register()
	return &PrometheusMetrics{
		healthyGauge:   healthyGauge,
		lastCheckGauge: lastCheckGauge,
	}, nil
}
