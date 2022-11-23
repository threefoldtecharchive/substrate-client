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

	t.Run("TestCreateNode", func(t *testing.T) {
		node := Node{
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
		nodeID, err = cl.CreateNode(identity, node)
		require.NoError(t, err)

	})

	node, err = cl.GetNode(nodeID)
	require.NoError(t, err)
	require.Equal(t, twinID, uint32(node.TwinID))
	require.Equal(t, farmID, uint32(node.FarmID))

	nodeID, err = cl.GetNodeByTwinID(uint32(node.TwinID))
	require.NoError(t, err)
	require.Equal(t, uint32(node.ID), nodeID)

}
