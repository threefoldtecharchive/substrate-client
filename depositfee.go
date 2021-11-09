package substrate

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

var ErrDepositFeeNotFound = fmt.Errorf("deposit fee not found")

func (s *Substrate) GetDepositFee(identity Identity) (int64, error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return 0, err
	}

	var fee types.U64
	key, err := types.CreateStorageKey(meta, "TFTBridgeModule", "DepositFee", nil, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return 0, err
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &fee)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, ErrDepositFeeNotFound
	}

	return int64(fee), nil
}
