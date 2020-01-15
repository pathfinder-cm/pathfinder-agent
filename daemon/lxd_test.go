package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	if !ok {
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
	if !ok {
		t.Errorf("Container not properly deleted")
	}
}

func TestCreateContainerBootstrapScript(t *testing.T) {
	bytes := []byte(`{
		"consul":{
			"hosts":["guro-consul-01"],
			"config":{
			"consul.json":{"bind_addr":null}}
		},
		"run_list":["role[consul]"]
	}`)
	var attributes interface{}
	json.Unmarshal(bytes, &attributes)

	bootstrappers := []pfmodel.Bootstrapper{
		pfmodel.Bootstrapper{
			Type:         "chef-solo",
			CookbooksUrl: "127.0.0.1",
			Attributes:   attributes,
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

	content, _, _ := util.GenerateBootstrapScriptContent(bootstrappers[0])

	exceptedContent := `
CHEF_FLAG_FILE=/tmp/chef_installed.txt
if [ ! -f "$CHEF_FLAG_FILE" ]; then
	echo "$CHEF_FLAG_FILE doesn't exist, creating file..."
	cd /tmp && curl -LO  && sudo bash ./install.sh -v  && rm install.sh && touch chef_installed.txt
fi

CHEF_REPO_DIR=/tmp/chef-repo-master
[ -d "$CHEF_REPO_DIR" ] && rm -rf $CHEF_REPO_DIR
mkdir $CHEF_REPO_DIR && wget 127.0.0.1 -O - | tar -xz -C /tmp/chef-repo-master --strip-components=1

SOLO_FILE=/tmp/solo.rb
if [ ! -f "$SOLO_FILE" ]; then
	echo "$SOLO_FILE doesn't exist, creating file..."
	cat > solo.rb << EOF
cookbook_path "/tmp/chef-repo-master/cookbooks"
role_path "/tmp/chef-repo-master/roles"
EOF
fi

cat > /tmp/attributes.json << EOF
{"consul":{"config":{"consul.json":{"bind_addr":null}},"hosts":["guro-consul-01"]},"run_list":["role[consul]"]}
EOF

chef-solo -c /tmp/solo.rb -j /tmp/attributes.json 
`
	compareResult := strings.Compare(content,
		exceptedContent)
	if compareResult != 0 {
		t.Errorf("Incorrect content generated, got: %v, want: %v.",
			compareResult,
			0)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOperation := mock.NewMockOperation(mockCtrl)
	mockOperation.EXPECT().Wait().Return(nil).AnyTimes()

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)

	mockContainerServer.EXPECT().CreateContainerFile(tables[0].container.Hostname, config.AbsoluteBootstrapScriptPath, gomock.Any()).
		Return(nil)

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _ := l.CreateContainerBootstrapScript(tables[0].container)
	if !ok {
		t.Errorf("Container file failed to create")
	}
}

func TestBootstrapContainer(t *testing.T) {
	bytes := []byte(`{
		"consul":{
			"hosts":["guro-consul-01"],
			"config":{
			"consul.json":{"bind_addr":null}}
		},
		"run_list":["role[consul]"]
	}`)
	var attributes interface{}
	json.Unmarshal(bytes, &attributes)

	bootstrappers := []pfmodel.Bootstrapper{
		pfmodel.Bootstrapper{
			Type:         "chef-solo",
			CookbooksUrl: "127.0.0.1",
			Attributes:   attributes,
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
		"bash",
	}
	execBootstrapCmd := fmt.Sprintf("%s", config.AbsoluteBootstrapScriptPath)
	commands = append(commands, execBootstrapCmd)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	execReq := api.ContainerExecPost{
		Command:     commands,
		WaitForWS:   true,
		Interactive: false,
	}

	// Setup the exec arguments (fds)
	args := client.ContainerExecArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	opAPIReturnValues := []string{
		`{"return":1}`,
		`{"return":1}`,
		`{"return":0}`,
	}

	mockContainerServer := mock.NewMockContainerServer(mockCtrl)
	var previousExecContainerMock *gomock.Call

	for _, opAPIReturnValue := range opAPIReturnValues {
		opAPI := api.Operation{}
		json.Unmarshal([]byte(opAPIReturnValue), &opAPI.Metadata)

		mockOperation := mock.NewMockOperation(mockCtrl)
		mockOperation.EXPECT().Wait().Return(nil).AnyTimes()
		mockOperation.EXPECT().Get().Return(opAPI).AnyTimes()

		execContainerMock := mockContainerServer.EXPECT().ExecContainer(tables[0].container.Hostname, execReq, &args).Return(mockOperation, nil)
		if previousExecContainerMock != nil {
			execContainerMock.After(previousExecContainerMock)
		}
		previousExecContainerMock = execContainerMock
	}

	l := LXD{localSrv: mockContainerServer, targetSrv: mockContainerServer}
	ok, _ := l.ValidateAndBootstrapContainer(tables[0].container)
	if !ok {
		t.Errorf("Container not properly bootstrapped")
	}
}
