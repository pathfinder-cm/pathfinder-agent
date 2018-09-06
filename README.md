# Pathfinder Agent

Agent for Pathfinder container manager. Ensure appropriate containers are running on the node in which this agent reside in.

## Setup

Please ensure pathfinder server up and running and set this configuration before starting the agent.

```
export PF_SERVER_ADDR=http://127.0.0.1
export PF_CLUSTER=default
export PF_CLUSTER_PASSWORD=ubuntu
```

## Configurations

These are possible configurations that you can set via environment variables.

```
"LXD_SOCKET_PATH"
"PF_SERVER_ADDR"
"PF_CLUSTER"
"PF_CLUSTER_PASSWORD"
```
