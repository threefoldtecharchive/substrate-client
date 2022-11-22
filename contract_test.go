package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestNameContract(t *testing.T) {
	var contractID uint64
	var nameContractID uint64

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
		nameContractID, err = cl.GetContractIDByNameRegistration(testName)
		require.NoError(t, err)
		require.Equal(t, contractID, nameContractID)
	})

	t.Run("TestCancelContract", func(t *testing.T) {
		err = cl.CancelContract(identity, contractID)
		require.NoError(t, err)
	})

}

func TestNodeContract(t *testing.T) {
	var nodeID uint32
	var contractID uint64
	var contractIDWithHash uint64
	var contract *Contract

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID, err = cl.CreateNode(identity,
		Node{
			FarmID: types.U32(farmID),
			TwinID: types.U32(twinID),
			Location: Location{
				City:      "SomeCity",
				Country:   "SomeCountry",
				Latitude:  "51.049999",
				Longitude: "3.733333",
			},
			Resources: Resources{
				HRU: 9001778946048,
				SRU: 5121101905921,
				CRU: 24,
				MRU: 202802929664,
			},
			BoardSerial: OptionBoardSerial{
				HasValue: true,
				AsValue:  "some_serial",
			},
		},
	)
	require.NoError(t, err)

	contractID, err = cl.CreateNodeContract(identity, nodeID, "", "", 0, nil)
	require.NoError(t, err)

	contract, err = cl.GetContract(contractID)
	require.NoError(t, err)

	contractIDWithHash, err = cl.GetContractWithHash(uint32(
		contract.ContractType.NodeContract.Node),
		contract.ContractType.NodeContract.DeploymentHash)

	require.NoError(t, err)
	require.Equal(t, contractID, contractIDWithHash)

	err = cl.CancelContract(identity, contractID)
	require.NoError(t, err)
}

func TestGetRentContract(t *testing.T) {
	var nodeID uint32
	var contractID uint64

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	t.Run("TestCreateRentContract", func(t *testing.T) {
		nodeID, err = cl.CreateNode(identity,
			Node{
				FarmID: types.U32(farmID),
				TwinID: types.U32(twinID),
				Location: Location{
					City:      "SomeCity",
					Country:   "SomeCountry",
					Latitude:  "51.049999",
					Longitude: "3.733333",
				},
				Resources: Resources{
					HRU: 9001778946048,
					SRU: 5121101905921,
					CRU: 24,
					MRU: 202802929664,
				},
				BoardSerial: OptionBoardSerial{
					HasValue: true,
					AsValue:  "some_serial",
				},
			},
		)
		require.NoError(t, err)

		contractID, err = cl.CreateRentContract(identity, nodeID, nil)
		require.NoError(t, err)
	})

	t.Run("TestGetContract", func(t *testing.T) {
		_, err = cl.GetContract(contractID)
		require.NoError(t, err)
	})

	err = cl.CancelContract(identity, contractID)
	require.NoError(t, err)
}
