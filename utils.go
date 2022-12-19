package substrate

import (
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vedhavyas/go-subkey"
	"golang.org/x/crypto/blake2b"
)

var (
	ErrIsUsurped = fmt.Errorf("Is Usurped")
)

// https://github.com/threefoldtech/tfchain_pallets/blob/bc9c5d322463aaf735212e428da4ea32b117dc24/pallet-smart-contract/src/lib.rs#L58
var smartContractModuleErrors = []string{
	"TwinNotExists",
	"NodeNotExists",
	"FarmNotExists",
	"FarmHasNotEnoughPublicIPs",
	"FarmHasNotEnoughPublicIPsFree",
	"FailedToReserveIP",
	"FailedToFreeIPs",
	"ContractNotExists",
	"TwinNotAuthorizedToUpdateContract",
	"TwinNotAuthorizedToCancelContract",
	"NodeNotAuthorizedToDeployContract",
	"NodeNotAuthorizedToComputeReport",
	"PricingPolicyNotExists",
	"ContractIsNotUnique",
	"NameExists",
	"NameNotValid",
	"InvalidContractType",
	"TFTPriceValueError",
	"NotEnoughResourcesOnNode",
	"NodeNotAuthorizedToReportResources",
	"MethodIsDeprecated",
	"NodeHasActiveContracts",
	"NodeHasRentContract",
	"NodeIsNotDedicated",
	"NodeNotAvailableToDeploy",
}

// Sign signs data with the private key under the given derivation path, returning the signature. Requires the subkey
// command to be in path
func signBytes(data []byte, privateKeyURI string, scheme subkey.Scheme) ([]byte, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	kyr, err := subkey.DeriveKeyPair(scheme, privateKeyURI)
	if err != nil {
		return nil, err
	}

	signature, err := kyr.Sign(data)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Sign adds a signature to the extrinsic
func (s *Substrate) sign(e *types.Extrinsic, signer Identity, o types.SignatureOptions) error {
	if e.Type() != types.ExtrinsicVersion4 {
		return fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(), e.Type())
	}

	mb, err := types.Encode(e.Method)
	if err != nil {
		return err
	}

	era := o.Era
	if !o.Era.IsMortalEra {
		era = types.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := types.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: o.TransactionVersion,
	}

	signerPubKey := types.NewMultiAddressFromAccountID(signer.PublicKey())

	b, err := types.Encode(payload)
	if err != nil {
		return err
	}

	sig, err := signer.Sign(b)

	if err != nil {
		return err
	}
	msig := signer.MultiSignature(sig)
	extSig := types.ExtrinsicSignatureV4{
		Signer:    signerPubKey,
		Signature: msig,
		Era:       era,
		Nonce:     o.Nonce,
		Tip:       o.Tip,
	}

	e.Signature = extSig

	// mark the extrinsic as signed
	e.Version |= types.ExtrinsicBitSigned

	return nil
}

// Call call this extrinsic and retry if Usurped
func (s *Substrate) Call(cl Conn, meta Meta, identity Identity, call types.Call) (hash types.Hash, err error) {
	for {
		hash, err := s.CallOnce(cl, meta, identity, call)

		if errors.Is(err, ErrIsUsurped) {
			continue
		}

		return hash, err
	}
}

func (s *Substrate) CallOnce(cl Conn, meta Meta, identity Identity, call types.Call) (hash types.Hash, err error) {
	// Create the extrinsic
	ext := types.NewExtrinsic(call)

	genesisHash, err := cl.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return hash, errors.Wrap(err, "failed to get genesisHash")
	}

	rv, err := cl.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return hash, err
	}

	//node.Address =identity.PublicKey
	account, err := s.getAccount(cl, meta, identity)
	if err != nil {
		return hash, errors.Wrap(err, "failed to get account")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(account.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = s.sign(&ext, identity, o)
	if err != nil {
		return hash, errors.Wrap(err, "failed to sign")
	}

	// Send the extrinsic
	sub, err := cl.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return hash, errors.Wrap(err, "failed to submit extrinsic")
	}

	defer sub.Unsubscribe()

	ch := sub.Chan()
	ech := sub.Err()

loop:
	for {
		select {
		case err := <-ech:
			return hash, errors.Wrap(err, "error failed on extrinsic status")
		case <-time.After(30 * time.Second):
			return hash, fmt.Errorf("extrinsic timeout waiting for block")
		case event := <-ch:
			if event.IsReady || event.IsBroadcast {
				continue
			} else if event.IsInBlock {
				hash = event.AsInBlock
				break loop
			} else if event.IsFinalized {
				// we shouldn't hit this case
				// any more since InBlock will always
				// happen first we leave it only
				// as a safety net
				hash = event.AsFinalized
				break loop
			} else if event.IsDropped || event.IsInvalid {
				return hash, fmt.Errorf("failed to make call")
			} else if event.IsUsurped {
				return hash, ErrIsUsurped
			} else {
				log.Error().Err(err).Msgf("extrinsic block in an unhandled state: %+v", event)
			}
		}
	}

	return hash, nil
}

func (s *Substrate) checkForError(cl Conn, meta Meta, blockHash types.Hash, signer types.AccountID) error {
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return err
	}

	raw, err := cl.RPC.State.GetStorageRaw(key, blockHash)
	if err != nil {
		return err
	}

	block, err := cl.RPC.Chain.GetBlock(blockHash)
	if err != nil {
		return err
	}

	events := EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)
	if err != nil {
		log.Debug().Msgf("Failed to decode event %+v", err)
		return nil
	}

	if len(events.System_ExtrinsicFailed) > 0 {
		for _, e := range events.System_ExtrinsicFailed {
			who := block.Block.Extrinsics[e.Phase.AsApplyExtrinsic].Signature.Signer.AsID
			if signer == who {
				if int(e.DispatchError.ModuleError.Error) >= len(smartContractModuleErrors) {
					return fmt.Errorf("error with code %d occured", e.DispatchError.ModuleError.Error)
				} else {
					return fmt.Errorf(smartContractModuleErrors[e.DispatchError.ModuleError.Error])
				}
			}
		}
	}

	return nil
}
