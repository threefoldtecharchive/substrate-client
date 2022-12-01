package substrate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Resources type
type Resources struct {
	HRU types.U64
	SRU types.U64
	CRU types.U64
	MRU types.U64
}

// Location type
type Location struct {
	City      string
	Country   string
	Latitude  string
	Longitude string
}

// Role type
type Role struct {
	IsNode    bool
	IsGateway bool
}

// Decode implementation for the enum type
func (r *Role) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsNode = true
	case 1:
		r.IsGateway = true
	default:
		return fmt.Errorf("unknown Role value")
	}

	return nil
}

// Encode implementation
func (r Role) Encode(encoder scale.Encoder) (err error) {
	if r.IsNode {
		err = encoder.PushByte(0)
	} else if r.IsGateway {
		err = encoder.PushByte(1)
	}

	return
}

type IP struct {
	IP string
	GW string
}

type OptionIP struct {
	HasValue bool
	AsValue  IP
}

// Encode implementation
func (m OptionIP) Encode(encoder scale.Encoder) (err error) {
	var i byte
	if m.HasValue {
		i = 1
	}
	err = encoder.PushByte(i)
	if err != nil {
		return err
	}

	if m.HasValue {
		err = encoder.Encode(m.AsValue)
	}

	return
}

// Decode implementation
func (m *OptionIP) Decode(decoder scale.Decoder) (err error) {
	var i byte
	if err := decoder.Decode(&i); err != nil {
		return err
	}

	switch i {
	case 0:
		return nil
	case 1:
		m.HasValue = true
		return decoder.Decode(&m.AsValue)
	default:
		return fmt.Errorf("unknown value for Option")
	}
}

type OptionDomain struct {
	HasValue bool
	AsValue  string
}

// Encode implementation
func (m OptionDomain) Encode(encoder scale.Encoder) (err error) {
	var i byte
	if m.HasValue {
		i = 1
	}
	err = encoder.PushByte(i)
	if err != nil {
		return err
	}

	if m.HasValue {
		err = encoder.Encode(m.AsValue)
	}

	return
}

// Decode implementation
func (m *OptionDomain) Decode(decoder scale.Decoder) (err error) {
	var i byte
	if err := decoder.Decode(&i); err != nil {
		return err
	}

	switch i {
	case 0:
		return nil
	case 1:
		m.HasValue = true
		return decoder.Decode(&m.AsValue)
	default:
		return fmt.Errorf("unknown value for Option")
	}
}

// PublicConfig type
type PublicConfig struct {
	IP4    IP
	IP6    OptionIP
	Domain OptionDomain
}

// OptionPublicConfig type
type OptionPublicConfig struct {
	HasValue bool
	AsValue  PublicConfig
}

// Encode implementation
func (m OptionPublicConfig) Encode(encoder scale.Encoder) (err error) {
	var i byte
	if m.HasValue {
		i = 1
	}
	err = encoder.PushByte(i)
	if err != nil {
		return err
	}

	if m.HasValue {
		err = encoder.Encode(m.AsValue)
	}

	return
}

// Decode implementation
func (m *OptionPublicConfig) Decode(decoder scale.Decoder) (err error) {
	var i byte
	if err := decoder.Decode(&i); err != nil {
		return err
	}

	switch i {
	case 0:
		return nil
	case 1:
		m.HasValue = true
		return decoder.Decode(&m.AsValue)
	default:
		return fmt.Errorf("unknown value for Option")
	}
}

type Interface struct {
	Name string
	Mac  string
	IPs  []string
}

// OptionBoardSerial type
type OptionBoardSerial struct {
	HasValue bool
	AsValue  string
}

// Encode implementation
func (m OptionBoardSerial) Encode(encoder scale.Encoder) (err error) {
	var i byte
	if m.HasValue {
		i = 1
	}
	err = encoder.PushByte(i)
	if err != nil {
		return err
	}

	if m.HasValue {
		err = encoder.Encode(m.AsValue)
	}

	return
}

// Decode implementation
func (m *OptionBoardSerial) Decode(decoder scale.Decoder) (err error) {
	var i byte
	if err := decoder.Decode(&i); err != nil {
		return err
	}

	switch i {
	case 0:
		return nil
	case 1:
		m.HasValue = true
		return decoder.Decode(&m.AsValue)
	default:
		return fmt.Errorf("unknown value for Option")
	}
}

