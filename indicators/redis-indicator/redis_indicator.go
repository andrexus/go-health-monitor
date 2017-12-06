package redis_indicator

import (
	"time"

	"strings"

	"log"

	monitor "github.com/andrexus/go-health-monitor"
	"github.com/andrexus/go-health-monitor/status"
	"github.com/garyburd/redigo/redis"
)

const (
	redisCheckInterval = time.Second * 3
)

func NewRedisHealthIndicator(pool *redis.Pool) monitor.HealthIndicator {
	return &RedisHealthIndicatorImpl{
		Pool: pool,
		health: &monitor.Health{
			Status:  status.UNKNOWN,
			Details: map[string]interface{}{},
		},
	}
}

type RedisHealthIndicatorImpl struct {
	Pool   *redis.Pool
	health *monitor.Health
}

func (c *RedisHealthIndicatorImpl) Name() string {
	return "redis"
}

func (c *RedisHealthIndicatorImpl) Health() *monitor.Health {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	return c.health
}

func (c *RedisHealthIndicatorImpl) StartHealthCheck() {
	go func() {
		c.checkRedisConnection()
		for range time.Tick(redisCheckInterval) {
			c.checkRedisConnection()
		}
	}()
}

func (c *RedisHealthIndicatorImpl) checkRedisConnection() {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	conn := c.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		if c.health.Status != status.DOWN {
			c.health.Status = status.DOWN
			c.health.Details = nil
			//logrus.Error("Lost connection to redis")
			log.Println("Lost connection to redis")
		}
	} else {
		if c.health.Status != status.UP {
			c.health.Status = status.UP
			c.health.Details = c.getDetails()
			log.Println("Reconnect to redis successful")
			//logrus.Info("Reconnect to redis successful")
		}
	}
}

func (c *RedisHealthIndicatorImpl) getDetails() map[string]interface{} {
	conn := c.Pool.Get()
	defer conn.Close()
	data, err := redis.String(conn.Do("INFO"))
	details := map[string]interface{}{}
	if err == nil {
		info := parseRedisInfo(data)
		if val, ok := info["redis_version"]; ok {
			details["version"] = val
		}
		if val, ok := info["used_memory"]; ok {
			details["used_memory"] = val
		}
		if val, ok := info["used_memory_human"]; ok {
			details["used_memory_human"] = val
		}
	}
	return details
}

func parseRedisInfo(in string) map[string]string {
	info := map[string]string{}
	lines := strings.Split(in, "\r\n")

	for _, line := range lines {
		values := strings.Split(line, ":")

		if len(values) > 1 {
			info[values[0]] = values[1]
		}
	}
	return info
}
