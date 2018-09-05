package agent

import (
	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-agent/metrics"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	log "github.com/sirupsen/logrus"
	"time"
)

type metricsAgent struct {
	nodeHostname    string
	containerDaemon daemon.ContainerDaemon
	pfclient        pfclient.Pfclient
}

func (a *metricsAgent) Run() {
	log.WithFields(log.Fields{}).Warn("Push Metrics")

	for {
		delay := 60 + util.RandomIntRange(1, 10)
		time.Sleep(time.Duration(delay) * time.Second)

		a.Process()
	}
}

func (a *metricsAgent) Process() bool {
	m := metrics.NewMetrics()
	collectedMetrics := m.Collect()
	err := a.pfclient.PushMetrics(collectedMetrics)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}
