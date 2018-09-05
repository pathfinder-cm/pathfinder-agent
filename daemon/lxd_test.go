package daemon

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lxc/lxd/shared/api"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
)

func TestToContainerList(t *testing.T) {
	tables := []api.Container{
		api.Container{Name: "test-01"},
		api.Container{Name: "test-02"},
		api.Container{Name: "test-03"},
	}

	cl := apiContainers(tables).toContainerList()
	for i, table := range tables {
		if (*cl)[i].Hostname != table.Name {
			t.Errorf("Incorrect container hostname generated, got: %s, want: %s.",
				(*cl)[i].Hostname,
				table.Name)
		}
	}
}

func TestListContainers(t *testing.T) {
	tables := []api.Container{
		api.Container{Name: "test-01"},
		api.Container{Name: "test-02"},
		api.Container{Name: "test-03"},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	mockContainerServer.EXPECT().GetContainers().Return(tables, nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	result, _ := l.ListContainers()
	for i, table := range tables {
		if (*result)[i].Hostname != table.Name {
			t.Errorf("Incorrect container hostname generated, got: %s, want: %s.",
				(*result)[i].Hostname,
				table.Name)
		}
	}
}

func TestCreateContainer(t *testing.T) {
	tables := []struct {
		hostname string
		image    string
	}{
		{"test-01", "16.04"},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	createReq := api.ContainersPost{
		Name: tables[0].hostname,
		Source: api.ContainerSource{
			Type:     "image",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Protocol: "simplestreams",
			Alias:    tables[0].image,
		},
	}

	startReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	mockOperation := mock.NewMockOperation(mockCtrl)
	mockOperation.EXPECT().Wait().Return(nil).AnyTimes()

	state := api.ContainerState{
		Network: map[string]api.ContainerStateNetwork{
			"eth0": api.ContainerStateNetwork{
				Addresses: []api.ContainerStateNetworkAddress{
					api.ContainerStateNetworkAddress{
						Family:  "inet",
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	mockContainerServer.EXPECT().IsClustered().Return(false)
	mockContainerServer.EXPECT().CreateContainer(createReq).Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		UpdateContainerState(tables[0].hostname, startReq, "").
		Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		GetContainerState(tables[0].hostname).
		Return(&state, "", nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _, _ := l.CreateContainer(tables[0].hostname, tables[0].image)
	if ok != true {
		t.Errorf("Container not properly generated")
	}
}

func TestDeleteContainer(t *testing.T) {
	tables := []struct {
		hostname string
	}{
		{"test-01"},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stopReq := api.ContainerStatePut{
		Action:  "stop",
		Timeout: 60,
	}

	mockOperation := mock.NewMockOperation(mockCtrl)
	mockOperation.EXPECT().Wait().Return(nil).AnyTimes()

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	mockContainerServer.EXPECT().DeleteContainer(tables[0].hostname).Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		UpdateContainerState(tables[0].hostname, stopReq, "").
		Return(mockOperation, nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _ := l.DeleteContainer(tables[0].hostname)
	if ok != true {
		t.Errorf("Container not properly deleted")
	}
}
