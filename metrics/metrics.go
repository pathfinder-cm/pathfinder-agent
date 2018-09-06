package metrics

import (
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

func Collect() *pfmodel.Metrics {
	memory := collectMemory()
	return &pfmodel.Metrics{
		Memory: memory,
	}
}

func collectMemory() *pfmodel.Memory {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &pfmodel.Memory{
		Used:  util.BToMb(vmStat.Used),
		Free:  util.BToMb(vmStat.Free),
		Total: util.BToMb(vmStat.Total),
	}
}
