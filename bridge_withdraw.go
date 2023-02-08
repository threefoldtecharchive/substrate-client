package substrate

import (
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

var (
	ErrWithdrawTransactionNotFound = fmt.Errorf("withdraw tx not found")
	ErrRefundTransactionNotFound   = fmt.Errorf("refund tx not found")
	ErrFailedToDecode              = fmt.Errorf("failed to decode events, skipping")
)

type WithdrawTransaction struct {
	ID             types.U64
	Block          types.U32
	Amount         types.U64
	Target         string
	Signatures     []StellarSignature
	SequenceNumber types.U64
}

func (s *Substrate) ProposeWithdrawTransactionOrAddSig(identity Identity, txID uint64, target string, amount *big.Int, signature string, stellarAddress string, sequence_number uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TFTBridgeModule.propose_withdraw_transaction_or_add_sig",
		txID, target, types.U64(amount.Uint64()), signature, stellarAddress, sequence_number,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	_, err = s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to propose withdraw transaction")
	}

	return nil
}

func (s *Substrate) SetWithdrawTransactionExecuted(identity Identity, txID uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TFTBridgeModule.set_withdraw_transaction_executed", txID)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	_, err = s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to set withdraw transaction executed")
	}

	return nil
}

func (s *Substrate) GetWithdrawTransaction(withdrawTransactionID types.U64) (*WithdrawTransaction, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(withdrawTransactionID)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	var withdrawTx WithdrawTransaction
	key, err := types.CreateStorageKey(meta, "TFTBridgeModule", "WithdrawTransactions", bytes, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return nil, err
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &withdrawTx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrWithdrawTransactionNotFound
	}

	return &withdrawTx, nil
}

func (s *Substrate) IsAlreadyWithdrawn(withdrawTransactionID types.U64) (exists bool, err error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return false, err
	}

	bytes, err := types.Encode(withdrawTransactionID)
	if err != nil {
		return false, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	var withdrawTx WithdrawTransaction
	key, err := types.CreateStorageKey(meta, "TFTBridgeModule", "ExecutedWithdrawTransactions", bytes, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &withdrawTx)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}

func (s *Substrate) GetPendingWithdraws() (*[]WithdrawTransaction, error) {
	cl, _, err := s.getClient()
	if err != nil {
		return nil, err
	}

	skey := createPrefixedKey("TFTBridgeModule", "BurnTransactions")

	keys, err := cl.RPC.State.GetKeysLatest(skey)
	if err != nil {
		return nil, err
	}

	var withdrawTxs []WithdrawTransaction
	for _, k := range keys {
		var withdrawTx WithdrawTransaction

		ok, err := cl.RPC.State.GetStorageLatest(k, &withdrawTx)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, ErrWithdrawTransactionNotFound
		}
		withdrawTxs = append(withdrawTxs, withdrawTx)

	}

	return &withdrawTxs, nil
}
