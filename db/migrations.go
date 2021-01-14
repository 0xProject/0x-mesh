// +build !js

package db

// TODO(albrow): Use a proper migration tool. We don't technically need this
// now but it will be necessary if we ever change the database schema.
// Note(albrow): If needed, we can optimize this by adding indexes to the
// orders and miniHeaders tables.
const schema = `
CREATE TABLE IF NOT EXISTS orders (
	hash                     TEXT UNIQUE NOT NULL,
	chainID                  TEXT NOT NULL,
	exchangeAddress          TEXT NOT NULL,
	makerAddress             TEXT NOT NULL,
	makerAssetData           TEXT NOT NULL,
	makerFeeAssetData        TEXT NOT NULL,
	makerAssetAmount         TEXT NOT NULL,
	makerFee                 TEXT NOT NULL,
	takerAddress             TEXT NOT NULL,
	takerAssetData           TEXT NOT NULL,
	takerFeeAssetData        TEXT NOT NULL,
	takerAssetAmount         TEXT NOT NULL,
	takerFee                 TEXT NOT NULL,
	senderAddress            TEXT NOT NULL,
	feeRecipientAddress      TEXT NOT NULL,
	expirationTimeSeconds    TEXT NOT NULL,
	salt                     TEXT NOT NULL,
	signature                TEXT NOT NULL,
	lastUpdated              DATETIME NOT NULL,
	fillableTakerAssetAmount TEXT NOT NULL,
	isRemoved                BOOLEAN NOT NULL,
	isPinned                 BOOLEAN NOT NULL,
	isUnfillable             BOOLEAN NOT NULL,
	isExpired                BOOLEAN NOT NULL,
	parsedMakerAssetData     TEXT NOT NULL,
	parsedMakerFeeAssetData  TEXT NOT NULL,
	lastValidatedBlockNumber TEXT NOT NULL,
	lastValidatedBlockHash   TEXT NOT NULL,
	keepCancelled            BOOLEAN NOT NULL,
	keepExpired              BOOLEAN NOT NULL,
	keepFullyFilled          BOOLEAN NOT NULL,
	keepUnfunded             BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS miniHeaders (
	hash      TEXT UNIQUE NOT NULL,
	number    TEXT UNIQUE NOT NULL,
	parent    TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	logs      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata (
	ethereumChainID                   BIGINT NOT NULL,
	ethRPCRequestsSentInCurrentUTCDay BIGINT NOT NULL,
	startOfCurrentUTCDay              DATETIME NOT NULL
);
`

const v4OrdersSchema = `
CREATE TABLE IF NOT EXISTS ordersv4 (
	hash                     TEXT UNIQUE NOT NULL,
	chainID                  TEXT NOT NULL,
	exchangeAddress          TEXT NOT NULL,
	makerToken               TEXT NOT NULL,
	takerToken               TEXT NOT NULL,
	makerAmount              TEXT NOT NULL,
	takerAmount              TEXT NOT NULL,
	takerTokenFeeAmount      TEXT NOT NULL,
	makerAddress             TEXT NOT NULL,
	takerAddress             TEXT NOT NULL,
	sender                   TEXT NOT NULL,
	feeRecipient             TEXT NOT NULL,
	pool                     TEXT NOT NULL,
	expiry                   TEXT NOT NULL,
	salt                     TEXT NOT NULL,
	signature                TEXT NOT NULL,
	lastUpdated              DATETIME NOT NULL,
	fillableTakerAssetAmount TEXT NOT NULL,
	isRemoved                BOOLEAN NOT NULL,
	isPinned                 BOOLEAN NOT NULL,
	isUnfillable             BOOLEAN NOT NULL,
	isExpired                BOOLEAN NOT NULL,
	parsedMakerAssetData     TEXT NOT NULL,
	parsedMakerFeeAssetData  TEXT NOT NULL,
	lastValidatedBlockNumber TEXT NOT NULL,
	lastValidatedBlockHash   TEXT NOT NULL,
	keepCancelled            BOOLEAN NOT NULL,
	keepExpired              BOOLEAN NOT NULL,
	keepFullyFilled          BOOLEAN NOT NULL,
	keepUnfunded             BOOLEAN NOT NULL
);
`

