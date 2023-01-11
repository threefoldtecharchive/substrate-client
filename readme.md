# TFchain go client

- This library is a go implementation of a client for the TFChain.
- Internally, our own [fork](https://github.com/threefoldtech/go-substrate-rpc-client) of <https://github.com/centrifuge/go-substrate-rpc-client> is used to make substrate rpc calls.
- Used in multiple repos like [zos](https://github.com/threefoldtech/zos), [rmb-go](https://github.com/threefoldtech/rmb_go), and [terraform-provider-gric](https://github.com/threefoldtech/terraform-provider-grid).

## Usage

To make substrate calls:

- First, start a substrate connection against the desired url for the chain:

  ```go
  manager := NewManager("wss://tfchain.grid.tf/ws")
  substrateConnection, err := manager.Substrate()
  ```

- These are the urls for different chain networks:
  
  - devnet:  <wss://tfchain.dev.grid.tf/ws>
  - testnet: <wss://tfchain.test.grid.tf/ws>
  - qanet:   <wss://tfchain.qa.grid.tf/ws>
  - mainnet: <wss://tfchain.grid.tf/ws>

- It is the user's responsibility to close the connection.

  ```go
  defer substrateConnection.Close()
  ```

- Then, a user could use the provided api calls to communicate with the chain. like:

  ```go
  contractID, err := substrateConnection.CreateNodeContract(identity, nodeID, body, hash, publicIPsCount, solutionProviderID)
  ```

- Also, if a connection is closed for some reason like timing out, internally, it is reopened if nothing blocks.
- All provided api calls are found under the Substrate struct.

## Run tests

To run the tests, execute:

  ```bash
  go test . -v
  ```

- **coverage**: 30.6% of statements

## Workflows

- ### Test
  
  - This workflow runs the tests found in the root directory against a local docker image of the [TFChain](https://github.com/threefoldtech/tfchain), found [here](https://hub.docker.com/r/threefolddev/tfchain) or against devnet.
