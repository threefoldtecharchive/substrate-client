package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	serviceID, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	consumerID, err := NewIdentityFromSr25519Phrase(BobMnemonics)
	require.NoError(t, err)

	// service_twinID := assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)
	// consumer_twinID := assertCreateTwin(t, cl, BobMnemonics, BobAddress)

	serviceAccount, err := FromAddress(AliceAddress)
	require.NoError(t, err)

	consumerAccount, err := FromAddress(BobAddress)
	require.NoError(t, err)

	serviceContractID, err := cl.ServiceContractCreate(serviceID, serviceAccount, consumerAccount)
	require.NoError(t, err)

	metadata := "some_metadata"
	err = cl.ServiceContractSetMetadata(consumerID, serviceContractID, metadata)
	require.NoError(t, err)

	var baseFee uint64 = 100
	var variableFee uint64 = 100
	err = cl.ServiceContractSetFees(serviceID, serviceContractID, baseFee, variableFee)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(serviceID, serviceContractID)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(serviceID, serviceContractID)
	require.NoError(t, err)

	err = cl.ServiceContractCancel(consumerID, serviceContractID)
	require.NoError(t, err)
}
