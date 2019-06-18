package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

func TestProvisionProcess(t *testing.T) {
	node := "test-01"

	scs := make(pfmodel.ContainerList, 4)
	scs[0] = pfmodel.Container{Hostname: "test-c-01", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", 
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	scs[1] = pfmodel.Container{Hostname: "test-c-02", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", 
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	scs[2] = pfmodel.Container{Hostname: "test-c-03", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", 
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	scs[3] = pfmodel.Container{Hostname: "test-c-04", Status: "SCHEDULE_DELETION", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", 
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}

	lcs := make(pfmodel.ContainerList, 3)
	lcs[0] = pfmodel.Container{Hostname: "test-c-01", Source: pfmodel.Source{Alias: "16.04"}}
	lcs[1] = pfmodel.Container{Hostname: "test-c-02", Source: pfmodel.Source{Alias: "16.04"}}
	lcs[2] = pfmodel.Container{Hostname: "test-c-04", Source: pfmodel.Source{Alias: "16.04"}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil).AnyTimes()
	mockContainerDaemon.EXPECT().CreateContainer(scs[2]).Return(true, "127.0.0.1", nil)

	mockContainerDaemon.EXPECT().DeleteContainer(scs[3].Hostname).Return(true, nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-03", "127.0.0.1").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-02").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-03").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsDeleted(node, "test-c-04").Return(true, nil)

	provisionAgent := NewProvisionAgent(node, mockContainerDaemon, mockPfClient)
	ok := provisionAgent.Process()
	if ok != true {
		t.Errorf("Agent does not process properly")
	}
}
