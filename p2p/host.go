// Note: the code in this file is largely copied from
// https://github.com/libp2p/go-libp2p/config with some modifications in order
// to support configuring the swarm.
//
// Copyright (c) 2014 Juan Batiz-Benet
// Modified work copyright (c) 2019 ZeroEx, Inc.
//
// The code in this file falls under the MIT license located at:
// https://github.com/libp2p/go-libp2p/blob/v0.3.0/LICENSE
//

package p2p

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	csms "github.com/libp2p/go-conn-security-multistream"
	"github.com/libp2p/go-conn-security/insecure"
	libp2p "github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-core/sec"
	discovery "github.com/libp2p/go-libp2p-discovery"
	swarm "github.com/libp2p/go-libp2p-swarm"
	transport "github.com/libp2p/go-libp2p-transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-libp2p/config"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	relay "github.com/libp2p/go-libp2p/p2p/host/relay"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	msmux "github.com/libp2p/go-stream-muxer-multistream"
	ma "github.com/multiformats/go-multiaddr"
)

var p2plog = logging.Logger("p2p-config")

func NewHost(ctx context.Context, opts ...libp2p.Option) (host.Host, *swarm.Swarm, error) {
	var cfg config.Config
	if err := cfg.Apply(append(opts, libp2p.FallbackDefaults)...); err != nil {
		return nil, nil, err
	}

	// Check this early. Prevents us from even *starting* without verifying this.
	if pnet.ForcePrivateNetwork && cfg.Protector == nil {
		p2plog.Error("tried to create a libp2p node with no Private" +
			" Network Protector but usage of Private Networks" +
			" is forced by the enviroment")
		// Note: This is *also* checked the upgrader itself so it'll be
		// enforced even *if* you don't use the libp2p constructor.
		return nil, nil, pnet.ErrNotInPrivateNetwork
	}

	if cfg.PeerKey == nil {
		return nil, nil, fmt.Errorf("no peer key specified")
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(cfg.PeerKey.GetPublic())
	if err != nil {
		return nil, nil, err
	}

	if cfg.Peerstore == nil {
		return nil, nil, fmt.Errorf("no peerstore specified")
	}

	if !cfg.Insecure {
		if err := cfg.Peerstore.AddPrivKey(pid, cfg.PeerKey); err != nil {
			return nil, nil, err
		}
		if err := cfg.Peerstore.AddPubKey(pid, cfg.PeerKey.GetPublic()); err != nil {
			return nil, nil, err
		}
	}

	// TODO: Make the swarm implementation configurable.
	swrm := swarm.NewSwarm(ctx, pid, cfg.Peerstore, cfg.Reporter)
	if cfg.Filters != nil {
		swrm.Filters = cfg.Filters
	}

	h, err := bhost.NewHost(ctx, swrm, &bhost.HostOpts{
		ConnManager:  cfg.ConnManager,
		AddrsFactory: cfg.AddrsFactory,
		NATManager:   cfg.NATManager,
		EnablePing:   !cfg.DisablePing,
	})

	if err != nil {
		swrm.Close()
		return nil, nil, err
	}

	if cfg.Relay {
		// If we've enabled the relay, we should filter out relay
		// addresses by default.
		//
		// TODO: We shouldn't be doing this here.
		oldFactory := h.AddrsFactory
		h.AddrsFactory = func(addrs []ma.Multiaddr) []ma.Multiaddr {
			return oldFactory(relay.Filter(addrs))
		}
	}

	upgrader := new(tptu.Upgrader)
	upgrader.Protector = cfg.Protector
	upgrader.Filters = swrm.Filters
	if cfg.Insecure {
		upgrader.Secure = makeInsecureTransport(pid)
	} else {
		upgrader.Secure, err = makeSecurityTransport(h, cfg.SecurityTransports)
		if err != nil {
			h.Close()
			return nil, nil, err
		}
	}

	upgrader.Muxer, err = makeMuxer(h, cfg.Muxers)
	if err != nil {
		h.Close()
		return nil, nil, err
	}

	tpts, err := makeTransports(h, upgrader, cfg.Transports)
	if err != nil {
		h.Close()
		return nil, nil, err
	}
	for _, t := range tpts {
		err = swrm.AddTransport(t)
		if err != nil {
			h.Close()
			return nil, nil, err
		}
	}

	if cfg.Relay {
		err := circuit.AddRelayTransport(swrm.Context(), h, upgrader, cfg.RelayOpts...)
		if err != nil {
			h.Close()
			return nil, nil, err
		}
	}

	// TODO: This method succeeds if listening on one address succeeds. We
	// should probably fail if listening on *any* addr fails.
	if err := h.Network().Listen(cfg.ListenAddrs...); err != nil {
		h.Close()
		return nil, nil, err
	}

	// Configure routing and autorelay
	var router routing.PeerRouting
	if cfg.Routing != nil {
		router, err = cfg.Routing(h)
		if err != nil {
			h.Close()
			return nil, nil, err
		}
	}

	if cfg.EnableAutoRelay {
		if !cfg.Relay {
			h.Close()
			return nil, nil, fmt.Errorf("cannot enable autorelay; relay is not enabled")
		}

		if router == nil {
			h.Close()
			return nil, nil, fmt.Errorf("cannot enable autorelay; no routing for discovery")
		}

		crouter, ok := router.(routing.ContentRouting)
		if !ok {
			h.Close()
			return nil, nil, fmt.Errorf("cannot enable autorelay; no suitable routing for discovery")
		}

		discovery := discovery.NewRoutingDiscovery(crouter)

		hop := false
		for _, opt := range cfg.RelayOpts {
			if opt == circuit.OptHop {
				hop = true
				break
			}
		}

		if hop {
			// advertise ourselves
			relay.Advertise(ctx, discovery)
		} else {
			_ = relay.NewAutoRelay(swrm.Context(), h, discovery, router)
		}
	}

	// start the host background tasks
	h.Start()

	if router != nil {
		return routed.Wrap(h, router), swrm, nil
	}
	return h, swrm, nil
}

func makeInsecureTransport(id peer.ID) sec.SecureTransport {
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(insecure.ID, insecure.New(id))
	return secMuxer
}

func makeSecurityTransport(h host.Host, tpts []config.MsSecC) (sec.SecureTransport, error) {
	secMuxer := new(csms.SSMuxer)
	transportSet := make(map[string]struct{}, len(tpts))
	for _, tptC := range tpts {
		if _, ok := transportSet[tptC.ID]; ok {
			return nil, fmt.Errorf("duplicate security transport: %s", tptC.ID)
		}
		transportSet[tptC.ID] = struct{}{}
	}
	for _, tptC := range tpts {
		tpt, err := tptC.SecC(h)
		if err != nil {
			return nil, err
		}
		if _, ok := tpt.(*insecure.Transport); ok {
			return nil, fmt.Errorf("cannot construct libp2p with an insecure transport, set the Insecure config option instead")
		}
		secMuxer.AddTransport(tptC.ID, tpt)
	}
	return secMuxer, nil
}

func makeMuxer(h host.Host, tpts []config.MsMuxC) (mux.Multiplexer, error) {
	muxMuxer := msmux.NewBlankTransport()
	transportSet := make(map[string]struct{}, len(tpts))
	for _, tptC := range tpts {
		if _, ok := transportSet[tptC.ID]; ok {
			return nil, fmt.Errorf("duplicate muxer transport: %s", tptC.ID)
		}
		transportSet[tptC.ID] = struct{}{}
	}
	for _, tptC := range tpts {
		tpt, err := tptC.MuxC(h)
		if err != nil {
			return nil, err
		}
		muxMuxer.AddTransport(tptC.ID, tpt)
	}
	return muxMuxer, nil
}

func makeTransports(h host.Host, u *tptu.Upgrader, tpts []config.TptC) ([]transport.Transport, error) {
	transports := make([]transport.Transport, len(tpts))
	for i, tC := range tpts {
		t, err := tC(h, u)
		if err != nil {
			return nil, err
		}
		transports[i] = t
	}
	return transports, nil
}
