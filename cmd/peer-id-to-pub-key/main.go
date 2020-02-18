package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	peer "github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expects exactly one argument")
	}
	peerIDString := os.Args[1]
	peerID, err := peer.IDB58Decode(peerIDString)
	if err != nil {
		log.Fatal(err)
	}
	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	rawPubKey, err := pubKey.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(rawPubKey))
}
