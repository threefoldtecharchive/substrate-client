package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZosVersion(t *testing.T) {
	zosVersion := "master"

	require := require.New(t)
	cl := startLocalConnection(t)
	defer cl.Close()

	/*t.Run("TestSetZosVersion", func(t *testing.T) {
		identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
		require.NoError(err)

		tempVersion := "master"
		zosVersion, err = cl.SetZosVersion(identity, tempVersion)
		require.NoError(err)
		require.Equal(zosVersion, tempVersion)
	})*/

	t.Run("TestSetZosVersion", func(t *testing.T) {
		version, err := cl.GetZosVersion()
		require.NoError(err)
		require.Equal(zosVersion, version)
	})
}
