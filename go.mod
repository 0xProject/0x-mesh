module github.com/0xProject/0x-mesh

go 1.13

replace (
	github.com/libp2p/go-flow-metrics => github.com/libp2p/go-flow-metrics v0.0.3
	github.com/libp2p/go-libp2p-pubsub => github.com/0xProject/go-libp2p-pubsub v0.1.1-0.20200228234556-aaa0317e068a
	github.com/libp2p/go-ws-transport => github.com/0xProject/go-ws-transport v0.1.1-0.20200201000210-2db3396fec39
	github.com/plaid/go-envvar => github.com/albrow/go-envvar v1.1.1-0.20200123010345-a6ece4436cb7
	github.com/syndtr/goleveldb => github.com/0xProject/goleveldb v1.0.1-0.20191115232649-6a187a47701c
)

require (
	github.com/0xProject/sql-datastore v0.0.0-20200129193319-32397013f115
	github.com/albrow/stringset v2.1.0+incompatible
	github.com/benbjohnson/clock v0.0.0-20161215174838-7dc76406b6d3
	github.com/cespare/cp v1.1.1 // indirect
	github.com/chromedp/cdproto v0.0.0-20190827000638-b5ac1e37ce90
	github.com/chromedp/chromedp v0.4.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.10.8
	github.com/gibson042/canonicaljson-go v1.0.3
	github.com/google/uuid v1.1.5
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/ipfs/go-datastore v0.3.1
	github.com/ipfs/go-ds-leveldb v0.4.0
	github.com/jpillora/backoff v0.0.0-20170918002102-8eab2debe79d
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/karlseguin/ccache v2.0.3+incompatible
	github.com/karlseguin/expect v1.0.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.2.0
	github.com/libp2p/go-libp2p v0.5.1
	github.com/libp2p/go-libp2p-autonat-svc v0.1.0
	github.com/libp2p/go-libp2p-circuit v0.1.4
	github.com/libp2p/go-libp2p-connmgr v0.2.1
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/libp2p/go-libp2p-discovery v0.2.0
	github.com/libp2p/go-libp2p-kad-dht v0.5.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.4
	github.com/libp2p/go-libp2p-pubsub v0.2.5
	github.com/libp2p/go-libp2p-swarm v0.2.2
	github.com/libp2p/go-maddr-filter v0.0.5
	github.com/libp2p/go-tcp-transport v0.1.1
	github.com/libp2p/go-ws-transport v0.2.0
	github.com/multiformats/go-multiaddr v0.2.0
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/ocdogan/rbt v0.0.0-20160425054511-de6e2b48be33
	github.com/plaid/go-envvar v1.1.0
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
)
