package core

import (
	"context"
	"io"

	ic "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	peer "gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	logging "gx/ipfs/QmRREK2CAZ5Re2Bd9zZFG6FeYDppUWt5cMgsoUEp3ktgSr/go-log"
	goprocess "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	pstore "gx/ipfs/QmeKD8YT7887Xu6Z86iZmpYNxrLogJexqxEugSmaf14k64/go-libp2p-peerstore"
	p2phost "gx/ipfs/QmfH9FKYv3Jp1xiyL8sPchGBUBg6JA6XviwajAo3qgnT3B/go-libp2p-host"
)

var log = logging.Logger("core")

// FogNode is Fogger Core module. It represents an Fogger instance.
type FogNode struct {
	// Self
	Identity peer.ID // the local node's identity

	// Local node
	PrivateKey      ic.PrivKey // the local node's private Key
	PNetFingerprint []byte     // fingerprint of private network

	// Services
	Peerstore pstore.Peerstore // storage for other Peer instances

	// Online
	PeerHost p2phost.Host

	proc goprocess.Process
	ctx  context.Context
}

// Process returns the Process object
func (n *FogNode) Process() goprocess.Process {
	return n.proc
}

// Close calls Close() on the Process object
func (n *FogNode) Close() error {
	return n.proc.Close()
}

// Context returns the IpfsNode context
func (n *FogNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

// teardown closes owned children. If any errors occur, this function returns
// the first error.
func (n *FogNode) teardown() error {
	log.Debug("core is shutting down...")
	// owned objects are closed in this teardown to ensure that they're closed
	// regardless of which constructor was used to add them to the node.
	var closers []io.Closer

	// NOTE: The order that objects are added(closed) matters, if an object
	// needs to use another during its shutdown/cleanup process, it should be
	// closed before that other object

	if n.PeerHost != nil {
		closers = append(closers, n.PeerHost)
	}

	var errs []error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
