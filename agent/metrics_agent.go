package agent

import (
	"time"

	"github.com/pathfinder-cm/pathfinder-agent/metrics"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	log "github.com/sirupsen/logrus"
)

type metricsAgent struct {
	nodeHostname string
	pfclient     pfclient.Pfclient
}

func NewMetricsAgent(
	nodeHostname string,
	pfclient pfclient.Pfclient) Agent {

	return &metricsAgent{
		nodeHostname: nodeHostname,
		pfclient:     pfclient,
	}
}

func (a *metricsAgent) Run() {
	log.WithFields(log.Fields{}).Warn("Starting metrics agent...")

	for {
		delay := 60 + util.RandomIntRange(1, 10)
		time.Sleep(time.Duration(delay) * time.Second)

		a.Process()
	}
}

func (a *metricsAgent) Process() bool {
	collectedMetrics := metrics.Collect()
	ok, err := a.pfclient.StoreMetrics(collectedMetrics)
	if !ok {
		log.Error(err.Error())
		return false
	}

	return true
}
