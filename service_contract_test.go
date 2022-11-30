package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	_ = assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)
	_ = assertCreateTwin(t, cl, BobMnemonics, BobAddress)

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

	var baseFee uint64 = 100
	var variableFee uint64 = 100
	err = cl.ServiceContractSetFees(serviceIdentity, serviceContractID, baseFee, variableFee)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(serviceIdentity, serviceContractID)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(serviceIdentity, serviceContractID)
	require.NoError(t, err)

	err = cl.ServiceContractCancel(consumerIdentity, serviceContractID)
	require.NoError(t, err)
}
