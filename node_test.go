package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
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

	t.Run("TestCreateNode", func(t *testing.T) {
		nodeID, err = cl.CreateNode(identity, Node{FarmID: types.U32(farmID), TwinID: types.U32(twinID)})
		require.NoError(t, err)
	})

	t.Run("TestGetNode", func(t *testing.T) {
		node, err = cl.GetNode(nodeID)
		require.NoError(t, err)
		require.Equal(t, twinID, uint32(node.TwinID))
		require.Equal(t, farmID, uint32(node.FarmID))
	})

	t.Run("TestGetNodeByTwinID", func(t *testing.T) {
		nodeID, err = cl.GetNodeByTwinID(uint32(node.TwinID))
		require.NoError(t, err)
		require.Equal(t, uint32(node.ID), nodeID)
	})

}