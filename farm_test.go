package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFarm(t *testing.T) {
	var twinID, farmID uint32

	cl := startLocalConnection(t)
	defer cl.Close()

	t.Run("TestCreateFarm", func(t *testing.T) {
		farmID, twinID = assertCreateFarm(t, cl)
	})

	t.Run("TestGetFarm", func(t *testing.T) {
		farm, err := cl.GetFarm(farmID)

		require.NoError(t, err)
		require.Equal(t, testName, farm.Name)
		require.Equal(t, twinID, uint32(farm.TwinID))
	})
}
