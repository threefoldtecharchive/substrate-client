package substrate

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestNameContract(t *testing.T) {
	var contractID uint64

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	assertCreateFarm(t, cl)

	t.Run("TestCreateNameContract", func(t *testing.T) {
		contractID, err = cl.CreateNameContract(identity, testName)
		require.NoError(t, err)
	})

	t.Run("TestGetContractIDByNameRegistration", func(t *testing.T) {
		nameContractID, err := cl.GetContractIDByNameRegistration(testName)
		require.NoError(t, err)
		require.Equal(t, contractID, nameContractID)
	})

	t.Run("TestCancelContract", func(t *testing.T) {
		err = cl.CancelContract(identity, contractID)
		require.NoError(t, err)
	})

}

func TestNodeContract(t *testing.T) {
	var contractID uint64
	var contract *Contract

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	t.Run("TestCreateNodeContract", func(t *testing.T) {
		nodeID, err := cl.CreateNode(identity, Node{
			FarmID: types.U32(farmID),
			TwinID: types.U32(twinID),
		})
		require.NoError(t, err)

		contractID, err = cl.CreateNodeContract(identity, nodeID, "", "", 0, nil)
		require.NoError(t, err)
	})

	t.Run("TestGetContract", func(t *testing.T) {
		contract, err = cl.GetContract(contractID)
		require.NoError(t, err)
	})

	t.Run("TestGetContractWithHash", func(t *testing.T) {
		contractIDWithHash, err := cl.GetContractWithHash(uint32(
			contract.ContractType.NodeContract.Node),
			contract.ContractType.NodeContract.DeploymentHash)

		require.NoError(t, err)
		require.Equal(t, contractID, contractIDWithHash)
	})

	err = cl.CancelContract(identity, contractID)
	require.NoError(t, err)
}

func TestGetRentContract(t *testing.T) {
	var contractID uint64

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	createdNode := Node{
		FarmID: types.U32(farmID),
		TwinID: types.U32(twinID),
	}
	nodeID, err := cl.CreateNode(identity, createdNode)

	require.NoError(t, err)

	contractID, err = cl.CreateRentContract(identity, nodeID, nil)
	require.NoError(t, err)

	cont, err := cl.GetContract(contractID)

	fmt.Println(cont)

	err = cl.CancelContract(identity, contractID)
	// require.NoError(t, err)
}