// Node type
type Node struct {
	Versioned
	ID              types.U32
	FarmID          types.U32
	TwinID          types.U32
	Resources       ConsumableResources
	Location        Location
	Power           Power
	PublicConfig    OptionPublicConfig
	Created         types.U64
	FarmingPolicy   types.U32
	Interfaces      []Interface
	Certification   NodeCertification
	SecureBoot      bool
	Virtualized     bool
	BoardSerial     OptionBoardSerial
	ConnectionPrice types.U32
}

type Power struct {
	Target     PowerTarget
	State      PowerState
	LastUptime types.U64
}

type PowerTarget struct {
	IsUp   bool
	IsDown bool
}

// Decode implementation for the enum type
func (r *PowerTarget) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsUp = true
	case 1:
		r.IsDown = true
	case 2:
	default:
		return fmt.Errorf("unknown power target value")
	}

	return nil
}

// Encode implementation
func (r PowerTarget) Encode(encoder scale.Encoder) (err error) {
	if r.IsUp {
		err = encoder.PushByte(0)
	} else if r.IsDown {
		err = encoder.PushByte(1)
	}
	return
}

type PowerState struct {
	IsUp   bool
	IsDown bool
	AsDown types.U32
}

// Decode implementation for the enum type
func (r *PowerState) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsUp = true
	case 1:
		r.IsDown = true
		if err := decoder.Decode(&r.AsDown); err != nil {
			return errors.Wrap(err, "failed to get power state")
		}
	default:
		return fmt.Errorf("unknown power state value")
	}

	return nil
}

// Encode implementation
func (r PowerState) Encode(encoder scale.Encoder) (err error) {
	if r.IsUp {
		err = encoder.PushByte(0)
	} else if r.IsDown {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(r.AsDown)
	}
	return
}

type NodeFeatures struct {
	IsPublicNode bool
}

// Decode implementation for the enum type
func (r *NodeFeatures) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsPublicNode = true
	default:
		return fmt.Errorf("unknown node feature value")
	}

	return nil
}

// Encode implementation
func (r NodeFeatures) Encode(encoder scale.Encoder) (err error) {
	if r.IsPublicNode {
		err = encoder.PushByte(0)
	}
	return
}

// Eq compare changes on node settable fields
func (n *Node) Eq(o *Node) bool {
	return n.FarmID == o.FarmID &&
		n.TwinID == o.TwinID &&
		reflect.DeepEqual(n.Resources, o.Resources) &&
		reflect.DeepEqual(n.Location, o.Location) &&
		reflect.DeepEqual(n.Interfaces, o.Interfaces) &&
		n.SecureBoot == o.SecureBoot &&
		n.Virtualized == o.Virtualized &&
		reflect.DeepEqual(n.BoardSerial, o.BoardSerial)
}

type NodeExtra struct {
	Secure       bool
	Virtualized  bool
	SerialNumber OptionBoardSerial
}

// GetNodeByTwinID gets a node by twin id
func (s *Substrate) GetNodeByTwinID(twin uint32) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}
	bytes, err := types.Encode(twin)
	if err != nil {
		return 0, err
	}
	key, err := types.CreateStorageKey(meta, "TfgridModule", "NodeIdByTwinID", bytes, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}
	var id types.U32
	ok, err := cl.RPC.State.GetStorageLatest(key, &id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup entity")
	}

	if !ok || id == 0 {
		return 0, errors.Wrap(ErrNotFound, "node not found")
	}

	return uint32(id), nil
}

func (s *Substrate) GetNodesByFarmID(id uint32) ([]uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}
	key, err := types.CreateStorageKey(meta, "TfgridModule", "NodesByFarmID", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	return s.getNodesByFarmID(cl, key)
}

func (s *Substrate) getNodesByFarmID(cl Conn, key types.StorageKey) ([]uint32, error) {
	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}

	var node []uint32
	if err := types.Decode(*raw, &node); err != nil {
		return nil, errors.Wrap(err, "failed to load object")
	}
	return node, nil
}

// GetNode with id
func (s *Substrate) GetNode(id uint32) (*Node, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}
	key, err := types.CreateStorageKey(meta, "TfgridModule", "Nodes", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	return s.getNode(cl, key)
}

type ScannedNode struct {
	ID   uint32
	Node Node
	Err  error
}

