package substrate

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// EntityProof struct
type EntityProof struct {
	EntityID  types.U32
	Signature string
}

// Twin struct
type Twin struct {
	ID       types.U32
	Account  AccountID
	Relay    OptionRelay
	Entities []EntityProof
	Pk       types.OptionBytes
}

// OptionRelay type
type OptionRelay struct {
	HasValue bool
	AsValue  string
}

// Encode implementation
func (m OptionRelay) Encode(encoder scale.Encoder) (err error) {
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

// Decode implementation for the enum type
func (r *OptionRelay) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.HasValue = false
		r.AsValue = ""
	case 1:
		r.HasValue = true
		return decoder.Decode(&r.AsValue)
	default:
		return fmt.Errorf("unknown deleted state value")
	}

	return nil
}

// GetTwinByPubKey gets a twin with public key
func (s *Substrate) GetTwinByPubKey(pk []byte) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	key, err := types.CreateStorageKey(meta, "TfgridModule", "TwinIdByAccountID", pk, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}
	var id types.U32
	ok, err := cl.RPC.State.GetStorageLatest(key, &id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup entity")
	}

	if !ok || id == 0 {
		return 0, errors.Wrap(ErrNotFound, "twin not found")
	}

	return uint32(id), nil
}

// GetTwin gets a twin
func (s *Substrate) GetTwin(id uint32) (*Twin, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}
	key, err := types.CreateStorageKey(meta, "TfgridModule", "Twins", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}

	if len(*raw) == 0 {
		return nil, errors.Wrap(ErrNotFound, "twin not found")
	}

	version, err := s.getVersion(*raw)
	if err != nil {
		return nil, err
	}

	var twin Twin

	switch version {
	case 1:
		if err := types.Decode(*raw, &twin); err != nil {
			return nil, errors.Wrap(err, "failed to load object")
		}
	default:
		return nil, ErrUnknownVersion
	}

	return &twin, nil
}

// CreateTwin creates a twin
func (s *Substrate) CreateTwin(identity Identity, relay *string, pk []byte) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	relayOption := OptionRelay{}
	pkOption := types.NewOptionBytesEmpty()
	if relay != nil {
		relayOption = OptionRelay{HasValue: true, AsValue: *relay}
	}
	if pk != nil {
		pkOption = types.NewOptionBytes(pk)
	}

	c, err := types.NewCall(meta, "TfgridModule.create_twin", relayOption, pkOption)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	if _, err := s.Call(cl, meta, identity, c); err != nil {
		return 0, errors.Wrap(err, "failed to create twin")
	}

	return s.GetTwinByPubKey(identity.PublicKey())
}

// UpdateTwin updates a twin
func (s *Substrate) UpdateTwin(identity Identity, relay *string, pk *string) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	relay_bytes := types.OptionBytes{}
	if relay != nil {
		relay_bytes = types.NewOptionBytes([]byte(*relay))
	}

	pk_bytes := types.OptionBytes{}
	if pk != nil {
		pk_bytes = types.NewOptionBytes([]byte(*pk))
	}

	c, err := types.NewCall(meta, "TfgridModule.update_twin", relay_bytes, pk_bytes)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	if _, err := s.Call(cl, meta, identity, c); err != nil {
		return 0, errors.Wrap(err, "failed to update twin")
	}

	return s.GetTwinByPubKey(identity.PublicKey())
}
