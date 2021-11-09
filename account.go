package substrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vedhavyas/go-subkey"
)

const (
	network        = 42
	KeyTypeEd25519 = "ed25519"
	KeyTypeSr25519 = "sr25519"
)

// AccountID type
type AccountID types.AccountID

//PublicKey gets public key from account id
func (a AccountID) PublicKey() []byte {
	return a[:]
}

// String return string representation of account
func (a AccountID) String() string {
	address, _ := subkey.SS58Address(a[:], network)
	return address
}

// MarshalJSON implementation
func (a AccountID) MarshalJSON() ([]byte, error) {
	address, err := subkey.SS58Address(a[:], network)
	if err != nil {
		return nil, err
	}

	return json.Marshal(address)
}

// FromAddress creates an AccountID from a SS58 address
func FromAddress(address string) (account AccountID, err error) {
	bytes := base58.Decode(address)
	if len(bytes) != 3+len(account) {
		return account, fmt.Errorf("invalid address length")
	}
	if bytes[0] != network {
		return account, fmt.Errorf("invalid address format")
	}

	copy(account[:], bytes[1:len(account)+1])
	return
}

func FromKeyBytes(address []byte) (string, error) {
	return subkey.SS58Address(address, network)
}

// keyringPairFromSecret creates KeyPair based on seed/phrase and network
// Leave network empty for default behavior
func keyringPairFromSecret(seedOrPhrase string, network uint8, keyType string) (signature.KeyringPair, error) {
	scheme, err := keyScheme(keyType)
	if err != nil {
		return signature.KeyringPair{}, err
	}
	kyr, err := subkey.DeriveKeyPair(scheme, seedOrPhrase)

	if err != nil {
		return signature.KeyringPair{}, err
	}

	ss58Address, err := kyr.SS58Address(network)
	if err != nil {
		return signature.KeyringPair{}, err
	}

	var pk = kyr.Public()

	return signature.KeyringPair{
		URI:       seedOrPhrase,
		Address:   ss58Address,
		PublicKey: pk,
	}, nil
}

var (
	ErrAccountNotFound = fmt.Errorf("account not found")
)

/*
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kycSignature": "", "data": {"name": "", "email": ""}, "substrateAccountID": "5DAprR72N6s7AWGwN7TzV9MyuyGk9ifrq8kVxoXG9EYWpic4"}' \
  https://api.substrate01.threefold.io/activate
*/

func (s *Substrate) activateAccount(identity *Identity, activationURL string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]string{
		"substrateAccountID": identity.Address,
	}); err != nil {
		return errors.Wrap(err, "failed to build required body")
	}

	response, err := http.Post(activationURL, "application/json", &buf)
	if err != nil {
		return errors.Wrap(err, "failed to call activation service")
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusConflict {
		// it went fine.
		return nil
	}

	return fmt.Errorf("failed to activate account: %s", response.Status)
}

// EnsureAccount makes sure account is available on blockchain
// if not, it uses activation service to create one
func (s *Substrate) EnsureAccount(identity *Identity, activationURL string) (info types.AccountInfo, err error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return info, err
	}
	info, err = s.getAccount(cl, meta, identity)
	if errors.Is(err, ErrAccountNotFound) {
		// account activation
		log.Debug().Msg("account not found ... activating")
		if err = s.activateAccount(identity, activationURL); err != nil {
			return
		}

		// after activation this can take up to 10 seconds
		// before the account is actually there !

		exp := backoff.NewExponentialBackOff()
		exp.MaxElapsedTime = 10 * time.Second
		exp.MaxInterval = 3 * time.Second

		err = backoff.Retry(func() error {
			info, err = s.getAccount(cl, meta, identity)
			return err
		}, exp)

		return
	}

	return

}

// Identity is a user identity
type Identity struct {
	signature.KeyringPair
	keyType string
}

// SecureKey returns subkey key pair from identity
func (i *Identity) KeyPair() (subkey.KeyPair, error) {
	scheme, err := keyScheme(i.keyType)
	if err != nil {
		return nil, err
	}
	kyr, err := subkey.DeriveKeyPair(scheme, i.URI)
	if err != nil {
		return nil, err
	}

	return kyr, nil
}

// IdentityFromSecureKey derive the correct substrate identity from ed25519 or sr25519 key
func IdentityFromSecureKey(sk []byte, keyType string) (Identity, error) {
	seed := sk[:32]
	str := types.HexEncodeToString(seed)
	krp, err := keyringPairFromSecret(str, network, keyType)
	if err != nil {
		return Identity{}, err
	}

	return Identity{krp, keyType}, nil
	// because 42 is the answer to life the universe and everything
	// no, seriously, don't change it, it has to be 42.
}

//IdentityFromPhrase gets identity from hex seed or mnemonics
func IdentityFromPhrase(seedOrPhrase, keyType string) (Identity, error) {
	krp, err := keyringPairFromSecret(seedOrPhrase, network, keyType)
	if err != nil {
		return Identity{}, err
	}

	return Identity{krp, keyType}, nil
}

func (s *Substrate) getAccount(cl Conn, meta Meta, identity *Identity) (info types.AccountInfo, err error) {
	key, err := types.CreateStorageKey(meta, "System", "Account", identity.PublicKey, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create storage key")
		return
	}

	ok, err := cl.RPC.State.GetStorageLatest(key, &info)
	if err != nil || !ok {
		if !ok {
			return info, ErrAccountNotFound
		}

		return
	}

	return
}

// GetAccount gets account info with secure key
func (s *Substrate) GetAccount(identity *Identity) (info types.AccountInfo, err error) {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return info, err
	}

	return s.getAccount(cl, meta, identity)
}
