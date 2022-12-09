package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestNameContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	assertCreateFarm(t, cl)

	contractID, err := cl.CreateNameContract(identity, testName)
	require.NoError(t, err)

	nameContractID, err := cl.GetContractIDByNameRegistration(testName)
	require.NoError(t, err)
	require.Equal(t, contractID, nameContractID)

	err = cl.CancelContract(identity, contractID)
	require.NoError(t, err)

}

func TestCreateDeployment(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	_ = assertCreateNode(t, cl, farmID, twinID, identity)

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
	capacityReservationID, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	deploymentID, err := cl.CreateDeployment(identity, capacityReservationID, "", "",
		Resources{
			SRU: types.U64(512 * Gigabyte),
			MRU: types.U64(8 * Gigabyte),
			CRU: types.U64(4),
			HRU: types.U64(512 * Gigabyte),
		}, 0)
	require.NoError(t, err)

	capacityReservationContract, err := cl.GetContract(capacityReservationID)
	require.NoError(t, err)

	deployment, err := cl.GetDeployment(deploymentID)
	require.NoError(t, err)

	contractIDWithHash, err := cl.GetContractWithHash(uint32(
		capacityReservationContract.ContractType.CapacityReservationContract.NodeID),
		deployment.DeploymentHash)
	require.NoError(t, err)
	require.Equal(t, deploymentID, contractIDWithHash)

	contractResources, err := cl.GetContractResources(capacityReservationID)
	require.NoError(t, err)
	require.Equal(t, contractResources.UsedResources, deployment.Resources)

	err = cl.CancelDeployment(identity, deploymentID)
	require.NoError(t, err)

	err = cl.CancelContract(identity, capacityReservationID)
	require.NoError(t, err)
}

func TestUpdateDeployment(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	_ = assertCreateNode(t, cl, farmID, twinID, identity)

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
	capacityReservationID, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	resources := Resources{
		SRU: types.U64(256 * Gigabyte),
		MRU: types.U64(4 * Gigabyte),
		CRU: types.U64(2),
		HRU: types.U64(256 * Gigabyte),
	}
	deploymentID, err := cl.CreateDeployment(identity, capacityReservationID, "", "", resources, 0)
	require.NoError(t, err)

	deployment, err := cl.GetDeployment(deploymentID)
	require.NoError(t, err)
	require.Equal(t, resources, deployment.Resources)

	resourcesUpdated := Resources{
		SRU: types.U64(512 * Gigabyte),
		MRU: types.U64(8 * Gigabyte),
		CRU: types.U64(4),
		HRU: types.U64(512 * Gigabyte),
	}
	err = cl.UpdateDeployment(identity, deploymentID, "", "", &resourcesUpdated)
	require.NoError(t, err)

	deployment, err = cl.GetDeployment(deploymentID)
	require.NoError(t, err)
	require.Equal(t, resourcesUpdated, deployment.Resources)

	err = cl.CancelDeployment(identity, deploymentID)
	require.NoError(t, err)

	err = cl.CancelContract(identity, capacityReservationID)
	require.NoError(t, err)
}

func TestCreateCapacityReservationContractPolicyAny(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID := assertCreateNode(t, cl, farmID, twinID, identity)

	resources := Resources{
		SRU: types.U64(512 * Gigabyte),
		MRU: types.U64(8 * Gigabyte),
		CRU: types.U64(4),
		HRU: types.U64(512 * Gigabyte),
	}
	policy := CapacityReservationPolicy{
		IsAny:       true,
		IsExclusive: false,
		IsNode:      false,
		AsAny: Any{
			Resources: resources,
		},
	}
	contractID, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	contract, err := cl.GetContract(contractID)
	require.NoError(t, err)
	require.Equal(t, contract.ContractType.CapacityReservationContract.NodeID, types.U32(nodeID))

	contractResources, err := cl.GetContractResources(contractID)
	require.NoError(t, err)
	require.Equal(t, contractResources.TotalResources, resources)

	err = cl.CancelContract(identity, contractID)
}

func TestCreateCapacityReservationContractPolicyNode(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID := assertCreateNode(t, cl, farmID, twinID, identity)

	policy := CapacityReservationPolicy{
		IsAny:       false,
		IsExclusive: false,
		IsNode:      true,
		AsNode: NodePolicy{
			NodeID: types.U32(nodeID),
		},
	}
	contractID, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	nodeResources, err := cl.GetNodeResources(nodeID)
	require.NoError(t, err)

	contractResources, err := cl.GetContractResources(contractID)
	require.NoError(t, err)
	require.Equal(t, contractResources.TotalResources, nodeResources.TotalResources)

	err = cl.CancelContract(identity, contractID)
}

