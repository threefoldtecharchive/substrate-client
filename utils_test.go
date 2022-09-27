package substrate

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testName       = "test-substrate"
	ip             = net.ParseIP("201:1061:b395:a8e3:5a0:f481:1102:e85a")
	AliceMnemonics = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"
	AliceAddress   = "5Engs9f8Gk6JqvVWz3kFyJ8Kqkgx7pLi8C1UTcr7EZ855wTQ"
)

func startLocalConnection(t *testing.T) *Substrate {
	var mgr Manager
	if _, ok := os.LookupEnv("CI"); ok {
		mgr = NewManager("ws://127.0.0.1:9944")
	} else {
		mgr = NewManager("wss://tfchain.dev.grid.tf")
	}

	cl, err := mgr.Substrate()

	require.NoError(t, err)

	return cl
}

func assertCreateTwin(t *testing.T, cl *Substrate) uint32 {

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	account, err := FromAddress(AliceAddress)
	require.NoError(t, err)

	termsAndConditions, err := cl.SignedTermsAndConditions(account)

	if len(termsAndConditions) == 0 {
		err = cl.AcceptTermsAndConditions(identity, "", "")
		require.NoError(t, err)
	}

	twnID, err := cl.GetTwinByPubKey(account.PublicKey())

	if err != nil {
		twnID, err = cl.CreateTwin(identity, ip)
		require.NoError(t, err)
	}

	return twnID
}

func assertCreateFarm(t *testing.T, cl *Substrate) (uint32, uint32) {

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	twnID := assertCreateTwin(t, cl)

	err = cl.CreateFarm(identity, testName, []PublicIPInput{})
	require.NoError(t, err)

	return 1, twnID
}
