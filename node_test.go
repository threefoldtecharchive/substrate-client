package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode(t *testing.T) {
	var nodeID uint32
	var node *Node

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID = assertCreateNode(t, cl, farmID, twinID, identity)

	node, err = cl.GetNode(nodeID)
	require.NoError(t, err)
	require.Equal(t, twinID, uint32(node.TwinID))
	require.Equal(t, farmID, uint32(node.FarmID))

	nodeID, err = cl.GetNodeByTwinID(uint32(node.TwinID))
	require.NoError(t, err)
	require.Equal(t, uint32(node.ID), nodeID)

}

func TestGetNodesByFarmID(t *testing.T) {
	var nodeID uint32
	var node *Node

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID = assertCreateNode(t, cl, farmID, twinID, identity)

	node, err = cl.GetNode(nodeID)
	require.NoError(t, err)
	require.Equal(t, twinID, uint32(node.TwinID))
	require.Equal(t, farmID, uint32(node.FarmID))

	nodes, err := cl.GetNodesByFarmID(farmID)
	require.NoError(t, err)

	// we can verify that this node is in this list
	for _, id := range nodes {
		if id == nodeID {
			return
		}
	}

	require.Fail(t, "node was not in list of farm nodes")
}

func TestNodeSetPowerState(t *testing.T) {
	var nodeID uint32
	var node *Node

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	farmID, twinID := assertCreateFarm(t, cl)

	nodeID = assertCreateNode(t, cl, farmID, twinID, identity)

	_, err = cl.SetNodePowerState(identity, PowerState{IsUp: true})
	require.NoError(t, err)

	node, err = cl.GetNode(nodeID)
	require.NoError(t, err)

	require.True(t, node.Power.State.IsUp)
	require.False(t, node.Power.State.IsDown)

	_, err = cl.SetNodePowerState(identity, PowerState{IsDown: true, AsDown: 100})
	require.NoError(t, err)

	node, err = cl.GetNode(nodeID)
	require.NoError(t, err)

	require.True(t, node.Power.State.IsDown)
	require.EqualValues(t, node.Power.State.AsDown, 100)
	require.False(t, node.Power.State.IsUp)
}
