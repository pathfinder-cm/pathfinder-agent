package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
)

func TestMetricsProcess(t *testing.T) {
	node := "test-01"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPfClient := mock.NewMockPfclient(mockCtrl)

	metricsAgent := NewMetricsAgent(node, mockPfClient)
	mockPfClient.EXPECT().StoreMetrics(gomock.Any()).Return(true, nil)
	ok := metricsAgent.Process()
	if !ok {
		t.Errorf("Metrics agent does not process properly")
	}
}
