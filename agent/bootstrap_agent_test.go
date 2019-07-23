package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

func TestBootstrapProcess(t *testing.T) {
	node := "test-01"

	bootstrappers := []pfmodel.Bootstrapper{
		pfmodel.Bootstrapper{
			Type:         "chef-solo",
			CookbooksUrl: "127.0.0.1",
			Attributes:   "{}",
		},
	}
	pcs := make(pfmodel.ContainerList, 4)
	pcs[0] = pfmodel.Container{
		Hostname: "test-c-01", Ipaddress: "127.0.0.1", Status: "PROVISIONED", Bootstrappers: bootstrappers,
		Source: pfmodel.Source{
			Type: "image", Alias: "16.04", Mode: "pull",
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	pcs[1] = pfmodel.Container{
		Hostname: "test-c-02", Ipaddress: "127.0.0.2", Status: "PROVISIONED", Bootstrappers: bootstrappers,
		Source: pfmodel.Source{
			Type: "image", Alias: "16.04", Mode: "pull",
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	pcs[2] = pfmodel.Container{
		Hostname: "test-c-03", Ipaddress: "127.0.0.3", Status: "PROVISIONED", Bootstrappers: bootstrappers,
		Source: pfmodel.Source{
			Type: "image", Alias: "16.04", Mode: "pull",
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}
	pcs[3] = pfmodel.Container{
		Hostname: "test-c-04", Ipaddress: "127.0.0.4", Status: "PROVISIONED", Bootstrappers: bootstrappers,
		Source: pfmodel.Source{
			Type: "image", Alias: "16.04", Mode: "pull",
			Remote: pfmodel.Remote{
				Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)

	for _, pc := range pcs {
		mockContainerDaemon.EXPECT().CreateContainerBootstrapFile(pc).Return(nil)
		mockContainerDaemon.EXPECT().ExecContainerBootstrap(pc).Return(true, nil)
	}

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchProvisionedContainersFromServer(node).Return(&pcs, nil)
	mockPfClient.EXPECT().MarkContainerAsBootstrapped(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsBootstrapped(node, "test-c-02").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsBootstrapped(node, "test-c-03").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsBootstrapped(node, "test-c-04").Return(true, nil)

	bootstrapAgent := NewBootstrapAgent(node, mockContainerDaemon, mockPfClient)
	ok := bootstrapAgent.Process()
	if ok != true {
		t.Errorf("Agent does not process properly")
	}
}