func TestCreateCapacityReservationContractPolicyExclusive(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identityA, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	identityB, err := NewIdentityFromSr25519Phrase(BobMnemonics)
	require.NoError(t, err)

	farmID, twinIDA := assertCreateFarm(t, cl)
	twinIDB := assertCreateTwin(t, cl, BobMnemonics, BobAddress)

	nodeIDA := assertCreateNode(t, cl, farmID, twinIDA, identityA)
	nodeIDB := assertCreateNode(t, cl, farmID, twinIDB, identityB)

	groupID, err := cl.CreateGroup(identityA)
	require.NoError(t, err)

	policy := CapacityReservationPolicy{
		IsAny:       false,
		IsExclusive: true,
		IsNode:      false,
		AsExclusive: Exclusive{
			GroupID: types.U32(groupID),
			Resources: Resources{
				SRU: types.U64(512 * Gigabyte),
				MRU: types.U64(8 * Gigabyte),
				CRU: types.U64(4),
				HRU: types.U64(512 * Gigabyte),
			},
		},
	}
	contractIDA, err := cl.CreateCapacityReservationContract(identityA, farmID, policy, nil)
	require.NoError(t, err)
	contractIDB, err := cl.CreateCapacityReservationContract(identityB, farmID, policy, nil)
	require.NoError(t, err)

	nodeResourcesA, err := cl.GetNodeResources(nodeIDA)
	require.NoError(t, err)
	nodeResourcesB, err := cl.GetNodeResources(nodeIDB)
	require.NoError(t, err)

	contractA, err := cl.GetContract(contractIDA)
	require.NoError(t, err)
	contractResourcesA, err := cl.GetContractResources(contractIDA)
	require.NoError(t, err)
	require.Equal(t, contractA.ContractType.CapacityReservationContract.NodeID, types.U32(nodeIDA))
	require.Equal(t, contractResourcesA.TotalResources, nodeResourcesA.UsedResources)
	contractB, err := cl.GetContract(contractIDB)
	require.NoError(t, err)
	contractResourcesB, err := cl.GetContractResources(contractIDB)
	require.NoError(t, err)
	require.Equal(t, contractB.ContractType.CapacityReservationContract.NodeID, types.U32(nodeIDB))
	require.Equal(t, contractResourcesB.TotalResources, nodeResourcesB.UsedResources)

	err = cl.CancelContract(identityA, contractIDA)
	require.NoError(t, err)
	err = cl.CancelContract(identityB, contractIDB)
	require.NoError(t, err)

	err = cl.DeleteGroup(identityA, groupID)
	require.NoError(t, err)
}

func TestUpdateCapacityReservationContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	_ = assertCreateNode(t, cl, farmID, twinID, identity)

	resources := Resources{
		SRU: types.U64(512 * Gigabyte),
		MRU: types.U64(8 * Gigabyte),
		CRU: types.U64(4),
		HRU: types.U64(512 * Gigabyte),
	}
	policy := CapacityReservationPolicy{
		IsAny:       true,
		IsExclusive: false,
		IsNode:      false,
		AsAny: Any{
			Resources: resources,
		},
	}
	contractID, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	contractResources, err := cl.GetContractResources(contractID)
	require.NoError(t, err)
	require.Equal(t, contractResources.TotalResources, resources)

	resourcesUpdate := Resources{
		SRU: types.U64(1024 * Gigabyte),
		MRU: types.U64(16 * Gigabyte),
		CRU: types.U64(8),
		HRU: types.U64(1024 * Gigabyte),
	}
	err = cl.UpdateCapacityReservationContract(identity, contractID, resourcesUpdate)
	require.NoError(t, err)

	contractResources, err = cl.GetContractResources(contractID)
	require.NoError(t, err)
	require.Equal(t, contractResources.TotalResources, resourcesUpdate)

	err = cl.CancelContract(identity, contractID)
}

func TestMultipleCapacityReservationContracts(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID := assertCreateNode(t, cl, farmID, twinID, identity)

	resources := Resources{
		SRU: types.U64(512 * Gigabyte),
		MRU: types.U64(8 * Gigabyte),
		CRU: types.U64(4),
		HRU: types.U64(512 * Gigabyte),
	}
	policy := CapacityReservationPolicy{
		IsAny:       true,
		IsExclusive: false,
		IsNode:      false,
		AsAny: Any{
			Resources: resources,
		},
	}
	contractID1, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	contractID2, err := cl.CreateCapacityReservationContract(identity, farmID, policy, nil)
	require.NoError(t, err)

	require.NotEqual(t, contractID1, contractID2)

	contract1, err := cl.GetContract(contractID1)
	require.NoError(t, err)
	require.Equal(t, contract1.ContractType.CapacityReservationContract.NodeID, types.U32(nodeID))

	contract2, err := cl.GetContract(contractID2)
	require.NoError(t, err)
	require.Equal(t, contract2.ContractType.CapacityReservationContract.NodeID, types.U32(nodeID))

	contractIDs, err := cl.GetCapacityReservationContracts(nodeID)
	require.NoError(t, err)
	require.Equal(t, []uint64{contractID1, contractID2}, contractIDs)

	err = cl.CancelContract(identity, contractID1)
	err = cl.CancelContract(identity, contractID2)

}
