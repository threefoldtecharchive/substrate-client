package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	validFarm = Farm{
		Versioned:       Versioned{Version: 4},
		ID:              125,
		Name:            "substrate-test",
		TwinID:          256,
		PricingPolicyID: 1,
	}
)

func TestGetFarm(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	frm, err := cl.GetFarm(uint32(validFarm.ID))

	require.NoError(t, err)

	require.Equal(t, validFarm.Versioned, frm.Versioned)
	require.Equal(t, validFarm.ID, frm.ID)
	require.Equal(t, validFarm.Name, frm.Name)
	require.Equal(t, validFarm.TwinID, frm.TwinID)
	require.Equal(t, validFarm.PublicIPs, frm.PublicIPs)
	require.Equal(t, validFarm.DedicatedFarm, frm.DedicatedFarm)
	require.Equal(t, validFarm.PricingPolicyID, frm.PricingPolicyID)

}
