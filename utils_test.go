package substrate

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

var (
	someDocumentUrl = "somedocument"
	testName        = "test-substrate"
	ip              = net.ParseIP("201:1061:b395:a8e3:5a0:f481:1102:e85a")
	AliceMnemonics  = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"
	AliceAddress    = "5Engs9f8Gk6JqvVWz3kFyJ8Kqkgx7pLi8C1UTcr7EZ855wTQ"
	documentLink   = "somedocumentlink"
	documentHash   = "thedocumenthash"
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
	require.NoError(t, err)

	if len(termsAndConditions) == 0 {
		hash := md5.New()
		hash.Write([]byte(someDocumentUrl))
		h := hex.EncodeToString(hash.Sum(nil))
		err = cl.AcceptTermsAndConditions(identity, someDocumentUrl, h)
		require.NoError(t, err)
	}

	twnID, err := cl.GetTwinByPubKey(account.PublicKey())

	if err != nil {
		twnID, err = cl.CreateTwin(identity, ip)
		require.NoError(t, err)
	}

	return twnID
}

func assertCreateNode(t *testing.T, cl *Substrate, node Node) uint32 {
	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	var id uint32
	id, err = cl.GetNodeByTwinID(uint32(node.TwinID))
	require.NoError(t, err)

	if id == 0 {
		id, err = cl.CreateNode(identity, node)
		require.NoError(t, err)
	}

	id
}

func assertCreateFarm(t *testing.T, cl *Substrate) (uint32, uint32) {

	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	twnID := assertCreateTwin(t, cl)

	id, err := cl.GetFarmByName(testName)
	if err == nil {
		return id, twnID
	}

	if errors.Is(err, ErrNotFound) {
		err = cl.CreateFarm(identity, testName, []PublicIPInput{})
		require.NoError(t, err)
	}

	id, err = cl.GetFarmByName(testName)
	require.NoError(t, err)

	return id, twnID
}

func assertCreateNode(t *testing.T, cl *Substrate, farmID uint32, twinID uint32, identity Identity) uint32 {

	nodeID, err := cl.GetNodeByTwinID(twinID)
	if err == nil {
		return nodeID
	} else if !errors.Is(err, ErrNotFound) {
		require.NoError(t, err)
	}
	// if not found create a node.
	nodeID, err = cl.CreateNode(identity,
		Node{
			FarmID: types.U32(farmID),
			TwinID: types.U32(twinID),
			Location: Location{
				City:      "SomeCity",
				Country:   "SomeCountry",
				Latitude:  "51.049999",
				Longitude: "3.733333",
			},
			Resources: Resources{
				HRU: 9001778946048,
				SRU: 5121101905921,
				CRU: 24,
				MRU: 202802929664,
			},
			BoardSerial: OptionBoardSerial{
				HasValue: true,
				AsValue:  "some_serial",
			},
		},
	)
	require.NoError(t, err)

	return nodeID
}
