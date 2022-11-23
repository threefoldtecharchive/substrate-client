package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Group struct {
	Id                             types.U32
	TwinID                         types.U32
	CapacityReservationContractIDs []types.U32
}

type NodeGroupConfig struct {
	Id      types.U32
	GroudID types.U32
}

func (s *Substrate) getGroup(cl Conn, key types.StorageKey) (*Group, error) {
	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup group")
	}

	if len(*raw) == 0 {
		return nil, errors.Wrap(ErrNotFound, "group not found")
	}

	var group Group
	if err := types.Decode(*raw, &group); err != nil {
		return nil, errors.Wrap(err, "failed to load object")
	}

	return &group, nil
}

// CreateGroup
func (s *Substrate) CreateGroup(identity Identity) (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	c, err := types.NewCall(meta, "SmartContractModule.create_group")

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create group")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	id, err := s.GetGroupID()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get group id after creating the group")
	}
	return id, nil
}

// Delete a group by ID
func (s *Substrate) DeleteGroup(identity Identity, groupID uint32) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.delete_group", groupID)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to create group")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// GetGroup
func (s *Substrate) GetGroup(id uint64) (*Group, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "Groups", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	return s.getGroup(cl, key)
}

// GetNode with id
func (s *Substrate) GetGroupID() (uint32, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "GroupID", nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}

	var id types.U32
	ok, err := cl.RPC.State.GetStorageLatest(key, &id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup entity")
	}
	log.Debug().Msgf("%d", id)
	log.Debug().Msgf("%t", ok)
	if !ok || id == 0 {
		return 0, errors.Wrap(ErrNotFound, "group id not found")
	}

	return uint32(id), nil
}
