package substrate

import "github.com/centrifuge/go-substrate-rpc-client/v3/types"

type BurnTransactionCreated struct {
	Phase  types.Phase
	Target AccountID
	// TODO check if this works ....
	Balance     types.U128
	BlockNumber types.U32
	Message     string
	Topics      []types.Hash
}
