package agent

import (
	"testing"

	"github.com/giosakti/pathfinder-agent/mock"
	"github.com/giosakti/pathfinder-agent/model"
	"github.com/golang/mock/gomock"
)

func TestProcessContainerExist(t *testing.T) {
	node := "test-01"

	scs := make(model.ContainerList, 3)
	scs[0] = model.Container{Hostname: "test-c-01", Image: "16.04"}
	scs[1] = model.Container{Hostname: "test-c-02", Image: "16.04"}
	scs[2] = model.Container{Hostname: "test-c-03", Image: "16.04"}

	lcs := make(model.ContainerList, 2)
	lcs[0] = model.Container{Hostname: "test-c-01", Image: "16.04"}
	lcs[1] = model.Container{Hostname: "test-c-02", Image: "16.04"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil)
	mockContainerDaemon.EXPECT().CreateContainer("test-c-03", "16.04").Return(true, nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-03").Return(true, nil)

	a := NewAgent(node, mockContainerDaemon, mockPfClient)
	ok := a.Process()
	if ok != true {
		t.Errorf("Agent does not process properly")
	}
}
