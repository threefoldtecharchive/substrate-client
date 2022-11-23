package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// ServiceContract struct
type ServiceContract struct {
	ServiceTwinID      types.U32
	ConsumerTwinID     types.U32
	BaseFee            types.U64
	VariableFee        types.U64
	Metadata           string
	AcceptedByService  bool
	AcceptedByConsumer bool
	LastBill           types.U64
	State              ServiceContractState
}

// ServiceContractBill struct
type ServiceContractBill struct {
	VariableAmount types.U64
	Window         types.U64
	Metadata       string
}

// ServiceContractState enum
type ServiceContractState struct {
	IsCreated        bool
	IsAgreementReady bool
	IsApprovedByBoth bool
}

// OK - service_contract_create()
// OK - service_contract_set_metadata()
// OK - service_contract_set_fees()
// OK - service_contract_approve()
// OK - service_contract_reject()
// OK - service_contract_cancel()
// OK - service_contract_bill()

// ServiceContractCreate creates a service contract
func (s *Substrate) ServiceContractCreate(identity Identity, service AccountID, consumer AccountID) (uint64, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return 0, err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_create",
		service, consumer,
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return 0, err
	}

	var contractID uint64 = 0
	return contractID, nil
}

// ServiceContractSetMetadata sets metadata for a service contract
func (s *Substrate) ServiceContractSetMetadata(identity Identity, contract uint64, metadata string) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_set_metadata",
		contract, metadata,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to set metadata for service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// ServiceContractSetFees sets fees for a service contract
func (s *Substrate) ServiceContractSetFees(identity Identity, contract uint64, base_fee uint64, variable_fee uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_set_fees",
		contract, base_fee, variable_fee,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to set fees for service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// ServiceContractApprove approves a service contract
func (s *Substrate) ServiceContractApprove(identity Identity, contract uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_approve",
		contract,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to approve service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// ServiceContractReject rejects a service contract
func (s *Substrate) ServiceContractReject(identity Identity, contract uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_reject",
		contract,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to reject service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// ServiceContractCancel cancels a service contract
func (s *Substrate) ServiceContractCancel(identity Identity, contract uint64) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_cancel",
		contract,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to cancel service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// ServiceContractBill bills a service contract
func (s *Substrate) ServiceContractBill(identity Identity, contract uint64, variable_amount uint64, metadata string) error {
	cl, meta, err := s.getClient()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "SmartContractModule.service_contract_bill",
		contract, variable_amount, metadata,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to bill service contract")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}

// GetServiceContract gets a service contract given the service contract id
func (s *Substrate) GetServiceContract(id uint64) (*ServiceContract, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return nil, err
	}

	bytes, err := types.Encode(id)
	if err != nil {
		return nil, errors.Wrap(err, "substrate: encoding error building query arguments")
	}

	key, err := types.CreateStorageKey(meta, "SmartContractModule", "ServiceContracts", bytes, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create substrate query key")
	}

	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup contract")
	}

	if len(*raw) == 0 {
		return nil, errors.Wrap(ErrNotFound, "service contract not found")
	}

	var contract ServiceContract
	if err := types.Decode(*raw, &contract); err != nil {
		return nil, errors.Wrap(err, "failed to load object")
	}

	return &contract, nil
}
