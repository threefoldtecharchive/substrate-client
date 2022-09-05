package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	account, _ = FromAddress("5Ey64mTdp2kz7paMvRYErntweFF4utuda2VqF3GABiupRRog")
	validTwin  = Twin{
		Versioned: Versioned{
			Version: 1,
		},
		ID:       256,
		Account:  account,
		IP:       "::1",
		Entities: nil,
	}
)

func TestGetTwin(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	twn, err := cl.GetTwin(uint32(validTwin.ID))

	require.NoError(t, err)
	require.Equal(t, &validTwin, twn)

}

func TestGetTwinByPubKey(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	twnID, err := cl.GetTwinByPubKey(validTwin.Account.PublicKey())

	require.NoError(t, err)
	require.Equal(t, uint32(validTwin.ID), twnID)

}
