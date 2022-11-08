package substrate

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

type DeletedState struct {
	IsCanceledByUser bool
	IsOutOfFunds     bool
}

// Decode implementation for the enum type
func (r *DeletedState) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsCanceledByUser = true
	case 1:
		r.IsOutOfFunds = true
	case 2:
	default:
		return fmt.Errorf("unknown deleted state value")
	}

	return nil
}

// Encode implementation
func (r DeletedState) Encode(encoder scale.Encoder) (err error) {
	if r.IsCanceledByUser {
		err = encoder.PushByte(0)
	} else if r.IsOutOfFunds {
		err = encoder.PushByte(1)
	}
	return
}

// ContractState enum
type ContractState struct {
	IsCreated                bool
	IsDeleted                bool
	AsDeleted                DeletedState
	IsGracePeriod            bool
	AsGracePeriodBlockNumber types.U64
}

// Decode implementation for the enum type
func (r *ContractState) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsCreated = true
	case 1:
		r.IsDeleted = true
		if err := decoder.Decode(&r.AsDeleted); err != nil {
			return errors.Wrap(err, "failed to get deleted state")
		}
	case 2:
		r.IsGracePeriod = true
		if err := decoder.Decode(&r.AsGracePeriodBlockNumber); err != nil {
			return errors.Wrap(err, "failed to get grace period")
		}
	default:
		return fmt.Errorf("unknown ContractState value")
	}

	return nil
}

// Encode implementation
func (r ContractState) Encode(encoder scale.Encoder) (err error) {
	if r.IsCreated {
		err = encoder.PushByte(0)
	} else if r.IsDeleted {
		if err = encoder.PushByte(1); err != nil {
			return err
		}
		err = encoder.Encode(r.AsDeleted)
	} else if r.IsGracePeriod {
		if err = encoder.PushByte(2); err != nil {
			return err
		}
		err = encoder.Encode(r.AsGracePeriodBlockNumber)
	}

	return
}

type HexHash [32]byte

func (h HexHash) String() string {
	return string(h[:])
}

// NewHexHash will create a new hash from a hex input (32 bytes)
func NewHexHash(hash string) (hexHash HexHash) {
	copy(hexHash[:], hash)
	return
}

type ConsumableResources struct {
	TotalResources Resources
	UsedResources  Resources
}

type CapacityReservationContract struct {
	NodeID              types.U32
	Resources           ConsumableResources
	GroupID             types.OptionU32
	PublicIPs           types.U32
	DeploymentContracts []types.U64
}

type DeploymentContract struct {
	CapacityReservationID types.U64
	DeploymentHash        HexHash
	DeploymentData        string
	PublicIPsCount        types.U32
	PublicIPs             []PublicIP
	Resources             Resources
}

type NameContract struct {
	Name string
}

type RentContract struct {
	Node types.U32
}

type ContractType struct {
	IsDeploymentContract bool
	DeploymentContract   DeploymentContract
	IsNameContract       bool
	NameContract         NameContract
	IsRentContract       bool
	RentContract         RentContract
}

type Group struct {
	Id                             types.U32
	TwinID                         types.U32
	CapacityReservationContractIDs []types.U32
}

type NodeGroupConfig struct {
	Id      types.U32
	GroudID types.U32
}

type CapacityReservationPolicy struct {
	IsAny       bool
	AsAny       Any
	IsExclusive bool
	AsExclusive Exclusive
	IsNode      bool
	AsNode      NodePolicy
}

type Any struct {
	Resources Resources
	Features  OptionFeatures
}

type Exclusive struct {
	GroupID   types.U32
	Resources Resources
	Features  OptionFeatures
}

type OptionFeatures struct {
	HasValue bool
	AsValue  []NodeFeatures
}

type NodePolicy struct {
	NodeID types.U32
}

// Encode implementation
func (m OptionFeatures) Encode(encoder scale.Encoder) (err error) {
	var i byte
	if m.HasValue {
		i = 1
	}
	err = encoder.PushByte(i)
	if err != nil {
		return err
	}

	if m.HasValue {
		err = encoder.Encode(m.AsValue)
	}

	return
}

// Decode implementation
func (m *OptionFeatures) Decode(decoder scale.Decoder) (err error) {
	var i byte
	if err := decoder.Decode(&i); err != nil {
		return err
	}

	switch i {
	case 0:
		return nil
	case 1:
		m.HasValue = true
		return decoder.Decode(&m.AsValue)
	default:
		return fmt.Errorf("unknown value for Option")
	}
}

