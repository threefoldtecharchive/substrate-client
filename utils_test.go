package substrate

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testName  = "test-substrate"
	ip        = net.ParseIP("201:1061:b395:a8e3:5a0:f481:1102:e85a")
	Mnemonics = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"
)

func startConnection(t *testing.T) *Substrate {
	mgr := NewManager("wss://tfchain.dev.grid.tf")

	cl, err := mgr.Substrate()

	require.NoError(t, err)

	return cl
}

func startLocalConnection(t *testing.T) *Substrate {
	mgr := NewManager("ws://0.0.0.0:9944")

	cl, err := mgr.Substrate()

	require.NoError(t, err)

	return cl
}

func assertCreateTwin(t *testing.T, cl *Substrate) uint32 {

	identity, err := NewIdentityFromSr25519Phrase(Mnemonics)
	require.NoError(t, err)

	twnID, err := cl.CreateTwin(identity, ip)
	require.NoError(t, err)

	return twnID
}

func assertCreateFarm(t *testing.T, cl *Substrate) uint32 {

	identity, err := NewIdentityFromSr25519Phrase(Mnemonics)
	require.NoError(t, err)

	err = cl.AcceptTermsAndConditions(identity, "", "")
	require.NoError(t, err)

	assertCreateTwin(t, cl)

	err = cl.CreateFarm(identity, testName, []PublicIPInput{})
	require.NoError(t, err)

	return 1
}
