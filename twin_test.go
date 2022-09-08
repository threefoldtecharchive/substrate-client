package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTwin(t *testing.T) {

	cl := startLocalConnection(t)
	defer cl.Close()

	twnID := assertCreateTwin(t, cl)

	twn, err := cl.GetTwin(twnID)

	require.NoError(t, err)
	require.Equal(t, twnID, uint32(twn.ID))
}

func TestGetTwinByPubKey(t *testing.T) {

	cl := startLocalConnection(t)
	defer cl.Close()

	twnID := assertCreateTwin(t, cl)

	twn, err := cl.GetTwin(twnID)
	require.NoError(t, err)

	ID, err := cl.GetTwinByPubKey(twn.Account.PublicKey())

	require.Equal(t, twnID, ID)
}
