package substrate

import (
	"testing"

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

	contractID, err = cl.CreateNameContract(identity, testName)
	require.NoError(t, err)

	nameContractID, err = cl.GetContractIDByNameRegistration(testName)
	require.NoError(t, err)
	require.Equal(t, contractID, nameContractID)

	err = cl.CancelContract(identity, contractID)
	require.NoError(t, err)

}

func TestDeploymentContract(t *testing.T) {
	var deploymentID uint64
	var deployment *Deployment

	var capacityReservationID uint64
	var capacityReservationContract *Contract

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	t.Run("TestCreateDeploymentContract", func(t *testing.T) {
		_, err := cl.CreateNode(identity, Node{
			FarmID: types.U32(farmID),
			TwinID: types.U32(twinID),
			Resources: ConsumableResources{
				TotalResources: Resources{
					SRU: types.U64(1024 * Gigabyte),
					MRU: types.U64(16 * Gigabyte),
					CRU: types.U64(8),
					HRU: types.U64(1024 * Gigabyte),
				},
			},
			Location: Location{
				City:      "Ghent",
				Country:   "Belgium",
				Longitude: "12",
				Latitude:  "15",
			},
		})
		require.NoError(t, err)

		policy := CapacityReservationPolicy{
			IsAny: true,
			AsAny: Any{
				Resources: Resources{
					SRU: types.U64(512 * Gigabyte),
					MRU: types.U64(8 * Gigabyte),
					CRU: types.U64(4),
					HRU: types.U64(512 * Gigabyte),
				},
				Features: OptionFeatures{
					HasValue: false,
				},
			},
			IsExclusive: false,
			IsNode:      false,
		}

		capacityReservationID, err = cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
		require.NoError(t, err)

		deploymentID, err = cl.CreateDeployment(identity, capacityReservationID, "", "", Resources{}, 0)
		require.NoError(t, err)
	})

	t.Run("TestGetContract", func(t *testing.T) {
		capacityReservationContract, err = cl.GetContract(capacityReservationID)
		require.NoError(t, err)
	})

	t.Run("TestGetDeployment", func(t *testing.T) {
		deployment, err = cl.GetDeployment(deploymentID)
		require.NoError(t, err)
	})

	t.Run("TestGetContractWithHash", func(t *testing.T) {
		contractIDWithHash, err := cl.GetContractWithHash(uint32(
			capacityReservationContract.ContractType.CapacityReservationContract.NodeID),
			deployment.DeploymentHash)

		require.NoError(t, err)
		require.Equal(t, deploymentID, contractIDWithHash)
	})

	err = cl.CancelDeployment(identity, deploymentID)
	require.NoError(t, err)

	err = cl.CancelContract(identity, capacityReservationID)
	require.NoError(t, err)
}

func TestCreateCapacityReservationContractPolicyNode(t *testing.T) {
	var contractID uint64

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	createdNode := Node{
		FarmID: types.U32(farmID),
		TwinID: types.U32(twinID),
		Resources: ConsumableResources{
			TotalResources: Resources{
				SRU: types.U64(1024 * Gigabyte),
				MRU: types.U64(16 * Gigabyte),
				CRU: types.U64(8),
				HRU: types.U64(1024 * Gigabyte),
			},
		},
		Location: Location{
			City:      "Ghent",
			Country:   "Belgium",
			Longitude: "12",
			Latitude:  "15",
		},
	}
	nodeID, err := cl.CreateNode(identity, createdNode)

	require.NoError(t, err)

	policy := CapacityReservationPolicy{
		IsAny:       false,
		IsExclusive: false,
		IsNode:      true,
		AsNode: NodePolicy{
			NodeID: types.U32(nodeID),
		},
	}
	contractID, err = cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	cont, err := cl.GetContract(contractID)

	fmt.Println(cont)

	err = cl.CancelContract(identity, contractID)
}
