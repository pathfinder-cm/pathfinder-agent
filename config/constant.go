package config

import (
	"github.com/BaritoLog/go-boilerplate/envkit"
)

const (
	DefaultLXDSocketPath                          = "/var/snap/lxd/common/lxd/unix.socket"
	DefaultPfCluster                              = "default"
	DefaultPfServerAddr                           = "http://192.168.33.1:3000"
	DefaultPfRegisterPath                         = "api/v1/node/register"
	DefaultPfStoreMetricsPath                     = "api/v1/node/nodes/store_metrics"
	DefaultPfListScheduledContainersPath          = "api/v2/node/containers/scheduled"
	DefaultPfListBootstrapScheduledContainersPath = "api/v2/node/containers/bootstrap_scheduled"
	DefaultPfUpdateIpaddressPath                  = "api/v1/node/containers/ipaddress"
	DefaultPfMarkProvisionedPath                  = "api/v1/node/containers/mark_provisioned"
	DefaultPfMarkProvisionErrorPath               = "api/v1/node/containers/mark_provision_error"
	DefaultPfMarkBootstrapStartedPath             = "api/v2/node/containers/mark_bootstrap_started"
	DefaultPfMarkBootstrappedPath                 = "api/v2/node/containers/mark_bootstrapped"
	DefaultPfMarkBootstrapErrorPath               = "api/v2/node/containers/mark_bootstrap_error"
	DefaultPfMarkDeletedPath                      = "api/v1/node/containers/mark_deleted"
	DefaultBootstrapInstallerUrl                  = ""
	DefaultBootstrapVersion                       = ""
	DefaultBootstrapFlagOptions                   = ""
	DefaultAbsoluteBootstrapScriptPath            = "/opt/bootstrap.sh"
	DefaultBootstrapContainerMaxRetry             = 2
	DefaultBootstrapMaxConcurrent                 = 5
	DefaultMetricsZpoolName                       = "local"

	EnvLXDSocketPath                          = "LXD_SOCKET_PATH"
	EnvPfCluster                              = "PF_CLUSTER"
	EnvPfClusterPassword                      = "PF_CLUSTER_PASSWORD"
	EnvPfServerAddr                           = "PF_SERVER_ADDR"
	EnvPfRegisterPath                         = "PF_REGISTER_PATH"
	EnvPfStoreMetricsPath                     = "PF_STORE_METRICS_PATH"
	EnvPfListScheduledContainersPath          = "PF_LIST_SCHEDULED_CONTAINERS_PATH"
	EnvPfListBootstrapScheduledContainersPath = "PF_LIST_BOOTSTRAP_SCHEDULED_CONTAINERS_PATH"
	EnvPfUpdateIpaddressPath                  = "PF_UPDATE_IPADDRESS_PATH"
	EnvPfMarkProvisionedPath                  = "PF_MARK_PROVISIONED_PATH"
	EnvPfMarkProvisionErrorPath               = "PF_MARK_PROVISION_ERROR_PATH"
	EnvPfMarkBootstrapStartedPath             = "PF_MARK_BOOTSTRAP_STARTED_PATH"
	EnvPfMarkBootstrappedPath                 = "PF_MARK_BOOTSTRAPPED_PATH"
	EnvPfMarkBootstrapErrorPath               = "PF_MARK_BOOTSTRAP_ERROR_PATH"
	EnvPfMarkDeletedPath                      = "PF_MARK_DELETED_PATH"
	EnvBootstrapInstallerUrl                  = "BOOTSTRAP_INSTALLER_URL"
	EnvBootstrapVersion                       = "BOOTSTRAP_VERSION"
	EnvBootstrapFlagOptions                   = "BOOTSTRAP_FLAG_OPTIONS"
	EnvAbsoluteBootstrapScriptPath            = "ABSOLUTE_BOOTSTRAP_SCRIPT_PATH"
	EnvBootstrapContainerMaxRetry             = "BOOTSTRAP_CONTAINER_MAX_RETRY"
	EnvBootstrapMaxConcurrent                 = "BOOTSTRAP_MAX_CONCURRENT"
	EnvMetricsZpoolName                       = "METRICS_ZPOOL_NAME"
)

var (
	LXDSocketPath               string
	PfCluster                   string
	PfClusterPassword           string
	PfServerAddr                string
	PfApiPath                   map[string]string
	BootstrapInstallerUrl       string
	BootstrapVersion            string
	BootstrapFlagOptions        string
	AbsoluteBootstrapScriptPath string
	BootstrapContainerMaxRetry  int
	BootstrapMaxConcurrent      int
	MetricsZpoolName            string
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
	PfApiPath["ListBootstrapScheduledContainers"], _ = envkit.GetString(EnvPfListBootstrapScheduledContainersPath, DefaultPfListBootstrapScheduledContainersPath)
	PfApiPath["UpdateIpaddress"], _ = envkit.GetString(EnvPfUpdateIpaddressPath, DefaultPfUpdateIpaddressPath)
	PfApiPath["MarkProvisioned"], _ = envkit.GetString(EnvPfMarkProvisionedPath, DefaultPfMarkProvisionedPath)
	PfApiPath["MarkProvisionError"], _ = envkit.GetString(EnvPfMarkProvisionErrorPath, DefaultPfMarkProvisionErrorPath)
	PfApiPath["MarkBootstrapStarted"], _ = envkit.GetString(EnvPfMarkBootstrapStartedPath, DefaultPfMarkBootstrapStartedPath)
	PfApiPath["MarkBootstrapped"], _ = envkit.GetString(EnvPfMarkBootstrappedPath, DefaultPfMarkBootstrappedPath)
	PfApiPath["MarkBootstrapError"], _ = envkit.GetString(EnvPfMarkBootstrapErrorPath, DefaultPfMarkBootstrapErrorPath)
	PfApiPath["MarkDeleted"], _ = envkit.GetString(EnvPfMarkDeletedPath, DefaultPfMarkDeletedPath)
	BootstrapInstallerUrl, _ = envkit.GetString(EnvBootstrapInstallerUrl, DefaultBootstrapInstallerUrl)
	BootstrapVersion, _ = envkit.GetString(EnvBootstrapVersion, DefaultBootstrapVersion)
	BootstrapFlagOptions, _ = envkit.GetString(EnvBootstrapFlagOptions, DefaultBootstrapFlagOptions)
	AbsoluteBootstrapScriptPath, _ = envkit.GetString(EnvAbsoluteBootstrapScriptPath, DefaultAbsoluteBootstrapScriptPath)
	BootstrapContainerMaxRetry, _ = envkit.GetInt(EnvBootstrapContainerMaxRetry, DefaultBootstrapContainerMaxRetry)
	BootstrapMaxConcurrent, _ = envkit.GetInt(EnvBootstrapMaxConcurrent, DefaultBootstrapMaxConcurrent)
	MetricsZpoolName, _ = envkit.GetString(EnvMetricsZpoolName, DefaultMetricsZpoolName)
}
