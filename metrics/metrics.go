package metrics

import (
	"github.com/pathfinder-cm/pathfinder-agent/model"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

type Metrics interface {
	Collect() model.Metrics
	collectMemoryMetrics() model.Memory
}

type metrics struct{}

func NewMetrics() Metrics {
	return &metrics{}
}

func (m *metrics) Collect() model.Metrics {
	memory := m.collectMemoryMetrics()

	return model.Metrics{
		Memory: memory,
	}
}

func (m *metrics) collectMemoryMetrics() model.Memory {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err.Error())
	}

	return model.Memory{
		Used:  util.BToMb(vmStat.Used),
		Free:  util.BToMb(vmStat.Free),
		Total: util.BToMb(vmStat.Total),
	}
}
