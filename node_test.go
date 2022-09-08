package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	validNode = Node{
		Versioned:       Versioned{Version: 5},
		ID:              14,
		FarmID:          1,
		TwinID:          22,
		Resources:       Resources{HRU: 9001778946048, SRU: 512110190592, CRU: 24, MRU: 202802929664},
		Location:        Location{Longitude: "4.3447", Latitude: "50.8509"},
		Country:         "Belgium",
		City:            "Unknown",
		PublicConfig:    OptionPublicConfig{HasValue: true, AsValue: PublicConfig{IP4: IP{IP: "185.206.122.31/24", GW: "185.206.122.1"}, IP6: OptionIP{HasValue: true, AsValue: IP{IP: "2a10:b600:1::0cc4:7a30:65b5/64", GW: "2a10:b600:1::1"}}, Domain: OptionDomain{HasValue: true, AsValue: "gent01.dev.grid.tf"}}},
		Created:         1649252280,
		FarmingPolicy:   1,
		Interfaces:      []Interface{{Name: "zos", Mac: "0c:c4:7a:30:7a:68", IPs: []string{"10.9.0.113", "2a10:b600:0:9:ec4:7aff:fe30:7a68"}}},
		Certification:   NodeCertification{IsDiy: true, IsCertified: false},
		SecureBoot:      false,
		Virtualized:     false,
		BoardSerial:     "NM141S001819",
		ConnectionPrice: 800,
	}
	createdNode = Node{
		Versioned: Versioned{
			Version: 5,
		},
		ID:            1,
		FarmID:        1,
		TwinID:        1,
		Created:       0,
		FarmingPolicy: 1,
		Certification: NodeCertification{
			IsDiy: true,
		},
		ConnectionPrice: 80,
	}
)

func TestGetNodeByTwinID(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	nodeID, err := cl.GetNodeByTwinID(uint32(validNode.TwinID))

	require.NoError(t, err)
	require.Equal(t, uint32(validNode.ID), nodeID)
}

func TestGetNode(t *testing.T) {

	cl := startConnection(t)
	defer cl.Close()

	node, err := cl.GetNode(uint32(validNode.ID))

	require.NoError(t, err)
	require.Equal(t, &validNode, node)
}

func TestCreateNode(t *testing.T) {

	cl := startLocalConnection(t)
	defer cl.Close()

	identity, err := NewIdentityFromSr25519Phrase(Mnemonics)
	require.NoError(t, err)

	assertCreateFarm(t, cl)

	nodeID, err := cl.CreateNode(identity, createdNode)
	node, err := cl.GetNode(nodeID)

	createdNode.Created = node.Created
	require.Equal(t, &createdNode, node)
}
