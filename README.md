# Pathfinder Agent

Agent for Pathfinder container manager. Ensure appropriate containers are running on the node in which this agent reside in.

## Minimum Setup
```
export PF_CLUSTER_PASSWORD=ubuntu
export PF_CLUSTER=default
export PF_SERVER_ADDR=http://127.0.0.1

```

## Environment Variables

```
"LXD_SOCKET_PATH"
"PF_CLUSTER"
"PF_CLUSTER_PASSWORD"
"PF_SERVER_ADDR"
"PF_REGISTER_PATH"
"PF_STORE_METRIC_PATH"
"PF_LIST_CONTAINERS_PATH"
"PF_UPDATE_IPADDRESS"
"PF_MARK_PROVISIONED_PATH"
"PF_MARK_PROVISION_ERROR_PATH"
"PF_MARK_DELETED_PATH"
```
