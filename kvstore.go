package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

func (s *Substrate) KVStoreSet(key string, value string, identity Identity) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TFKVStore.set",
		key, value,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to create contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

func (s *Substrate) KVSToreDelete(key string, identity Identity) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TFKVStore.delete",
		key,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to create contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}
	return nil
}

func (s *Substrate) KVStoreGet(key string, Identity Identity) ([]byte, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}
	bytes, err := types.Encode(key)
	if err != nil {
		return nil, err
	}
	storageKey, err := types.CreateStorageKey(meta, "TFKVStore", "TFKVStore", Identity.PublicKey(), bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}
	var value []byte
	ok, err := cl.RPC.State.GetStorageLatest(storageKey, &value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}
	if !ok {
		return nil, errors.Wrap(ErrNotFound, "key not found")
	}
	return value, nil
}

func (s *Substrate) KVStoreList(identity Identity) ([]byte, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}
	bytes, err := types.Encode("")
	if err != nil {
		return nil, err
	}

	storageKey, err := types.CreateStorageKey(meta, "TFKVStore", "TFKVStore.entries", identity.PublicKey(), bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}
	var value []byte

	ok, err := cl.RPC.State.GetStorageLatest(storageKey, &value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}
	if !ok {
		return nil, errors.Wrap(ErrNotFound, "key not found")
	}
	return value, nil
}
