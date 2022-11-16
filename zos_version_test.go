package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZosVersion(t *testing.T) {
	var zosVersion *string
	var err error

	cl := startLocalConnection(t)
	defer cl.Close()

	t.Run("TestGetZosVersion", func(t *testing.T) {
		zosVersion, err = cl.GetZosVersion()

		require.NoError(t, err)
		require.NotEqual(t, zosVersion, nil)
	})
}
