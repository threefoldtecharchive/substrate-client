package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFarm(t *testing.T) {
	var twinID, farmID uint32

	cl := startLocalConnection(t)
	defer cl.Close()

	farmID, twinID = assertCreateFarm(t, cl)

	farm, err := cl.GetFarm(farmID)

	require.NoError(t, err)
	require.Equal(t, testName, farm.Name)
	require.Equal(t, twinID, uint32(farm.TwinID))

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	publicIP := PublicIPInput{
		IP: "185.206.122.33/24",
		GW: "185.206.122.1",
	}
	err = cl.RemovePublicIpFromFarm(identity, farmID, publicIP.IP)
	require.NoError(t, err)

	err = cl.AddPublicIpToFarm(identity, farmID, publicIP)
	require.NoError(t, err)
}
