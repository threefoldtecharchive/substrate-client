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

type FarmStored struct {
	Phase  types.Phase
	Farm   Farm
	Topics []types.Hash
}

type FarmDeleted struct {
	Phase  types.Phase
	Farm   types.U32
	Topics []types.Hash
}

type NodeStored struct {
	Phase  types.Phase
	Node   Node
	Topics []types.Hash
}

type NodeDeleted struct {
	Phase  types.Phase
	Node   types.U32
	Topics []types.Hash
}
type NodeUptimeReported struct {
	Phase     types.Phase
	Node      types.U32
	Timestamp types.U64
	Uptime    types.U64
	Topics    []types.Hash
}

type EntityStored struct {
	Phase  types.Phase
	Entity Entity
	Topics []types.Hash
}

type EntityDeleted struct {
	Phase  types.Phase
	Entity types.U32
	Topics []types.Hash
}

type TwinStored struct {
	Phase  types.Phase
	Twin   Twin
	Topics []types.Hash
}

type TwinDeleted struct {
	Phase  types.Phase
	Twin   types.U32
	Topics []types.Hash
}

type TwinEntityStored struct {
	Phase     types.Phase
	Twin      types.U32
	Entity    types.U32
	Signature []byte
	Topics    []types.Hash
}

type TwinEntityRemoved struct {
	Phase  types.Phase
	Twin   types.U32
	Entity types.U32
	Topics []types.Hash
}

// numeric enum for unit
type Unit byte

type Policy struct {
	Value types.U32
	Unit  Unit
}
type PricingPolicy struct {
	Versioned
	ID                    types.U32
	Name                  string
	SU                    Policy
	CU                    Policy
	NU                    Policy
	IPU                   Policy
	UniqueName            Policy
	DomainName            Policy
	FoundationAccount     AccountID
	CertifiedSalesAccount AccountID
}

type PricingPolicyStored struct {
	Phase  types.Phase
	Policy PricingPolicy
	Topics []types.Hash
}

type FarmingPolicy struct {
	Versioned
	ID                types.U32
	Name              string
	CU                types.U32
	SU                types.U32
	NU                types.U32
	IPv4              types.U32
	Timestamp         types.U64
	CertificationType CertificationType
}

type FarmingPolicyStored struct {
	Phase  types.Phase
	Policy FarmingPolicy
	Topics []types.Hash
}

type CertificationCodes struct {
	Versioned
	ID                    types.U32
	Name                  string
	Description           string
	CertificationCodeType byte
}

type CertificationCodeStored struct {
	Phase  types.Phase
	Codes  CertificationCodes
	Topics []types.Hash
}

type FarmPayoutV2AddressRegistered struct {
	Phase   types.Phase
	Farm    types.U32
	Address string
	Topics  []types.Hash
}

// EventRecords is a struct that extends the default events with our events
type EventRecords struct {
	types.EventRecords
	SmartContractModule_ContractCreated  []ContractCreated  //nolint:stylecheck,golint
	SmartContractModule_ContractUpdated  []ContractUpdated  //nolint:stylecheck,golint
	SmartContractModule_ContractCanceled []ContractCanceled //nolint:stylecheck,golint

	// farm events
	TfgridModule_FarmStored  []FarmStored
	TfgridModule_FarmUpdated []FarmStored
	TfgridModule_FarmDeleted []FarmDeleted

	// node events
	TfgridModule_NodeStored             []NodeStored
	TfgridModule_NodeUpdated            []NodeStored
	TfgridModule_NodeDeleted            []NodeDeleted
	TfgridModule_NodeUptimeReported     []NodeUptimeReported
	TfgridModule_NodePublicConfigStored []NodePublicConfig

	// entity events
	TfgridModule_EntityStored  []EntityStored
	TfgridModule_EntityUpdated []EntityStored
	TfgridModule_EntityDeleted []EntityDeleted

	// twin events
	TfgridModule_TwinStored        []TwinStored
	TfgridModule_TwinUpdated       []TwinStored
	TfgridModule_TwinDeleted       []TwinDeleted
	TfgridModule_TwinEntityStored  []TwinEntityStored
	TfgridModule_TwinEntityRemoved []TwinEntityRemoved

	// policy events
	TfgridModule_PricingPolicyStored []PricingPolicyStored
	TfgridModule_FarmingPolicyStored []FarmingPolicyStored

	// other events
	TfgridModule_CertificationCodeStored       []CertificationCodeStored
	TfgridModule_FarmPayoutV2AddressRegistered []FarmPayoutV2AddressRegistered
}
