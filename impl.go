package substrate

import (
	"fmt"
	"math/rand"
	"net"
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v3"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	//ErrInvalidVersion is returned if version 4bytes is invalid
	ErrInvalidVersion = fmt.Errorf("invalid version")
	//ErrUnknownVersion is returned if version number is not supported
	ErrUnknownVersion = fmt.Errorf("unknown version")
	//ErrNotFound is returned if an object is not found
	ErrNotFound = fmt.Errorf("object not found")
)

// Versioned base for all types
type Versioned struct {
	Version uint32
}

type Conn = *gsrpc.SubstrateAPI
type Meta = *types.Metadata

type Pool interface {
	Get() (Conn, Meta, error)
}

type poolImpl struct {
	urls []string
	cl   *gsrpc.SubstrateAPI
	r    int
	m    sync.Mutex
}

func NewPool(url ...string) (Pool, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("at least one url is required")
	}

	return &poolImpl{
		urls: url,
		r:    rand.Intn(len(url)), // start with random url, then roundrobin
	}, nil
}

// endpoint return the next endpoint to use
// in roundrobin fashion. need to be called
// while lock is acquired.
func (p *poolImpl) endpoint() string {
	defer func() {
		p.r = (p.r + 1) % len(p.urls)
	}()

	return p.urls[p.r]
}

// Get implements Pool interface. Get will try the next url only if the
// client timesout.
func (p *poolImpl) Get() (Conn, Meta, error) {
	// right now this pool implementation just tests the connection
	// makes sure that it is still active, otherwise, tries again
	// until the connection is restored.
	// A better pool implementation can be done later were multiple connections
	// can be handled
	// TODO: thread safety!
	p.m.Lock()
	defer p.m.Unlock()

	for {
		if p.cl == nil {
			cl, err := gsrpc.NewSubstrateAPI(p.endpoint())
			if err != nil {
				return nil, nil, err
			}
			p.cl = cl
		}
		meta, err := p.cl.RPC.State.GetMetadataLatest()
		if errors.Is(err, net.ErrClosed) {
			p.cl = nil
			log.Debug().Msg("reconnecting")
			continue
		} else if err != nil {
			return nil, nil, err
		}

		return p.cl, meta, nil
	}
}

// Substrate client
type Substrate struct {
	pool Pool
}

// NewSubstrate creates a substrate client
func NewSubstrate(url ...string) (*Substrate, error) {
	pool, err := NewPool(url...)
	if err != nil {
		return nil, err
	}

	return &Substrate{
		pool: pool,
	}, nil
}

func (s *Substrate) GetClient() (Conn, Meta, error) {
	return s.pool.Get()
}

func (s *Substrate) getVersion(b types.StorageDataRaw) (uint32, error) {
	var ver Versioned
	if err := types.DecodeFromBytes(b, &ver); err != nil {
		return 0, errors.Wrapf(ErrInvalidVersion, "failed to load version (reason: %s)", err)
	}

	return ver.Version, nil
}
