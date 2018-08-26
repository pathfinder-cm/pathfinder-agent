package daemon

import (
	"testing"

	"github.com/giosakti/pathfinder-agent/mock"
	"github.com/golang/mock/gomock"
	"github.com/lxc/lxd/shared/api"
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

	l := LXD{conn: mockContainerServer}
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

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	mockContainerServer.EXPECT().CreateContainer(createReq).Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		UpdateContainerState(tables[0].hostname, startReq, "").
		Return(mockOperation, nil)

	l := LXD{conn: mockContainerServer}
	ok, _ := l.CreateContainer(tables[0].hostname, tables[0].image)
	if ok != true {
		t.Errorf("Container not properly generated")
	}
}
