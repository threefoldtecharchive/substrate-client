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

	twinID = assertCreateTwin(t, cl, AccountAlice)

	twin, err = cl.GetTwin(twinID)

	require.NoError(t, err)
	require.Equal(t, twinID, uint32(twin.ID))

	ID, err := cl.GetTwinByPubKey(twin.Account.PublicKey())
	require.NoError(t, err)

	require.Equal(t, uint32(twin.ID), ID)
}
