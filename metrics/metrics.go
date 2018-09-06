package metrics

import (
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

type Metrics interface {
	Collect() pfmodel.Metrics
	collectMemoryMetrics() pfmodel.Memory
}

type metrics struct{}

func NewMetrics() Metrics {
	return &metrics{}
}

func (m *metrics) Collect() pfmodel.Metrics {
	memory := m.collectMemoryMetrics()

	return pfmodel.Metrics{
		Memory: memory,
	}
}

func (m *metrics) collectMemoryMetrics() pfmodel.Memory {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err.Error())
	}

	return pfmodel.Memory{
		Used:  util.BToMb(vmStat.Used),
		Free:  util.BToMb(vmStat.Free),
		Total: util.BToMb(vmStat.Total),
	}
}
