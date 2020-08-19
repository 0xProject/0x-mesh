module github.com/0xProject/0x-mesh

go 1.14

replace (
	// github.com/ethereum/go-ethereum => github.com/0xProject/go-ethereum wasm-support
	github.com/ethereum/go-ethereum => github.com/0xProject/go-ethereum v1.8.8-0.20200603225022-cb1f52043425
	github.com/libp2p/go-flow-metrics => github.com/libp2p/go-flow-metrics v0.0.3
	github.com/libp2p/go-libp2p-pubsub => github.com/0xProject/go-libp2p-pubsub v0.1.1-0.20200228234556-aaa0317e068a
	// github.com/libp2p/go-ws-transport => github.com/0xProject/go-ws-transport upgrade-go-1.14
	github.com/libp2p/go-ws-transport => github.com/0xProject/go-ws-transport v0.1.1-0.20200602173532-300f0ff55a11
	github.com/plaid/go-envvar => github.com/albrow/go-envvar v1.1.1-0.20200123010345-a6ece4436cb7
	// github.com/syndtr/goleveldb => github.com/0xProject/goleveldb upgrade-go-1.14
	github.com/syndtr/goleveldb => github.com/0xProject/goleveldb v1.0.1-0.20200602173211-6ee893c9b83a
)

require (
	github.com/99designs/gqlgen v0.11.3
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/albrow/stringset v2.1.0+incompatible
	github.com/aristanetworks/goarista v0.0.0-20190712234253-ed1100a1c015 // indirect
	github.com/benbjohnson/clock v0.0.0-20161215174838-7dc76406b6d3
	github.com/cespare/cp v1.1.1 // indirect
	github.com/chromedp/cdproto v0.0.0-20190827000638-b5ac1e37ce90
	github.com/chromedp/chromedp v0.4.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/gibson042/canonicaljson-go v1.0.3
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ido50/sqlz v0.0.0-20200308174337-487b8faf612c
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-ds-leveldb v0.4.0
	github.com/ipfs/go-ds-sql v0.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/jpillora/backoff v0.0.0-20170918002102-8eab2debe79d
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/karlseguin/ccache v2.0.3+incompatible
	github.com/karlseguin/expect v1.0.1 // indirect
	github.com/lib/pq v1.3.0
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
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.4.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/multiformats/go-multiaddr v0.2.1
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pkg/errors v0.9.1 // indirect
	github.com/plaid/go-envvar v1.1.0
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/vektah/gqlparser/v2 v2.0.1
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
