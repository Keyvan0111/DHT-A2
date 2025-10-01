## Description
This readme describes how to run a DHT written in golang using the gin web framework.

## How to run
The following code will initiate the cluster with a given amount of nodes with a CLI argument. It produces log files where the server output is piped into. Run clean.sh script to clean this after a run.

```sh
./clean.sh && ./run.sh NUM_NODES
# Example: ./run.sh 16
```
