package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	identityA, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	_ = assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)

	groupID, err := cl.CreateGroup(identityA)
	require.NoError(t, err)

	err = cl.DeleteGroup(identityA, groupID)
	require.NoError(t, err)
}
