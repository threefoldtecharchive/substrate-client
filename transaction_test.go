package substrate

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	t.Skip("manual test")

	require := require.New(t)
	key := "begin hope purse present final sea wing someone devote drastic praise scrap"

	identity, err := NewIdentityFromSr25519Phrase(key)
	require.NoError(err)

	cl := NewManager("wss://tfchain.dev.grid.tf/ws")
	con, err := cl.Substrate()
	require.NoError(err)
	defer con.Close()

	err = con.AcceptTermsAndConditions(identity, "", "")
	require.NoError(err)

	twin, err := con.GetTwinByPubKey(identity.PublicKey())
	if errors.Is(err, ErrNotFound) {
		twin, err = con.CreateTwin(identity, net.ParseIP("10.20.30.40"))
		require.NoError(err)
	} else if err != nil {
		t.Fatal(err)
	}
	const farmID = 41

	farm, err := con.GetFarm(41)
	require.NoError(err)

	nodeId, err := con.GetNodeByTwinID(twin)
	if errors.Is(err, ErrNotFound) {
		nodeId, err = con.CreateNode(identity, Node{
			FarmID:      farm.ID,
			TwinID:      types.U32(twin),
			Virtualized: true,
		})
		require.NoError(err)
	} else if err != nil {
		t.Fatal(err)
	} else {
		// node is fine. Let's do an update
		fmt.Println("updating node")
		_, err := con.UpdateNode(identity, Node{
			ID:          types.U32(nodeId),
			FarmID:      farm.ID,
			TwinID:      types.U32(twin),
			Virtualized: true,
			Country:     fmt.Sprintf("EG-%d", time.Now().Unix()),
			City:        fmt.Sprintf("CA-%d", time.Now().Unix()),
			BoardSerial: fmt.Sprint(time.Now().Unix()),
		})

		require.NoError(err)
	}

	fmt.Println("twin id: ", twin)
	fmt.Println("node id: ", nodeId)
}
