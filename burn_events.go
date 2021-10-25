package substrate

import "github.com/centrifuge/go-substrate-rpc-client/v3/types"

type BurnTransactionCreated struct {
	Phase  types.Phase
	Target AccountID
	// TODO check if this works ....
	Balance     types.BalanceStatus
	BlockNumber types.U32
	Message     string
	Topics      []types.Hash
}