// Decode implementation for the enum type
func (r *ContractType) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsDeploymentContract = true
		if err := decoder.Decode(&r.DeploymentContract); err != nil {
			return err
		}
	case 1:
		r.IsNameContract = true
		if err := decoder.Decode(&r.NameContract); err != nil {
			return err
		}
	case 2:
		r.IsRentContract = true
		if err := decoder.Decode(&r.RentContract); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown contract type value")
	}

	return nil
}

// Encode implementation
func (r ContractType) Encode(encoder scale.Encoder) (err error) {
	if r.IsDeploymentContract {
		if err = encoder.PushByte(0); err != nil {
			return err
		}
		if err = encoder.Encode(r.DeploymentContract); err != nil {
			return err
		}
	} else if r.IsNameContract {
		if err = encoder.PushByte(1); err != nil {
			return err
		}

		if err = encoder.Encode(r.NameContract); err != nil {
			return err
		}
	} else if r.IsRentContract {
		if err = encoder.PushByte(2); err != nil {
			return err
		}

		if err = encoder.Encode(r.RentContract); err != nil {
			return err
		}
	}

	return
}

// Decode implementation for the enum type
func (r *CapacityReservationPolicy) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsAny = true
		if err := decoder.Decode(&r.AsAny); err != nil {
			return err
		}
	case 1:
		r.IsExclusive = true
		if err := decoder.Decode(&r.AsExclusive); err != nil {
			return err
		}
	case 2:
		r.IsNode = true
		if err := decoder.Decode(&r.AsNode); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown capacity reservation policy value")
	}

	return nil
}

// Encode implementation
func (r CapacityReservationPolicy) Encode(encoder scale.Encoder) (err error) {
	if r.IsAny {
		if err = encoder.PushByte(0); err != nil {
			return err
		}
		if err = encoder.Encode(r.AsAny); err != nil {
			return err
		}
	} else if r.IsExclusive {
		if err = encoder.PushByte(1); err != nil {
			return err
		}

		if err = encoder.Encode(r.AsExclusive); err != nil {
			return err
		}
	} else if r.IsNode {
		if err = encoder.PushByte(2); err != nil {
			return err
		}

		if err = encoder.Encode(r.AsNode); err != nil {
			return err
		}
	}

	return
}

// Contract structure
type Contract struct {
	Versioned
	State              ContractState
	ContractID         types.U64
	TwinID             types.U32
	ContractType       ContractType
	SolutionProviderID types.OptionU64
}

// CreateDeploymentContract creates a contract for deployment
func (s *Substrate) CreateDeploymentContract(identity Identity, node uint32, body string, hash string, publicIPs uint32, solutionProviderID *uint64) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	var providerID types.OptionU64
	if solutionProviderID != nil {
		providerID = types.NewOptionU64(types.U64(*solutionProviderID))
	}

	h := NewHexHash(hash)
	c, err := types.NewCall(meta, "SmartContractModule.create_node_contract",
		node, h, body, publicIPs, providerID,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	return s.GetContractWithHash(node, h)
}

// CreateNameContract creates a contract for deployment
func (s *Substrate) CreateNameContract(identity Identity, name string) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	c, err := types.NewCall(meta, "SmartContractModule.create_name_contract",
		name,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	return s.GetContractIDByNameRegistration(name)
}

// CreateRentContract creates a rent contract on a node
func (s *Substrate) CreateRentContract(identity Identity, node uint32, solutionProviderID *uint64) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	var providerID types.OptionU64
	if solutionProviderID != nil {
		providerID = types.NewOptionU64(types.U64(*solutionProviderID))
	}

	c, err := types.NewCall(meta, "SmartContractModule.create_rent_contract",
		node, providerID,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create rent contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	return s.GetNodeRentContract(node)
}

// UpdateDeploymentContract updates existing contract
func (s *Substrate) UpdateDeploymentContract(identity Identity, contract uint64, body string, hash string) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	h := NewHexHash(hash)
	c, err := types.NewCall(meta, "SmartContractModule.update_node_contract",
		contract, h, body,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	return contract, nil
}

