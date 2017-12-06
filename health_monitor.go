package go_health_monitor

import (
	"sort"

	"net/http"

	"github.com/andrexus/go-health-monitor/status"
)

type Health struct {
	Status  status.HealthStatus    `json:"status"`
	Details map[string]interface{} `json:"details"`
}

type HealthIndicator interface {
	Name() string
	Health() *Health
	StartHealthCheck()
}

type HealthMonitorService interface {
	HealthStatus() status.HealthStatus
	StatusCode() int
	HealthReport(showDetails bool) map[string]interface{}
	RegisterHealthIndicator(indicator HealthIndicator)
}

type HealthMonitorServiceImpl struct {
	healthStatus     status.HealthStatus
	healthAggregator HealthAggregator
	indicators       []HealthIndicator
}

func NewHealthMonitorService() HealthMonitorService {
	s := &HealthMonitorServiceImpl{healthStatus: status.UNKNOWN}
	s.healthAggregator = &OrderedHealthAggregator{
		statusOrder: []status.HealthStatus{status.DOWN, status.OUT_OF_SERVICE, status.UP, status.UNKNOWN},
	}
	return s
}

func (c *HealthMonitorServiceImpl) HealthStatus() status.HealthStatus {
	return c.healthAggregator.aggregateStatus(c.getAggregatorCandidates())
}

func (c *HealthMonitorServiceImpl) StatusCode() int {
	if c.HealthStatus() == status.UP || c.HealthStatus() == status.UNKNOWN {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

func (c *HealthMonitorServiceImpl) HealthReport(showDetails bool) map[string]interface{} {
	c.healthStatus = c.healthAggregator.aggregateStatus(c.getAggregatorCandidates())
	result := map[string]interface{}{"status": c.healthStatus}
	if showDetails {
		for _, indicator := range c.indicators {
			name := indicator.Name()
			health := indicator.Health()
			indicatorStatus := map[string]interface{}{
				"status": health.Status,
			}
			for k, v := range health.Details {
				indicatorStatus[k] = v
			}
			result[name] = indicatorStatus
		}
	}
	return result
}

func (c *HealthMonitorServiceImpl) RegisterHealthIndicator(indicator HealthIndicator) {
	c.indicators = append(c.indicators, indicator)
	indicator.StartHealthCheck()
}

func (c *HealthMonitorServiceImpl) getAggregatorCandidates() []status.HealthStatus {
	var candidates []status.HealthStatus
	for _, indicator := range c.indicators {
		candidates = append(candidates, indicator.Health().Status)
	}
	return candidates
}

type HealthAggregator interface {
	aggregateStatus([]status.HealthStatus) status.HealthStatus
}

type OrderedHealthAggregator struct {
	statusOrder []status.HealthStatus
}

func (c *OrderedHealthAggregator) aggregateStatus(candidates []status.HealthStatus) status.HealthStatus {
	if len(candidates) == 0 {
		return status.UNKNOWN
	}
	sort.Slice(candidates, func(i, j int) bool {
		i1 := index(c.statusOrder, candidates[i])
		i2 := index(c.statusOrder, candidates[j])
		if i1 < i2 {
			return true
		} else if i1 == i2 {
			return candidates[i] == candidates[j]
		}
		return false
	})
	//fmt.Println("Sorted status: ", candidates)
	return candidates[0]
}

func index(vs []status.HealthStatus, t status.HealthStatus) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
