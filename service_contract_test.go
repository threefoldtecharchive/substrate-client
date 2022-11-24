package substrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceContract(t *testing.T) {
	cl := startLocalConnection(t)
	defer cl.Close()

	service_id, err := NewIdentityFromSr25519Phrase(AliceMnemonics)
	require.NoError(t, err)

	consumer_id, err := NewIdentityFromSr25519Phrase(BobMnemonics)
	require.NoError(t, err)

	// service_twinID := assertCreateTwin(t, cl, AliceMnemonics, AliceAddress)
	// consumer_twinID := assertCreateTwin(t, cl, BobMnemonics, BobAddress)

	service_account, err := FromAddress(AliceAddress)
	require.NoError(t, err)

	consumer_account, err := FromAddress(BobAddress)
	require.NoError(t, err)

	contract, err := cl.ServiceContractCreate(service_id, service_account, consumer_account)
	require.NoError(t, err)

	metadata := "some_metadata"
	err = cl.ServiceContractSetMetadata(service_id, contract, metadata)
	require.NoError(t, err)

	var base_fee uint64 = 100
	var variable_fee uint64 = 100
	err = cl.ServiceContractSetFees(service_id, contract, base_fee, variable_fee)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(service_id, contract)
	require.NoError(t, err)

	err = cl.ServiceContractApprove(consumer_id, contract)
	require.NoError(t, err)

	err = cl.ServiceContractCancel(service_id, contract)
	require.NoError(t, err)
}
