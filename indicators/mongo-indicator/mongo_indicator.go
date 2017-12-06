package mongo_indicator

import (
	"time"

	"log"

	monitor "github.com/andrexus/go-health-monitor"
	"github.com/andrexus/go-health-monitor/status"
	"gopkg.in/mgo.v2"
)

const (
	mongoCheckInterval = time.Second * 3
)

func NewMongoHealthIndicator(session *mgo.Session) monitor.HealthIndicator {
	return &MongoHealthIndicatorImpl{
		Session: session,
		health: &monitor.Health{
			Status:  status.UNKNOWN,
			Details: map[string]interface{}{},
		},
	}
}

type MongoHealthIndicatorImpl struct {
	Session *mgo.Session
	health  *monitor.Health
}

type MongoHealthStatus struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (c *MongoHealthIndicatorImpl) Name() string {
	return "mongo"
}

func (c *MongoHealthIndicatorImpl) Health() *monitor.Health {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	return c.health
}

func (c *MongoHealthIndicatorImpl) StartHealthCheck() {
	go func() {
		c.checkMongoConnection()
		c.health.Details = c.getDetails()
		for range time.Tick(mongoCheckInterval) {
			c.checkMongoConnection()
		}
	}()
}

func (c *MongoHealthIndicatorImpl) checkMongoConnection() {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	if err := c.Session.Ping(); err != nil {
		c.health.Status = status.DOWN
		c.health.Details = nil
		log.Println("Lost connection to mongodb")
		//logrus.Error("Lost connection to mongodb")
		c.Session.Refresh()
		if err := c.Session.Ping(); err == nil {
			c.health.Status = status.UP
			c.health.Details = c.getDetails()
			log.Println("Reconnect to mongodb successful")
			//logrus.Info("Reconnect to mongodb successful")
		}
	} else {
		if c.health.Status != status.UP {
			c.health.Status = status.UP
			c.health.Details = c.getDetails()
		}
	}
}

func (c *MongoHealthIndicatorImpl) getDetails() map[string]interface{} {
	details := map[string]interface{}{}
	buildInfo, err := c.Session.BuildInfo()
	if err == nil {
		details["version"] = buildInfo.Version
	}

	return details
}
