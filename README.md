# Pathfinder Agent

Agent for Pathfinder container manager. Ensure appropriate containers are running on the node in which this agent reside in.

## Getting Started

Precompiled binaries are availabe on the release page, [here][pathfinder-agent-releases]

Please ensure pathfinder server is up and running and also exports appropriate configurations before starting the agent (see configurations section).

Example:
```
export PF_SERVER_ADDR=http://127.0.0.1
export PF_CLUSTER=default
export PF_CLUSTER_PASSWORD=ubuntu
```

Run the downloaded binary normally.

## Development Setup

1. Ensure that you have golang installed, with version >= 1.11.4 (because this project uses go modules).
2. Run `go build`

### Running tests

Run `go test ./...`

## Configurations

These are possible configurations that you can set via environment variables.

```
LXD_SOCKET_PATH
PF_SERVER_ADDR
PF_CLUSTER
PF_CLUSTER_PASSWORD
BOOTSTRAP_INSTALLER_URL
BOOTSTRAP_VERSION
BOOTSTRAP_FLAG_OPTIONS
```

## Getting Help

If you have any questions or feedback regarding pathfinder-agent:

- [File an issue](https://github.com/pathfinder-cm/pathfinder-agent/issues/new) for bugs, issues and feature suggestions.

Your feedback is always welcome.

## Further Reading

- [Pathfinder Container Manager Wiki][pathfinder-cm-wiki]

[pathfinder-cm-wiki]: https://github.com/pathfinder-cm/wiki
[pathfinder-agent-releases]: https://github.com/pathfinder-cm/pathfinder-agent/releases

## License

Apache License v2, see [LICENSE](LICENSE).
