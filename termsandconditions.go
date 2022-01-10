package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
)

// AcceptTermsAndConditions accepts terms and conditions
func (s *Substrate) AcceptTermsAndConditions(identity Identity, documentLink string, documentHash string) error {
	cl, meta, err := s.pool.Get()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "TfgridModule.user_accept_tc",
		documentLink, documentHash,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create call")
	}

	blockHash, err := s.Call(cl, meta, identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to accept terms and conditions")
	}

	if err := s.checkForError(cl, meta, blockHash, types.NewAccountID(identity.PublicKey())); err != nil {
		return err
	}

	return nil
}
