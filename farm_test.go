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
		CertificationType: FarmCertification{
			isNotCertified: true,
		},
	}
)

func TestGetFarm(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	frm, err := cl.GetFarm(uint32(validFarm.ID))

	require.NoError(t, err)
	require.Equal(t, &validFarm, frm)

}

func TestCreateFarm(t *testing.T) {

	cl := startLocalConnection(t)
	defer cl.Close()

	frmID := assertCreateFarm(t, cl)

	frm, err := cl.GetFarm(frmID)

	require.NoError(t, err)
	require.Equal(t, testName, frm.Name)
}
