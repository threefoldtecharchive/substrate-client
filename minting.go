package substrate

import (
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

var (
	ErrMintTransactionNotFound = fmt.Errorf("mint tx not found")
)

type MintTransaction struct {
	Amount types.U64
	Target types.AccountID
	Block  types.U32
	Votes  types.U32
}

func (s *Substrate) IsMintedAlready(identity Identity, mintTxID string) (exists bool, err error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return false, err
	}

	bytes, err := types.EncodeToBytes(mintTxID)
	if err != nil {
		return false, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	var mintTX MintTransaction
	key, err := types.CreateStorageKey(meta, "TFTBridgeModule", "ExecutedMintTransactions", bytes, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &mintTX)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, ErrMintTransactionNotFound
	}

	return true, nil
}

func (s *Substrate) ProposeOrVoteMintTransaction(identity Identity, txID string, target AccountID, amount *big.Int) (*types.Call, error) {
	_, meta, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	c, err := types.NewCall(meta, "TFTBridgeModule.propose_or_vote_mint_transaction",
		txID, target, types.U64(amount.Uint64()),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create call")
	}

	return &c, nil
}
