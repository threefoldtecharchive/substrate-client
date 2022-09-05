package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

var (
	nameContract = Contract{
		Versioned: Versioned{
			Version: 4,
		},
		State: ContractState{
			IsCreated: true,
		},
		ContractID: 7399,
		TwinID:     256,
		ContractType: ContractType{
			IsNameContract: true,
			NameContract: NameContract{
				Name: "substrate-testing",
			},
		},
		SolutionProviderID: types.OptionU64{},
	}
	nodeContract = Contract{
		ContractID: 7406,
		TwinID:     256,
		ContractType: ContractType{
			IsNodeContract: true,
			NodeContract: NodeContract{
				Node: 14,
			},
		},
		SolutionProviderID: types.OptionU64{},
	}
)

func TestGetContract(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	contract, err := cl.GetContract(uint64(nameContract.ContractID))

	require.NoError(t, err)
	require.Equal(t, &nameContract, contract)
}

func TestGetContractIDByNameRegistration(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	contractID, err := cl.GetContractIDByNameRegistration(nameContract.ContractType.NameContract.Name)

	require.NoError(t, err)
	require.Equal(t, uint64(nameContract.ContractID), contractID)
}

func TestGetContractWithHash(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	contractID, err := cl.GetContractWithHash(uint32(nodeContract.ContractType.NodeContract.Node),
		nodeContract.ContractType.NodeContract.DeploymentHash)

	require.NoError(t, err)
	require.Equal(t, uint64(nodeContract.ContractID), contractID)
}

func TestGetNodeContracts(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	contracts, err := cl.GetNodeContracts(uint32(nodeContract.ContractType.NodeContract.Node))

	require.NoError(t, err)
	require.Contains(t, contracts, nodeContract.ContractID)
}
