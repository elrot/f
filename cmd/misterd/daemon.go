package main

import (
	"context"
	"fmt"
	"time"

	"github.com/elrot/mister/core"
)

func daemon(ctx context.Context) {
	// let the user know we're going.
	fmt.Printf("Initializing daemon...\n")

	go func() {
		<-ctx.Done()
		fmt.Println("Received interrupt signal, shutting down...")
		fmt.Println("(Hit ctrl-c again to force-shutdown the daemon.)")
	}()

	node, err := core.NewNode(ctx)

	if err != nil {
		log.Error("error from node construction: ", err)
		return
	}

	defer func() {
		// We wait for the node to close first, as the node has children
		// that it will wait for before closing, such as the API server.
		node.Close()

		select {
		case <-ctx.Done():
			log.Info("Gracefully shut down daemon")
		default:
		}
	}()

	rpErrc, err := runLogger(node)
	if err != nil {
		return
	}

	fmt.Printf("Daemon is ready\n")
	// collect long-running errors and block for shutdown
	for err := range merge(rpErrc) {
		if err != nil {
			log.Error(err)
		}
	}
}

func runLogger(node *core.FogNode) (<-chan error, error) {
	errc := make(chan error)
	go func() {
		errc <- logger(node, "2s")
		close(errc)
	}()
	return errc, nil
}

func logger(node *core.FogNode, period string) error {
	p, err := time.ParseDuration(period)
	if err != nil {
		return err
	}

	for {
		select {
		case <-node.Context().Done():
			return nil
		case <-time.After(p):
			log.Info(fmt.Sprintf("Identity: %s", node.Identity.Pretty()))
		}
	}
}
