package config

import (
	"github.com/BaritoLog/go-boilerplate/envkit"
)

const (
	DefaultLXDSocketPath                   = "/var/snap/lxd/common/lxd/unix.socket"
	DefaultPfCluster                       = "default"
	DefaultPfServerAddr                    = "http://192.168.33.1:3000"
	DefaultPfRegisterPath                  = "api/v1/node/register"
	DefaultPfStoreMetricsPath              = "api/v1/node/nodes/store_metrics"
	DefaultPfListScheduledContainersPath   = "api/v2/node/containers/scheduled"
	DefaultPfListProvisionedContainersPath = "api/v2/node/containers/provisioned"
	DefaultPfUpdateIpaddressPath           = "api/v1/node/containers/ipaddress"
	DefaultPfMarkProvisionedPath           = "api/v1/node/containers/mark_provisioned"
	DefaultPfMarkProvisionErrorPath        = "api/v1/node/containers/mark_provision_error"
	DefaultPfMarkBootstrappedPath          = "api/v2/node/containers/mark_bootstrapped"
	DefaultPfMarkBootstrapErrorPath        = "api/v2/node/containers/mark_bootstrap_error"
	DefaultPfMarkDeletedPath               = "api/v1/node/containers/mark_deleted"
	DefaultChefInstaller                   = "https://www.chef.io/chef/install.sh"
	DefaultChefVersion                     = "14.12.3"
	DefaultAbsoluteBootstrapScriptPath     = "/opt/bootstrap.sh"

	EnvLXDSocketPath                   = "LXD_SOCKET_PATH"
	EnvPfCluster                       = "PF_CLUSTER"
	EnvPfClusterPassword               = "PF_CLUSTER_PASSWORD"
	EnvPfServerAddr                    = "PF_SERVER_ADDR"
	EnvPfRegisterPath                  = "PF_REGISTER_PATH"
	EnvPfStoreMetricsPath              = "PF_STORE_METRICS_PATH"
	EnvPfListScheduledContainersPath   = "PF_LIST_SCHEDULED_CONTAINERS_PATH"
	EnvPfListProvisionedContainersPath = "PF_LIST_PROVISIONED_CONTAINERS_PATH"
	EnvPfUpdateIpaddressPath           = "PF_UPDATE_IPADDRESS_PATH"
	EnvPfMarkProvisionedPath           = "PF_MARK_PROVISIONED_PATH"
	EnvPfMarkProvisionErrorPath        = "PF_MARK_PROVISION_ERROR_PATH"
	EnvPfMarkBootstrappedPath          = "PF_MARK_BOOTSTRAPPED_PATH"
	EnvPfMarkBootstrapErrorPath        = "PF_MARK_BOOTSTRAP_ERROR_PATH"
	EnvPfMarkDeletedPath               = "PF_MARK_DELETED_PATH"
	EnvChefInstaller                   = "CHEF_INSTALLER"
	EnvChefVersion                     = "CHEF_VERSION"
	EnvAbsoluteBootstrapScriptPath     = "ABSOLUTE_BOOTSTRAP_SCRIPT_PATH"
)

var (
	LXDSocketPath               string
	PfCluster                   string
	PfClusterPassword           string
	PfServerAddr                string
	PfApiPath                   map[string]string
	ChefInstaller               string
	ChefVersion                 string
	AbsoluteBootstrapScriptPath string
)

func init() {
	LXDSocketPath, _ = envkit.GetString(EnvLXDSocketPath, DefaultLXDSocketPath)
	PfCluster, _ = envkit.GetString(EnvPfCluster, DefaultPfCluster)
	PfClusterPassword, _ = envkit.GetString(EnvPfClusterPassword, "")
	PfServerAddr, _ = envkit.GetString(EnvPfServerAddr, DefaultPfServerAddr)
	PfApiPath = make(map[string]string)
	PfApiPath["Register"], _ = envkit.GetString(EnvPfRegisterPath, DefaultPfRegisterPath)
	PfApiPath["StoreMetrics"], _ = envkit.GetString(EnvPfStoreMetricsPath, DefaultPfStoreMetricsPath)
	PfApiPath["ListScheduledContainers"], _ = envkit.GetString(EnvPfListScheduledContainersPath, DefaultPfListScheduledContainersPath)
	PfApiPath["ListProvisionedContainers"], _ = envkit.GetString(EnvPfListProvisionedContainersPath, DefaultPfListProvisionedContainersPath)
	PfApiPath["UpdateIpaddress"], _ = envkit.GetString(EnvPfUpdateIpaddressPath, DefaultPfUpdateIpaddressPath)
	PfApiPath["MarkProvisioned"], _ = envkit.GetString(EnvPfMarkProvisionedPath, DefaultPfMarkProvisionedPath)
	PfApiPath["MarkProvisionError"], _ = envkit.GetString(EnvPfMarkProvisionErrorPath, DefaultPfMarkProvisionErrorPath)
	PfApiPath["MarkBootstrapped"], _ = envkit.GetString(EnvPfMarkBootstrappedPath, DefaultPfMarkBootstrappedPath)
	PfApiPath["MarkBootstrapError"], _ = envkit.GetString(EnvPfMarkBootstrapErrorPath, DefaultPfMarkBootstrapErrorPath)
	PfApiPath["MarkDeleted"], _ = envkit.GetString(EnvPfMarkDeletedPath, DefaultPfMarkDeletedPath)
	ChefInstaller, _ = envkit.GetString(EnvChefInstaller, DefaultChefInstaller)
	ChefVersion, _ = envkit.GetString(EnvChefVersion, DefaultChefVersion)
	AbsoluteBootstrapScriptPath, _ = envkit.GetString(EnvAbsoluteBootstrapScriptPath, DefaultAbsoluteBootstrapScriptPath)
}
