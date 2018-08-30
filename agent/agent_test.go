package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-agent/model"
)

func TestProcess(t *testing.T) {
	node := "test-01"

	scs := make(model.ContainerList, 4)
	scs[0] = model.Container{Hostname: "test-c-01", Image: "16.04", Status: "SCHEDULED"}
	scs[1] = model.Container{Hostname: "test-c-02", Image: "16.04", Status: "SCHEDULED"}
	scs[2] = model.Container{Hostname: "test-c-03", Image: "16.04", Status: "SCHEDULED"}
	scs[3] = model.Container{Hostname: "test-c-04", Image: "16.04", Status: "SCHEDULE_DELETION"}

	lcs := make(model.ContainerList, 3)
	lcs[0] = model.Container{Hostname: "test-c-01", Image: "16.04"}
	lcs[1] = model.Container{Hostname: "test-c-02", Image: "16.04"}
	lcs[2] = model.Container{Hostname: "test-c-04", Image: "16.04"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil)
	mockContainerDaemon.EXPECT().CreateContainer("test-c-03", "16.04").Return(true, "127.0.0.1", nil)
	mockContainerDaemon.EXPECT().DeleteContainer("test-c-04").Return(true, nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-03", "127.0.0.1").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-02").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-03").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsDeleted(node, "test-c-04").Return(true, nil)

	a := NewAgent(node, mockContainerDaemon, mockPfClient)
	ok := a.Process()
	if ok != true {
		t.Errorf("Agent does not process properly")
	}
}
