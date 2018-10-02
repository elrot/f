package core

import (
	"context"
	"os"
	"syscall"

	goprocessctx "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess/context"
	pstore "gx/ipfs/QmeKD8YT7887Xu6Z86iZmpYNxrLogJexqxEugSmaf14k64/go-libp2p-peerstore"
)

// NewNode constructs and returns an FogNode using the given cfg.
func NewNode(ctx context.Context) (*FogNode, error) {

	n := &FogNode{
		Peerstore: pstore.NewPeerstore(),
		ctx:       ctx,
	}

	// TODO: this is a weird circular-ish dependency, rework it
	n.proc = goprocessctx.WithContextAndTeardown(ctx, n.teardown)

	// if err := setupNode(ctx, n); err != nil {
	// 	n.Close()
	// 	return nil, err
	// }

	return n, nil
}

func isTooManyFDError(err error) bool {
	perr, ok := err.(*os.PathError)
	if ok && perr.Err == syscall.EMFILE {
		return true
	}

	return false
}