func (s *Substrate) ScanNodes(ctx context.Context, from, to uint32) (<-chan ScannedNode, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}
	ch := make(chan ScannedNode)

	getNode := func(id uint32) (*Node, error) {
		bytes, err := types.Encode(id)
		if err != nil {
			return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
		}

		key, err := types.CreateStorageKey(meta, "TfgridModule", "Nodes", bytes, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create substrate query key")
		}

		return s.getNode(cl, key)
	}

	go func(from, to uint32) {
		defer close(ch)

		for ; from <= to; from++ {
			var scanned ScannedNode
			scanned.ID = from
			node, err := getNode(from)
			if err != nil {
				scanned.Err = err
			} else {
				scanned.Node = *node
			}

			select {
			case <-ctx.Done():
				return
			case ch <- scanned:
			}
		}

	}(from, to)

	return ch, nil
}

func (s *Substrate) getNode(cl Conn, key types.StorageKey) (*Node, error) {
	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}

	if len(*raw) == 0 {
		return nil, errors.Wrap(ErrNotFound, "node not found")
	}

	var node Node
	if err := types.Decode(*raw, &node); err != nil {
		return nil, errors.Wrap(err, "failed to load object")
	}

	return &node, nil
}

// CreateNode creates a node, this ignores public_config since
// this is only setable by the farmer
func (s *Substrate) CreateNode(identity Identity, node Node) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	if node.TwinID == 0 {
		return 0, fmt.Errorf("twin id is required")
	}

	c, err := types.NewCall(meta, "TfgridModule.create_node",
		node.FarmID,
		node.Resources.TotalResources,
		node.Location,
		node.Interfaces,
		node.SecureBoot,
		node.Virtualized,
		node.BoardSerial,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	_, err = s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create node")
	}

	return s.GetNodeByTwinID(uint32(node.TwinID))

}

// UpdateNode updates a node, this ignores public_config and only keep the value
// set by the farmer
func (s *Substrate) UpdateNode(identity Identity, node Node) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	if node.ID == 0 {
		return 0, fmt.Errorf("node id is required")
	}
	if node.TwinID == 0 {
		return 0, fmt.Errorf("twin id is required")
	}

	c, err := types.NewCall(meta, "TfgridModule.update_node",
		node.ID,
		node.FarmID,
		node.Resources.TotalResources,
		node.Location,
		node.Interfaces,
		node.SecureBoot,
		node.Virtualized,
		node.BoardSerial,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	callResponse, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update node")
	} else {
		log.Debug().Str("hash", callResponse.Hash.Hex()).Msg("update call hash")
	}

	return s.GetNodeByTwinID(uint32(node.TwinID))
}

// UpdateNodeUptime updates the node uptime to given value
func (s *Substrate) UpdateNodeUptime(identity Identity, uptime uint64) (hash types.Hash, err error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return hash, err
	}

	c, err := types.NewCall(meta, "TfgridModule.report_uptime", uptime)

	if err != nil {
		return hash, errors.Wrap(err, "failed to create call")
	}

	callResponse, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return callResponse.Hash, errors.Wrap(err, "failed to update node uptime")
	}

	return callResponse.Hash, nil
}

// SetNodeCertificate sets the node certificate type
func (s *Substrate) SetNodeCertificate(sudo Identity, id uint32, cert NodeCertification) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TfgridModule.set_node_certification",
		id, cert,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	su, err := types.NewCall(meta, "Sudo.sudo", c)
	if err != nil {
		return errors.Wrap(err, "failed to create sudo call")
	}

	if _, err := s.Call(cl, meta, sudo, su); err != nil {
		return errors.Wrap(err, "failed to set node certificate")
	}

	return nil
}

func (s *Substrate) SetNodePowerState(identity Identity, state PowerState) (hash types.Hash, err error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return hash, err
	}

	c, err := types.NewCall(meta, "TfgridModule.change_power_state", state)

	if err != nil {
		return hash, errors.Wrap(err, "failed to create call")
	}

	response, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return hash, errors.Wrap(err, "failed to set node power state")
	}

	return response.Hash, nil
}

// GetNode with id
func (s *Substrate) GetLastNodeID() (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	key, err := types.CreateStorageKey(meta, "TfgridModule", "NodeID")
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}

	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup node id")
	}

	if len(*raw) == 0 {
		return 0, errors.Wrap(ErrNotFound, "no value for last nodeid")
	}

	var v types.U32
	if err := types.Decode(*raw, &v); err != nil {
		return 0, err
	}

	return uint32(v), nil
}
