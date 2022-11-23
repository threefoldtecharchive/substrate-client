package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	require := require.New(t)

	address := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"

	account, err := FromAddress(address)
	require.NoError(err)

	require.Equal(address, account.String())
}

func TestGetAccountByAddress(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	require := require.New(t)

	address := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	account, err := FromAddress(address)
	require.NoError(err)

	_, err = cl.GetAccountPublicInfo(account)
	require.NoError(err)
}

func TestGetBalanceByAddress(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	require := require.New(t)

	address := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	account, err := FromAddress(address)
	require.NoError(err)

	_, err = cl.GetBalance(account)
	require.NoError(err)
}