const peerstoreSchema = `
CREATE TABLE IF NOT EXISTS peerstore (
	key  TEXT NOT NULL UNIQUE,
	data BYTEA NOT NULL
);
`
const dhtSchema = `
CREATE TABLE IF NOT EXISTS dhtstore (
	key  TEXT NOT NULL UNIQUE,
	data BYTEA NOT NULL
);
`

// Note(albrow): If needed, we can optimize this by using prepared
// statements for inserts instead of just a string.
const insertOrderQuery = `INSERT INTO orders (
	hash,
	chainID,
	exchangeAddress,
	makerAddress,
	makerAssetData,
	makerFeeAssetData,
	makerAssetAmount,
	makerFee,
	takerAddress,
	takerAssetData,
	takerFeeAssetData,
	takerAssetAmount,
	takerFee,
	senderAddress,
	feeRecipientAddress,
	expirationTimeSeconds,
	salt,
	signature,
	lastUpdated,
	fillableTakerAssetAmount,
	isRemoved,
	isPinned,
	isUnfillable,
	isExpired,
	parsedMakerAssetData,
	parsedMakerFeeAssetData,
	lastValidatedBlockNumber,
	lastValidatedBlockHash,
	keepCancelled,
	keepExpired,
	keepFullyFilled,
	keepUnfunded
) VALUES (
	:hash,
	:chainID,
	:exchangeAddress,
	:makerAddress,
	:makerAssetData,
	:makerFeeAssetData,
	:makerAssetAmount,
	:makerFee,
	:takerAddress,
	:takerAssetData,
	:takerFeeAssetData,
	:takerAssetAmount,
	:takerFee,
	:senderAddress,
	:feeRecipientAddress,
	:expirationTimeSeconds,
	:salt,
	:signature,
	:lastUpdated,
	:fillableTakerAssetAmount,
	:isRemoved,
	:isPinned,
	:isUnfillable,
	:isExpired,
	:parsedMakerAssetData,
	:parsedMakerFeeAssetData,
	:lastValidatedBlockNumber,
	:lastValidatedBlockHash,
	:keepCancelled,
	:keepExpired,
	:keepFullyFilled,
	:keepUnfunded
) ON CONFLICT DO NOTHING
`

const updateOrderQuery = `UPDATE orders SET
	chainID = :chainID,
	exchangeAddress = :exchangeAddress,
	makerAddress = :makerAddress,
	makerAssetData = :makerAssetData,
	makerFeeAssetData = :makerFeeAssetData,
	makerAssetAmount = :makerAssetAmount,
	makerFee = :makerFee,
	takerAddress = :takerAddress,
	takerAssetData = :takerAssetData,
	takerFeeAssetData = :takerFeeAssetData,
	takerAssetAmount = :takerAssetAmount,
	takerFee = :takerFee,
	senderAddress = :senderAddress,
	feeRecipientAddress = :feeRecipientAddress,
	expirationTimeSeconds = :expirationTimeSeconds,
	salt = :salt,
	signature = :signature,
	lastUpdated = :lastUpdated,
	fillableTakerAssetAmount = :fillableTakerAssetAmount,
	isRemoved = :isRemoved,
	isPinned = :isPinned,
	isUnfillable = :isUnfillable,
	isExpired = :isExpired,
	parsedMakerAssetData = :parsedMakerAssetData,
	parsedMakerFeeAssetData = :parsedMakerFeeAssetData,
	lastValidatedBlockNumber = :lastValidatedBlockNumber,
	lastValidatedBlockHash = :lastValidatedBlockHash,
	keepCancelled = :keepCancelled,
	keepExpired = :keepExpired,
	keepFullyFilled = :keepFullyFilled,
	keepUnfunded = :keepUnfunded
WHERE orders.hash = :hash
`

const insertMiniHeaderQuery = `INSERT INTO miniHeaders (
	hash,
	parent,
	number,
	timestamp,
	logs
) VALUES (
	:hash,
	:parent,
	:number,
	:timestamp,
	:logs
) ON CONFLICT DO NOTHING`

const insertMetadataQuery = `INSERT INTO metadata (
	ethereumChainID,
	ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay
) VALUES (
	:ethereumChainID,
	:ethRPCRequestsSentInCurrentUTCDay,
	:startOfCurrentUTCDay
)`

const updateMetadataQuery = `UPDATE metadata SET
	ethereumChainID = :ethereumChainID,
	ethRPCRequestsSentInCurrentUTCDay = :ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay = :startOfCurrentUTCDay
`
