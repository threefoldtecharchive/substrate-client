package substrate

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

var (
	someDocumentUrl = "somedocument"
	testName        = "test-substrate"
	ip              = net.ParseIP("201:1061:b395:a8e3:5a0:f481:1102:e85a")
	AliceMnemonics  = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"
	BobMnemonics    = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Bob"
	AliceAddress    = "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	BobAddress      = "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"
	documentLink    = "somedocumentlink"
	documentHash    = "thedocumenthash"
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

func assertCreateTwin(t *testing.T, cl *Substrate, phrase string, address string) uint32 {
	identity, err := NewIdentityFromSr25519Phrase(phrase)
	require.NoError(t, err)

	account, err := FromAddress(address)
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
		log.Debug().Msgf("%s", err)
		twnID, err = cl.CreateTwin(identity, ip)
		require.NoError(t, err)
	}

	return twnID
}

func assertCreateFarm(t *testing.T, cl *Substrate) (uint32, uint32) {
	identity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	twnID := assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)

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
		twinID,
		farmID,
		Resources{
			SRU: types.U64(1024 * Gigabyte),
			MRU: types.U64(16 * Gigabyte),
			CRU: types.U64(8),
			HRU: types.U64(1024 * Gigabyte),
		},
		Location{
			City:      "SomeCity",
			Country:   "SomeCountry",
			Latitude:  "51.049999",
			Longitude: "3.733333",
		},
		[]Interface{},
		false,
		false,
		OptionBoardSerial{
			HasValue: true,
			AsValue:  "some_serial",
		})
	require.NoError(t, err)

	return nodeID
}
