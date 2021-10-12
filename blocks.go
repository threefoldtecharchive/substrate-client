package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

func (s *Substrate) GetCurrentHeight() (uint32, error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return 0, err
	}

	var blockNumber uint32
	key, err := types.CreateStorageKey(meta, "System", "Number", nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return 0, err
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &blockNumber)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("block number not found")
	}

	return blockNumber, nil
}

func (s *Substrate) FetchEventsForBlockRange(start uint32, end uint32) (types.StorageKey, []types.StorageChangeSet, error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return nil, nil, err
	}

	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return key, nil, err
	}

	lbh, err := cl.RPC.Chain.GetBlockHash(uint64(start))
	if err != nil {
		return key, nil, err
	}

	uph, err := cl.RPC.Chain.GetBlockHash(uint64(end))
	if err != nil {
		return key, nil, err
	}

	rawSet, err := cl.RPC.State.QueryStorage([]types.StorageKey{key}, lbh, uph)
	if err != nil {
		return key, nil, err
	}

	return key, rawSet, nil
}

func (s *Substrate) GetBlock(block types.Hash) (*types.SignedBlock, error) {
	cl, _, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	return cl.RPC.Chain.GetBlock(block)
}
