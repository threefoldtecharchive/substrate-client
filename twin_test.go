package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwin(t *testing.T) {
	var twinID uint32
	var twin *Twin
	var err error

	cl := startLocalConnection(t)
	defer cl.Close()

	t.Run("TestCreateTwin", func(t *testing.T) {
		twinID = assertCreateTwin(t, cl, AccountBob)
	})

	t.Run("TestGetTwin", func(t *testing.T) {
		twin, err = cl.GetTwin(twinID)

		require.NoError(t, err)
		require.Equal(t, twinID, uint32(twin.ID))
	})

	t.Run("TestGetTwinByPubKey", func(t *testing.T) {
		ID, err := cl.GetTwinByPubKey(twin.Account.PublicKey())
		require.NoError(t, err)

		require.Equal(t, uint32(twin.ID), ID)
	})
}
