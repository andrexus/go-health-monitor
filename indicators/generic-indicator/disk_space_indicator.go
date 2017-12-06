package generic_indicator

import (
	"os"
	"syscall"

	"time"

	monitor "github.com/andrexus/go-health-monitor"
	"github.com/andrexus/go-health-monitor/status"
	"github.com/dustin/go-humanize"
)

const (
	megabytes         = 1024 * 1024
	defaultThreshold  = 10 * megabytes
	diskCheckInterval = time.Second * 30
)

type DiskSpaceIndicator struct {
	health *monitor.Health
}

func (c *DiskSpaceIndicator) Name() string {
	return "diskSpace"
}

func (c *DiskSpaceIndicator) Health() *monitor.Health {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	return c.health
}

func (c *DiskSpaceIndicator) StartHealthCheck() {
	go func() {
		c.checkDiskUsage()
		for range time.Tick(diskCheckInterval) {
			c.checkDiskUsage()
		}
	}()
}

func (c *DiskSpaceIndicator) checkDiskUsage() {
	if c.health == nil {
		c.health = &monitor.Health{Status: status.UNKNOWN}
	}
	if wd, err := os.Getwd(); err == nil {
		diskUsage := getDiskUsage(wd)
		if diskUsage.Free > defaultThreshold {
			c.health.Status = status.UP
		} else {
			c.health.Status = status.DOWN
		}
		c.health.Details = map[string]interface{}{
			"total":           diskUsage.All,
			"total_human":     humanize.Bytes(diskUsage.All),
			"free":            diskUsage.Free,
			"free_human":      humanize.Bytes(diskUsage.Free),
			"used":            diskUsage.Used,
			"used_human":      humanize.Bytes(diskUsage.Used),
			"threshold":       defaultThreshold,
			"threshold_human": humanize.Bytes(defaultThreshold),
		}
	} else {
		c.health.Status = status.UNKNOWN
		c.health.Details = nil
	}
}

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func getDiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
