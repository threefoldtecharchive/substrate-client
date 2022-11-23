package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// GetZosVersion gets the latest version for each network
func (s *Substrate) GetZosVersion() (string, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return "", err
	}

	key, err := types.CreateStorageKey(meta, "TfgridModule", "ZosVersion", nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create substrate query key")
	}

	raw, err := cl.RPC.State.GetStorageRawLatest(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup entity")
	}

	if len(*raw) == 0 {
		return "", errors.Wrap(ErrNotFound, "zos version not found")
	}

	var zosVersion string

	if err := types.Decode(*raw, &zosVersion); err != nil {
		return "", errors.Wrap(err, "failed to load object")
	}

	return zosVersion, nil
}

// CreateTwin creates a twin
func (s *Substrate) SetZosVersion(identity Identity, version string) (string, error) {
	cl, meta, err := s.getClient()
	if err != nil {
		return "", err
	}

	c, err := types.NewCall(meta, "TfgridModule.set_zos_version", version)
	if err != nil {
		return "", errors.Wrap(err, "failed to create call")
	}

	if _, err := s.Call(cl, meta, identity, c); err != nil {
		return "", errors.Wrap(err, "failed to set zos version")
	}

	return s.GetZosVersion()
}
