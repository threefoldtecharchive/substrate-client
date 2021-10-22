package substrate

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

// ContractBill structure
type ContractBill struct {
	ContractID    types.U64
	Timestamp     types.U64
	DiscountLevel DiscountLevel
	AmountBilled  types.U128
}

// DiscountLevel enum
type DiscountLevel struct {
	IsNone    bool
	IsDefault bool
	IsBronze  bool
	IsSilver  bool
	IsGold    bool
}

// Decode implementation for the enum type
func (r *DiscountLevel) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		r.IsNone = true
	case 1:
		r.IsDefault = true
	case 2:
		r.IsBronze = true
	case 3:
		r.IsSilver = true
	case 4:
		r.IsGold = true
	default:
		return fmt.Errorf("unknown CertificateType value")
	}

	return nil
}

// Encode implementation
func (r DiscountLevel) Encode(encoder scale.Encoder) (err error) {
	if r.IsNone {
		err = encoder.PushByte(0)
	} else if r.IsDefault {
		err = encoder.PushByte(1)
	} else if r.IsBronze {
		err = encoder.PushByte(2)
	} else if r.IsSilver {
		err = encoder.PushByte(3)
	} else if r.IsGold {
		err = encoder.PushByte(4)
	}

	return
}

// ContractCanceled
type ContractDeployed struct {
	Phase      types.Phase
	ContractID types.U64
	AccountID  AccountID
	Topics     []types.Hash
}

// ContractCanceled
type ConsumptionReportReceived struct {
	Phase       types.Phase
	Consumption Consumption
	Topics      []types.Hash
}

// ContractBilled
type ContractBilled struct {
	Phase        types.Phase
	ContractBill ContractBill
	Topics       []types.Hash
}

// ContractCanceled
type IPsReserved struct {
	Phase      types.Phase
	ContractID types.U64
	IPs        []PublicIP
	Topics     []types.Hash
}

// ContractCanceled
type IPsFreed struct {
	Phase      types.Phase
	ContractID types.U64
	IPs        []string
	Topics     []types.Hash
}
