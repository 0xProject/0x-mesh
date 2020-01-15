module github.com/0xProject/0x-mesh

go 1.12

replace (
	github.com/ethereum/go-ethereum => github.com/0xProject/go-ethereum v1.8.8-0.20191104224527-9d5c202240be
	github.com/libp2p/go-ws-transport => github.com/libp2p/go-ws-transport v0.0.0-20191008032742-3098bba549e8
	github.com/syndtr/goleveldb => github.com/0xProject/goleveldb v1.0.1-0.20191115232649-6a187a47701c
)

require (
	github.com/albrow/stringset v2.1.0+incompatible
	github.com/allegro/bigcache v0.0.0-20190618191010-69ea0af04088
	github.com/aristanetworks/goarista v0.0.0-20190712234253-ed1100a1c015
	github.com/benbjohnson/clock v0.0.0-20161215174838-7dc76406b6d3
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
	github.com/chromedp/cdproto v0.0.0-20190827000638-b5ac1e37ce90
	github.com/chromedp/chromedp v0.4.0
	github.com/coreos/go-semver v0.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/edsrzf/mmap-go v1.0.0
	github.com/elastic/gosigar v0.10.5
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fd/go-nat v1.0.0 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20190528105824-2fd9b619dd3c
	github.com/go-stack/stack v1.8.0
	github.com/gobwas/httphead v0.0.0-20180130184737-2c6c146eadee
	github.com/gobwas/pool v0.2.0
	github.com/gobwas/ws v1.0.2
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/golang-lru v0.5.3
	github.com/huin/goupnp v1.0.0
	github.com/ipfs/go-cid v0.0.3
	github.com/ipfs/go-datastore v0.1.1
	github.com/ipfs/go-ds-leveldb v0.0.2
	github.com/ipfs/go-ipfs-util v0.0.1
	github.com/ipfs/go-log v0.0.1
	github.com/ipfs/go-todocounter v0.0.1
	github.com/jackpal/gateway v1.0.5
	github.com/jackpal/go-nat-pmp v1.0.1
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jbenet/go-temp-err-catcher v0.0.0-20150120210811-aac704a3f4f2
	github.com/jbenet/goprocess v0.1.3
	github.com/jpillora/backoff v0.0.0-20170918002102-8eab2debe79d
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9
	github.com/karlseguin/ccache v2.0.3+incompatible
	github.com/knq/sysutil v0.0.0-20181215143952-f05b59f0f307
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/koron/go-ssdp v0.0.0-20180514024734-4a0ed625a78b
	github.com/lib/pq v1.2.0
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-buffer-pool v0.0.2
	github.com/libp2p/go-conn-security v0.0.0-20190226201940-b2fb4ac68c41 // indirect
	github.com/libp2p/go-conn-security-multistream v0.1.0
	github.com/libp2p/go-eventbus v0.1.0
	github.com/libp2p/go-flow-metrics v0.0.1
	github.com/libp2p/go-libp2p v0.3.1
	github.com/libp2p/go-libp2p-autonat v0.1.0
	github.com/libp2p/go-libp2p-autonat-svc v0.1.0
	github.com/libp2p/go-libp2p-circuit v0.1.2
	github.com/libp2p/go-libp2p-connmgr v0.1.1
	github.com/libp2p/go-libp2p-core v0.2.4
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-discovery v0.2.0
	github.com/libp2p/go-libp2p-kad-dht v0.3.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/libp2p/go-libp2p-protocol v0.0.0-20171212212132-b29f3d97e3a2 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.1.0
	github.com/libp2p/go-libp2p-record v0.1.1
	github.com/libp2p/go-libp2p-routing v0.1.0
	github.com/libp2p/go-libp2p-secio v0.2.0
	github.com/libp2p/go-libp2p-swarm v0.2.1
	github.com/libp2p/go-libp2p-transport v0.0.0-20190226201958-e8580c8a519d // indirect
	github.com/libp2p/go-libp2p-transport-upgrader v0.1.1
	github.com/libp2p/go-libp2p-yamux v0.2.1
	github.com/libp2p/go-maddr-filter v0.0.5
	github.com/libp2p/go-mplex v0.1.0
	github.com/libp2p/go-msgio v0.0.4
	github.com/libp2p/go-nat v0.0.3
	github.com/libp2p/go-openssl v0.0.3
	github.com/libp2p/go-reuseport v0.0.1
	github.com/libp2p/go-reuseport-transport v0.0.2
	github.com/libp2p/go-stream-muxer-multistream v0.2.0
	github.com/libp2p/go-tcp-transport v0.1.1
	github.com/libp2p/go-testutil v0.0.0-20190226202041-873eaa1a26ba // indirect
	github.com/libp2p/go-ws-transport v0.1.0
	github.com/libp2p/go-yamux v1.2.3
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e
	github.com/mattn/go-colorable v0.1.2
	github.com/mattn/go-isatty v0.0.8
	github.com/mattn/go-runewidth v0.0.4
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/minio/sha256-simd v0.1.1
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-base32 v0.0.3
	github.com/multiformats/go-multiaddr v0.1.2
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/multiformats/go-multiaddr-fmt v0.1.0
	github.com/multiformats/go-multiaddr-net v0.1.0
	github.com/multiformats/go-multibase v0.0.1
	github.com/multiformats/go-multihash v0.0.8
	github.com/multiformats/go-multistream v0.1.0
	github.com/multiformats/go-varint v0.0.1
	github.com/ocdogan/rbt v0.0.0-20160425054511-de6e2b48be33
	github.com/olekukonko/tablewriter v0.0.1
	github.com/opaolini/go-ds-sql v0.0.0-20191105113501-15a4c8e2fea5
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pkg/errors v0.8.1
	github.com/plaid/go-envvar v1.1.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/tsdb v0.10.0
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572
	github.com/spaolacci/murmur3 v1.1.0
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/steakknife/bloomfilter v0.0.0-20180906043351-99ee86d9200f
	github.com/steakknife/hamming v0.0.0-20180906055317-003c143a81c2
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1
	github.com/whyrusleeping/go-logging v0.0.1
	github.com/whyrusleeping/go-notifier v0.0.0-20170827234753-097c5d47330f
	github.com/whyrusleeping/go-smux-multiplex v3.0.16+incompatible // indirect
	github.com/whyrusleeping/go-smux-multistream v2.0.2+incompatible // indirect
	github.com/whyrusleeping/go-smux-yamux v2.0.8+incompatible // indirect
	github.com/whyrusleeping/mafmt v1.2.8
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/whyrusleeping/timecache v0.0.0-20160911033111-cfcb2f1abfee
	github.com/whyrusleeping/yamux v1.1.5 // indirect
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	github.com/xeipuuv/gojsonschema v1.1.0
	go.opencensus.io v0.22.1
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)
