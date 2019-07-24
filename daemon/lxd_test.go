package daemon

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/pathfinder-cm/pathfinder-agent/config"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
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
		container pfmodel.Container
	}{
		{
			pfmodel.Container{
				Hostname: "test-01",
				Source: pfmodel.Source{
					Type:  "image",
					Mode:  "pull",
					Alias: "16.04",
					Remote: pfmodel.Remote{
						Server:      "https://cloud-images.ubuntu.com/releases",
						Protocol:    "simplestream",
						AuthType:    "tls",
						Certificate: "random",
					},
				},
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	createReq := api.ContainersPost{
		Name: tables[0].container.Hostname,
		Source: api.ContainerSource{
			Type:        tables[0].container.Source.Type,
			Server:      tables[0].container.Source.Remote.Server,
			Protocol:    tables[0].container.Source.Remote.Protocol,
			Alias:       tables[0].container.Source.Alias,
			Mode:        tables[0].container.Source.Mode,
			Certificate: tables[0].container.Source.Remote.Certificate,
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
	mockContainerServer.EXPECT().CreateContainer(createReq).Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		UpdateContainerState(tables[0].container.Hostname, startReq, "").
		Return(mockOperation, nil)
	mockContainerServer.EXPECT().
		GetContainerState(tables[0].container.Hostname).
		Return(&state, "", nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _, _ := l.CreateContainer(tables[0].container)
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

func TestCreateContainerBootstrapFile(t *testing.T) {
	bootstrappers := []pfmodel.Bootstrapper{
		pfmodel.Bootstrapper{
			Type:         "chef-solo",
			CookbooksUrl: "127.0.0.1",
			Attributes:   "{}",
		},
	}

	tables := []struct {
		container pfmodel.Container
	}{
		{
			pfmodel.Container{
				Hostname: "test-01",
				Source: pfmodel.Source{
					Type:  "image",
					Mode:  "pull",
					Alias: "16.04",
					Remote: pfmodel.Remote{
						Server:      "https://cloud-images.ubuntu.com/releases",
						Protocol:    "simplestream",
						AuthType:    "tls",
						Certificate: "random",
					},
				},
				Bootstrappers: bootstrappers,
			},
		},
	}

	content, _, err := util.GenerateBootstrapFileContent(bootstrappers[0])

	exceptedContent := `cd /tmp && curl -LO https://www.chef.io/chef/install.sh && sudo bash ./install.sh -v 14.12.3 && rm install.sh
cat > solo.rb << EOF
root = File.absolute_path(File.dirname(__FILE__))
cookbook_path root + "/cookbooks"
EOF'`
	execChefSoloCmd := fmt.Sprintf("chef-solo -c ~/tmp/solo.rb -j %s %s", bootstrappers[0].Attributes, bootstrappers[0].CookbooksUrl)
	exceptedContent = exceptedContent + "\n" + execChefSoloCmd

	if err != nil {
		t.Errorf("Incorrect content generated, got: %s, want: %s.",
			content,
			exceptedContent)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOperation := mock.NewMockOperation(mockCtrl)
	mockOperation.EXPECT().Wait().Return(nil).AnyTimes()

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)

	mockContainerServer.EXPECT().CreateContainerFile(tables[0].container.Hostname, config.RelativeBootstrapPath, gomock.Any()).
		Return(nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	err = l.CreateContainerBootstrapFile(tables[0].container)
	if err != nil {
		t.Errorf("Container file failed to create")
	}
}

func TestExecContainerBootstrap(t *testing.T) {
	bootstrappers := []pfmodel.Bootstrapper{
		pfmodel.Bootstrapper{
			Type:         "chef-solo",
			CookbooksUrl: "127.0.0.1",
			Attributes:   "{}",
		},
	}

	tables := []struct {
		container pfmodel.Container
	}{
		{
			pfmodel.Container{
				Hostname: "test-01",
				Source: pfmodel.Source{
					Type:  "image",
					Mode:  "pull",
					Alias: "16.04",
					Remote: pfmodel.Remote{
						Server:      "https://cloud-images.ubuntu.com/releases",
						Protocol:    "simplestream",
						AuthType:    "tls",
						Certificate: "random",
					},
				},
				Bootstrappers: bootstrappers,
			},
		},
	}

	commands := []string{
		"sh",
		"-c",
	}
	execBootstrapShCmd := fmt.Sprintf("./%s", config.RelativeBootstrapPath)
	commands = append(commands, execBootstrapShCmd)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	execReq := api.ContainerExecPost{
		Command:     commands,
		WaitForWS:   true,
		Interactive: true,
		Width:       80,
		Height:      15,
	}

	// Setup the exec arguments (fds)
	args := client.ContainerExecArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	mockOperation := mock.NewMockOperation(mockCtrl)
	mockOperation.EXPECT().Wait().Return(nil).AnyTimes()

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	mockContainerServer.EXPECT().ExecContainer(tables[0].container.Hostname, execReq, &args).Return(mockOperation, nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _ := l.ExecContainerBootstrap(tables[0].container)
	if ok != true {
		t.Errorf("Container not properly bootstrapped")
	}
}
