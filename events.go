package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

// TODO: add all events from SmartContractModule and TfgridModule

// ContractCreated is the contract created event
type ContractCreated struct {
	Phase    types.Phase
	Contract Contract
	Topics   []types.Hash
}

// ContractUpdated is the contract updated event
type ContractUpdated struct {
	Phase    types.Phase
	Contract Contract
	Topics   []types.Hash
}

// ContractCanceled is the contract canceled event
type ContractCanceled struct {
	Phase      types.Phase
	ContractID types.U64
	Topics     []types.Hash
}

type NodePublicConfig struct {
	Phase  types.Phase
	Node   types.U32
	Config PublicConfig
	Topics []types.Hash
}

// EventRecords is a struct that extends the default events with our events
type EventRecords struct {
	types.EventRecords
	SmartContractModule_ContractCreated  []ContractCreated  //nolint:stylecheck,golint
	SmartContractModule_ContractUpdated  []ContractUpdated  //nolint:stylecheck,golint
	SmartContractModule_ContractCanceled []ContractCanceled //nolint:stylecheck,golint

	// TfgridModule events
	TfgridModule_NodePublicConfigStored []NodePublicConfig
}
