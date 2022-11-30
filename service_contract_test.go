package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestServiceContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	serviceTwinID := assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)
	consumerTwinID := assertCreateTwin(t, cl, BobMnemonics, BobAddress)

	serviceIdentity, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	consumerIdentity, err := NewIdentityFromSr25519Phrase(BobMnemonics)
	require.NoError(t, err)

	serviceAccount, err := FromAddress(AliceAddress)
	require.NoError(t, err)

	consumerAccount, err := FromAddress(BobAddress)
	require.NoError(t, err)

	serviceContractID, err := cl.ServiceContractCreate(serviceIdentity, serviceAccount, consumerAccount)
	require.NoError(t, err)

	metadata := "some_metadata"
	err = cl.ServiceContractSetMetadata(consumerIdentity, serviceContractID, metadata)
	require.NoError(t, err)

	var baseFee uint64 = 1000
	var variableFee uint64 = 1000
	err = cl.ServiceContractSetFees(serviceIdentity, serviceContractID, baseFee, variableFee)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(serviceIdentity, serviceContractID)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(consumerIdentity, serviceContractID)
	require.NoError(t, err)

	scID, err := cl.GetServiceContractID()
	require.NoError(t, err)
	require.Equal(t, serviceContractID, scID)

	serviceContract, err := cl.GetServiceContract(serviceContractID)
	require.NoError(t, err)
	require.Equal(t, serviceContract.ServiceTwinID, types.U32(serviceTwinID))
	require.Equal(t, serviceContract.ConsumerTwinID, types.U32(consumerTwinID))
	require.Equal(t, serviceContract.BaseFee, types.U64(baseFee))
	require.Equal(t, serviceContract.VariableFee, types.U64(variableFee))
	require.Equal(t, serviceContract.Metadata, metadata)
	require.Equal(t, serviceContract.AcceptedByService, true)
	require.Equal(t, serviceContract.AcceptedByService, true)
	require.Equal(t, serviceContract.State, ServiceContractState{
		IsCreated:        false,
		IsAgreementReady: false,
		IsApprovedByBoth: true,
	})

	// should be able to go to future block to test varaible amount greater than 0
	var variableAmount uint64 = 0
	billMetadata := "some_bill_metadata"
	err = cl.ServiceContractBill(serviceIdentity, serviceContractID, variableAmount, billMetadata)
	require.NoError(t, err)

	err = cl.ServiceContractCancel(consumerIdentity, serviceContractID)
	require.NoError(t, err)
}
