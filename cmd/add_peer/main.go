// demo/add_peer is a short program that adds a new peer to 0x Mesh via RPC.
package main

import (
	"github.com/0xProject/0x-mesh/rpc"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type clientEnvVars struct {
	// RPCAddress is the address of the 0x Mesh node to communicate with.
	RPCAddress string `envvar:"RPC_ADDRESS"`
	// PeerID is the base58-encoded peer ID of the peer to connect to.
	PeerID string `envvar:"PEER_ID"`
	// PeerAddr is the Multiaddress of the peer to connect to.
	PeerAddr string `envvar:"PEER_ADDR"`
}

func main() {
	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	// Parse peer ID and peer address.
	parsedPeerID, err := peer.IDB58Decode(env.PeerID)
	if err != nil {
		log.Fatal(err)
	}
	parsedMultiAddr, err := ma.NewMultiaddr(env.PeerAddr)
	if err != nil {
		log.Fatal(err)
	}
	peerInfo := peerstore.PeerInfo{
		ID:    parsedPeerID,
		Addrs: []ma.Multiaddr{parsedMultiAddr},
	}

	client, err := rpc.NewClient(env.RPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}

	if err := client.AddPeer(peerInfo); err != nil {
		log.WithError(err).Fatal("error from AddPeer")
	} else {
		log.Printf("successfully added peer: %s", env.PeerID)
	}
}
