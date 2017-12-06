package go_health_monitor

import (
	"testing"

	"github.com/andrexus/go-health-monitor/status"
)

type TestHealthIndicator struct {
	health *Health
}

func (c *TestHealthIndicator) Name() string {
	return "test"
}

func (c *TestHealthIndicator) Health() *Health {
	return c.health
}

func (c *TestHealthIndicator) StartHealthCheck() {
}

func TestConfigWithOverrides(t *testing.T) {
	healthMonitorService := NewHealthMonitorService()

	indicator1 := &TestHealthIndicator{health: &Health{Status: status.UP}}
	indicator2 := &TestHealthIndicator{health: &Health{Status: status.DOWN}}
	indicator3 := &TestHealthIndicator{health: &Health{Status: status.DOWN}}
	indicator4 := &TestHealthIndicator{health: &Health{Status: status.OUT_OF_SERVICE}}
	indicator5 := &TestHealthIndicator{health: &Health{Status: status.UP}}
	indicator6 := &TestHealthIndicator{health: &Health{Status: status.UNKNOWN}}
	indicator7 := &TestHealthIndicator{health: &Health{Status: status.OUT_OF_SERVICE}}

	healthMonitorService.RegisterHealthIndicator(indicator1)
	healthMonitorService.RegisterHealthIndicator(indicator2)
	healthMonitorService.RegisterHealthIndicator(indicator3)
	healthMonitorService.RegisterHealthIndicator(indicator4)
	healthMonitorService.RegisterHealthIndicator(indicator5)
	healthMonitorService.RegisterHealthIndicator(indicator6)
	healthMonitorService.RegisterHealthIndicator(indicator7)

	if healthMonitorService.HealthStatus() != status.DOWN {
		t.Error("Expected status.DOWN, got ", healthMonitorService.HealthStatus())
	}
}
