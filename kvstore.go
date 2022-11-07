package substrate

import (
	"encoding/hex"
	"log"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

func (s *Substrate) KVStoreSet(identity Identity, key string, value string) error {
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

func (s *Substrate) KVStoreDelete(identity Identity, key string) error {
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

func (s *Substrate) KVStoreGet(pk []byte, key string) ([]byte, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}
	bytes, err := types.Encode(key)
	if err != nil {
		return nil, err
	}
	storageKey, err := types.CreateStorageKey(meta, "TFKVStore", "TFKVStore", pk, bytes)
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

type Val struct {
	Id  string
	Key string
}

func (s *Substrate) KVStoreList(identity Identity) (map[string]string, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(types.NewOptionBytes1024Empty())
	if err != nil {
		return nil, err
	}

	log.Printf("bytes: %+v", bytes)
	// id := types.NewAccountID([]byte(identity.Address()))

	storageKey, err := types.CreateStorageKey(meta, "TFKVStore", "TFKVStore", identity.PublicKey(), bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}
	// v, err := hex.DecodeString("bd05d43493846b6a9b68833881cd4092bd05d43493846b6a9b68833881cd409222838b57e06152572297591bf067d778a6ce805400dfd431f12e495d6e14e56388bb1da5144096387d83e11eb8afbd0d")

	// sk := types.NewStorageKey(v)
	log.Printf("sk: %+v", hex.EncodeToString(storageKey))
	// log.Printf("sk: %+v", hex.EncodeToString(sk))
	// log.Printf("storageKey: %+v", storageKey)
	// var value []byte
	keys, err := cl.RPC.State.GetKeysLatest(storageKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup entity")
	}

	query, err := cl.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return nil, err
	}
	for _, q := range query {
		// key := make([][]byte, 0)
		// value := ""
		changes := q.Changes

		for _, c := range changes {
			log.Printf("key: %+v", c.StorageKey.Hex())
			log.Printf("value: %+v", string(c.StorageData))
		}
	}

	return nil, nil
}
