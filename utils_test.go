package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func startConnection(t *testing.T) *Substrate {
	mgr := NewManager("wss://tfchain.dev.grid.tf")

	cl, err := mgr.Substrate()

	require.NoError(t, err)

	return cl
}