// CancelContract creates a contract for deployment
func (s *Substrate) CancelContract(identity Identity, contract uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.cancel_contract", contract)

	if err != nil {
		return errors.Wrap(err, "failed to cancel call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to cancel contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// SetContractConsumption can only be called by the node that owns the contract to set the used
// resources associated with the node.
func (s *Substrate) SetContractConsumption(identity Identity, resources ...ContractResources) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.report_contract_resources",
		resources,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to set contract used resources")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// GetContract we should not have calls to create contract, instead only get
func (s *Substrate) GetContract(id uint64) (*Contract, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "Contracts", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	return s.getContract(cl, key)
}

// GetContractWithHash gets a contract given the node id and hash
func (s *Substrate) GetContractWithHash(node uint32, hash HexHash) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	nodeBytes, err := types.Encode(node)
	if err != nil {
		return 0, err
	}
	hashBytes, err := types.Encode(hash)
	if err != nil {
		return 0, err
	}
	key, err := types.CreateStorageKey(meta, "SmartContractModule", "ContractIDByNodeIDAndHash", nodeBytes, hashBytes, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}
	var contract types.U64
	_, err = cl.RPC.State.GetStorageLatest(key, &contract)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup contracts")
	}

	if contract == 0 {
		return 0, errors.Wrap(ErrNotFound, "contract not found")
	}

	return uint64(contract), nil
}

// GetContractIDByNameRegistration gets a contract given the its name
func (s *Substrate) GetContractIDByNameRegistration(name string) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	nameBytes, err := types.Encode(name)
	if err != nil {
		return 0, err
	}
	key, err := types.CreateStorageKey(meta, "SmartContractModule", "ContractIDByNameRegistration", nameBytes, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}
	var contract types.U64
	_, err = cl.RPC.State.GetStorageLatest(key, &contract)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup contracts")
	}

	if contract == 0 {
		return 0, errors.Wrap(ErrNotFound, "contract not found")
	}

	return uint64(contract), nil
}

// GetDeploymentContracts gets all contracts on a node (pk) in given state
func (s *Substrate) GetDeploymentContracts(node uint32) ([]types.U64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	nodeBytes, err := types.Encode(node)
	if err != nil {
		return nil, err
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "ActiveDeploymentContracts", nodeBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}
	var contracts []types.U64
	_, err = cl.RPC.State.GetStorageLatest(key, &contracts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup contracts")
	}

	return contracts, nil
}

// GetNodeContracts gets all contracts on a node (pk) in given state
func (s *Substrate) GetNodeRentContract(node uint32) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	nodeBytes, err := types.Encode(node)
	if err != nil {
		return 0, err
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "ActiveRentContractForNode", nodeBytes)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create substrate query key")
	}

	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return 0, errors.Wrap(err, "failed to lookup contract")
	}

	if len(*raw) == 0 {
		return 0, errors.Wrap(ErrNotFound, "contract not found")
	}

	var contract uint64
	err = types.Decode(*raw, &contract)
	return contract, err
}

func (s *Substrate) getContract(cl Conn, key types.StorageKey) (*Contract, error) {
	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup contract")
	}

	if len(*raw) == 0 {
		return nil, errors.Wrap(ErrNotFound, "contract not found")
	}

	var contract Contract
	if err := types.Decode(*raw, &contract); err != nil {
		return nil, errors.Wrap(err, "failed to load object")
	}

	return &contract, nil
}

// Consumption structure
type NruConsumption struct {
	ContractID types.U64
	Timestamp  types.U64
	Window     types.U64
	NRU        types.U64
}

// Consumption structure
type Consumption struct {
	ContractID types.U64
	Timestamp  types.U64
	CRU        types.U64 `json:"cru"`
	SRU        types.U64 `json:"sru"`
	HRU        types.U64 `json:"hru"`
	MRU        types.U64 `json:"mru"`
	NRU        types.U64 `json:"nru"`
}

// IsEmpty true if consumption is zero
func (s *NruConsumption) IsEmpty() bool {
	return s.NRU == 0
}

// Report send a capacity report to substrate
func (s *Substrate) Report(identity Identity, consumptions []NruConsumption) (hash types.Hash, err error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return hash, err
	}

	c, err := types.NewCall(meta, "SmartContractModule.add_nru_reports", consumptions)
	if err != nil {
		return hash, errors.Wrap(err, "failed to create call")
	}

	hash, err = s.Call(cl, meta, identity, c)
	if err != nil {
		return hash, errors.Wrap(err, "failed to create report")
	}

	return hash, nil
}

type ContractResources struct {
	ContractID types.U64
	Used       Resources
}
