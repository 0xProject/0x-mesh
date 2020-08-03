import { SignedOrder } from '@0x/order-utils';
import { BigNumber } from '@0x/utils';
import { SupportedProvider, ZeroExProvider } from 'ethereum-types';

export { SignedOrder } from '@0x/order-utils';
export { BigNumber } from '@0x/utils';
export { SupportedProvider } from 'ethereum-types';

/** @ignore */
export interface WrapperGetOrdersResponse {
    timestamp: string;
    ordersInfos: WrapperOrderInfo[];
}

export interface GetOrdersResponse {
    timestamp: number;
    ordersInfos: OrderInfo[];
}

/** @ignore */
export interface WrapperOrderInfo {
    orderHash: string;
    signedOrder: WrapperSignedOrder;
    fillableTakerAssetAmount: string;
}

export interface OrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: BigNumber;
}

/**
 * An interface for JSON schema types, which are used for custom order filters.
 */
export interface JsonSchema {
    id?: string;
    $schema?: string;
    $ref?: string;
    title?: string;
    description?: string;
    multipleOf?: number;
    maximum?: number;
    exclusiveMaximum?: boolean;
    minimum?: number;
    exclusiveMinimum?: boolean;
    maxLength?: number;
    minLength?: number;
    pattern?: string | RegExp;
    additionalItems?: boolean | JsonSchema;
    items?: JsonSchema | JsonSchema[];
    maxItems?: number;
    minItems?: number;
    uniqueItems?: boolean;
    maxProperties?: number;
    minProperties?: number;
    required?: string[];
    additionalProperties?: boolean | JsonSchema;
    definitions?: {
        [name: string]: JsonSchema;
    };
    properties?: {
        [name: string]: JsonSchema;
    };
    patternProperties?: {
        [name: string]: JsonSchema;
    };
    dependencies?: {
        [name: string]: JsonSchema | string[];
    };
    enum?: any[];
    // NOTE(albrow): This interface type is based on
    // https://github.com/tdegrunt/jsonschema/blob/9cb2cf847a33abb76b694c6ed4d8d12ef2037201/lib/index.d.ts#L50
    // but modified to include the 'const' field from the JSON Schema
    // specification draft 6 (https://json-schema.org/understanding-json-schema/reference/generic.html#constant-values)
    // See also: https://github.com/tdegrunt/jsonschema/issues/271
    const?: any;
    type?: string | string[];
    format?: string;
    allOf?: JsonSchema[];
    anyOf?: JsonSchema[];
    oneOf?: JsonSchema[];
    not?: JsonSchema;
}

// Note(albrow): This is currently copied over from core/core.go. We need to keep
// both definitions in sync, so if you change one you must also change the
// other.
/**
 * A set of configuration options for Mesh.
 */
export interface Config {
    // Verbosity is the logging verbosity. Defaults to Verbosity.Error meaning
    // only errors will be logged.
    verbosity?: Verbosity;
    // The URL of an Ethereum node which supports the Ethereum JSON RPC API.
    // Used to validate and watch orders.
    ethereumRPCURL?: string;
    // EthereumChainID is the chain ID specifying which Ethereum chain you wish to
    // run your Mesh node for
    ethereumChainID: number;
    // UseBootstrapList is whether to bootstrap the DHT by connecting to a
    // specific set of peers.
    useBootstrapList?: boolean;
    // bootstrapList is a list of multiaddresses to use for bootstrapping the
    // DHT (e.g.,
    // "/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF").
    // Defaults to the hard-coded default bootstrap list.
    bootstrapList?: string[];
    // The polling interval (in seconds) to wait before checking for a new
    // Ethereum block that might contain transactions that impact the
    // fillability of orders stored by Mesh. Different chains have different
    // block producing intervals: POW chains are typically slower (e.g.,
    // Mainnet) and POA chains faster (e.g., Kovan) so one should adjust the
    // polling interval accordingly. Defaults to 5.
    blockPollingIntervalSeconds?: number;
    // The maximum request Content-Length accepted by the backing Ethereum RPC
    // endpoint used by Mesh. Geth & Infura both limit a request's content
    // length to 1024 * 512 Bytes. Parity and Alchemy have much higher limits.
    // When batch validating 0x orders, we will fit as many orders into a
    // request without crossing the max content length. The default value is
    // appropriate for operators using Geth or Infura. If using Alchemy or
    // Parity, feel free to double the default max in order to reduce the number
    // of RPC calls made by Mesh. Defaults to 524288 bytes.
    ethereumRPCMaxContentLength?: number;
    // Determines whether or not Mesh should limit the number of Ethereum RPC
    // requests it sends. It defaults to true. Disabling Ethereum RPC rate
    // limiting can reduce latency for receiving order events in some network
    // conditions, but can also potentially lead to higher costs or other rate
    // limiting issues outside of Mesh, depending on your Ethereum RPC provider.
    // If set to false, ethereumRPCMaxRequestsPer24HrUTC and
    // ethereumRPCMaxRequestsPerSecond will have no effect.
    enableEthereumRPCRateLimiting?: boolean;
    // A cap on the number of Ethereum JSON-RPC requests a Mesh node will make
    // per 24hr UTC time window (time window starts and ends at midnight UTC).
    // It defaults to 200k but can be increased well beyond this limit depending
    // on your infrastructure or Ethereum RPC provider.
    ethereumRPCMaxRequestsPer24HrUTC?: number;
    // A cap on the number of Ethereum JSON-RPC requests a Mesh node will make
    // per second. This limits the concurrency of these requests and prevents
    // the Mesh node from getting rate-limited. It defaults to the recommended
    // 30 rps for Infura's free tier, and can be increased to 100 rpc for pro
    // users, and potentially higher on alternative infrastructure.
    ethereumRPCMaxRequestsPerSecond?: number;
    // A set of custom addresses to use for the configured network ID. The
    // contract addresses for most common networks are already included by
    // default, so this is typically only needed for testing on custom networks.
    // The given addresses are added to the default list of addresses for known
    // chains and overriding any contract addresses for known chains is not
    // allowed. The addresses for exchange, devUtils, erc20Proxy, and
    // erc721Proxy are required for each chain. For example:
    //
    //    {
    //        exchange: "0x48bacb9266a570d521063ef5dd96e61686dbe788",
    //        devUtils: "0x38ef19fdf8e8415f18c307ed71967e19aac28ba1",
    //        erc20Proxy: "0x1dc4c1cefef38a777b15aa20260a54e584b16c48",
    //        erc721Proxy: "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"
    //    }
    //
    customContractAddresses?: ContractAddresses;
    // The maximum number of orders that Mesh will keep in storage. As the
    // number of orders in storage grows, Mesh will begin enforcing a limit on
    // maximum expiration time for incoming orders and remove any orders with an
    // expiration time too far in the future. Defaults to 100,000.
    maxOrdersInStorage?: number;
    // A a JSON Schema object which will be used for validating incoming orders.
    // If provided, Mesh will only receive orders from other peers in the
    // network with the same filter.
    //
    // Here is an example filter which will only allow orders with a specific
    // makerAssetData:
    //
    //    {
    //        properties: {
    //            makerAssetData: {
    //                const: "0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"
    //            }
    //        }
    //    }
    //
    // Note that you only need to include the requirements for your specific
    // application in the filter. The default requirements for a valid order (e.g.
    // all the required fields) are automatically included. For more information
    // on JSON Schemas, see https://json-schema.org/
    customOrderFilter?: JsonSchema;
    // Offers the ability to use your own web3 provider for all Ethereum RPC
    // requests instead of the default.
    web3Provider?: SupportedProvider;
    // The maximum number of bytes per second that a peer is allowed to send before
    // failing the bandwidth check. Defaults to 5 MiB.
    maxBytesPerSecond?: number;
}

export interface ContractAddresses {
    exchange: string;
    devUtils: string;
    erc20Proxy: string;
    erc721Proxy: string;
    erc1155Proxy: string;
    weth9?: string;
    zrxToken?: string;
}

export enum Verbosity {
    Panic = 0,
    Fatal = 1,
    Error = 2,
    Warn = 3,
    Info = 4,
    Debug = 5,
    Trace = 6,
}

/**
 * The global entrypoint for creating a new MeshWrapper.
 * @ignore
 */
export interface ZeroExMesh {
    newWrapperAsync(config: WrapperConfig): Promise<MeshWrapper>;
}

/**
 * A direct translation of the MeshWrapper type in Go. Its API exposes only
 * simple JavaScript types like number and string, some of which will be
 * converted. For example, we will convert some strings to BigNumbers.
 * @ignore
 */
export interface MeshWrapper {
    startAsync(): Promise<void>;
    onError(handler: (err: Error) => void): void;
    onOrderEvents(handler: (events: WrapperOrderEvent[]) => void): void;
    getStatsAsync(): Promise<WrapperStats>;
    getOrdersForPageAsync(perPage: number, minOrderHash?: string): Promise<WrapperGetOrdersResponse>;
    addOrdersAsync(orders: WrapperSignedOrder[], pinned: boolean): Promise<WrapperValidationResults>;
}

/**
 * The type for configuration exposed by MeshWrapper.
 * @ignore
 */
export interface WrapperConfig {
    verbosity?: number;
    ethereumRPCURL?: string;
    ethereumChainID: number;
    useBootstrapList?: boolean;
    bootstrapList?: string; // comma-separated string instead of an array of strings.
    blockPollingIntervalSeconds?: number;
    ethereumRPCMaxContentLength?: number;
    ethereumRPCMaxRequestsPer24HrUTC?: number;
    ethereumRPCMaxRequestsPerSecond?: number;
    enableEthereumRPCRateLimiting?: boolean;
    customContractAddresses?: string; // json-encoded string instead of Object.
    maxOrdersInStorage?: number;
    customOrderFilter?: string; // json-encoded string instead of Object
    web3Provider?: ZeroExProvider; // Standardized ZeroExProvider instead the more permissive SupportedProvider interface
    maxBytesPerSecond?: number;
}

/**
 * The type for signed orders exposed by MeshWrapper. Unlike other types, the
 * analog isn't defined here. Instead we re-use the definition in
 * @0x/order-utils.
 * @ignore
 */
export interface WrapperSignedOrder {
    makerAddress: string;
    makerAssetData: string;
    makerAssetAmount: string;
    makerFee: string;
    makerFeeAssetData: string;
    takerAddress: string;
    takerAssetData: string;
    takerFeeAssetData: string;
    takerAssetAmount: string;
    takerFee: string;
    senderAddress: string;
    feeRecipientAddress: string;
    expirationTimeSeconds: string;
    salt: string;
    signature: string;
    exchangeAddress: string;
    chainId: number;
}

export interface ERC20TransferEvent {
    from: string;
    to: string;
    value: BigNumber;
}

/** @ignore */
export interface WrapperERC20TransferEvent {
    from: string;
    to: string;
    value: string;
}

export interface ERC20ApprovalEvent {
    owner: string;
    spender: string;
    value: BigNumber;
}

/** @ignore */
export interface WrapperERC20ApprovalEvent {
    owner: string;
    spender: string;
    value: string;
}

export interface ERC721TransferEvent {
    from: string;
    to: string;
    tokenId: BigNumber;
}

/** @ignore */
export interface WrapperERC721TransferEvent {
    from: string;
    to: string;
    tokenId: string;
}

export interface ERC721ApprovalEvent {
    owner: string;
    approved: string;
    tokenId: BigNumber;
}

/** @ignore */
export interface WrapperERC721ApprovalEvent {
    owner: string;
    approved: string;
    tokenId: string;
}

export interface ERC721ApprovalForAllEvent {
    owner: string;
    operator: string;
    approved: boolean;
}

export interface ERC1155TransferSingleEvent {
    operator: string;
    from: string;
    to: string;
    id: BigNumber;
    value: BigNumber;
}

/** @ignore */
export interface WrapperERC1155TransferSingleEvent {
    operator: string;
    from: string;
    to: string;
    id: string;
    value: string;
}

export interface ERC1155TransferBatchEvent {
    operator: string;
    from: string;
    to: string;
    ids: BigNumber[];
    values: BigNumber[];
}

/** @ignore */
export interface WrapperERC1155TransferBatchEvent {
    operator: string;
    from: string;
    to: string;
    ids: string[];
    values: string[];
}

export interface ERC1155ApprovalForAllEvent {
    owner: string;
    operator: string;
    approved: boolean;
}

export interface ExchangeFillEvent {
    makerAddress: string;
    takerAddress: string;
    senderAddress: string;
    feeRecipientAddress: string;
    makerAssetFilledAmount: BigNumber;
    takerAssetFilledAmount: BigNumber;
    makerFeePaid: BigNumber;
    takerFeePaid: BigNumber;
    protocolFeePaid: BigNumber;
    orderHash: string;
    makerAssetData: string;
    takerAssetData: string;
    makerFeeAssetData: string;
    takerFeeAssetData: string;
}

/** @ignore */
export interface WrapperExchangeFillEvent {
    makerAddress: string;
    takerAddress: string;
    senderAddress: string;
    feeRecipientAddress: string;
    makerAssetFilledAmount: string;
    takerAssetFilledAmount: string;
    makerFeePaid: string;
    takerFeePaid: string;
    protocolFeePaid: string;
    orderHash: string;
    makerAssetData: string;
    takerAssetData: string;
    makerFeeAssetData: string;
    takerFeeAssetData: string;
}

export interface ExchangeCancelEvent {
    makerAddress: string;
    senderAddress: string;
    feeRecipientAddress: string;
    orderHash: string;
    makerAssetData: string;
    takerAssetData: string;
}

export interface ExchangeCancelUpToEvent {
    makerAddress: string;
    orderSenderAddress: string;
    orderEpoch: BigNumber;
}

/** @ignore */
export interface WrapperExchangeCancelUpToEvent {
    makerAddress: string;
    orderSenderAddress: string;
    orderEpoch: string;
}

export interface WethWithdrawalEvent {
    owner: string;
    value: BigNumber;
}

/** @ignore */
export interface WrapperWethWithdrawalEvent {
    owner: string;
    value: string;
}

export interface WethDepositEvent {
    owner: string;
    value: BigNumber;
}

/** @ignore */
export interface WrapperWethDepositEvent {
    owner: string;
    value: string;
}

export enum ContractEventKind {
    ERC20TransferEvent = 'ERC20TransferEvent',
    ERC20ApprovalEvent = 'ERC20ApprovalEvent',
    ERC721TransferEvent = 'ERC721TransferEvent',
    ERC721ApprovalEvent = 'ERC721ApprovalEvent',
    ERC721ApprovalForAllEvent = 'ERC721ApprovalForAllEvent',
    ERC1155ApprovalForAllEvent = 'ERC1155ApprovalForAllEvent',
    ERC1155TransferSingleEvent = 'ERC1155TransferSingleEvent',
    ERC1155TransferBatchEvent = 'ERC1155TransferBatchEvent',
    ExchangeFillEvent = 'ExchangeFillEvent',
    ExchangeCancelEvent = 'ExchangeCancelEvent',
    ExchangeCancelUpToEvent = 'ExchangeCancelUpToEvent',
    WethDepositEvent = 'WethDepositEvent',
    WethWithdrawalEvent = 'WethWithdrawalEvent',
}

/** @ignore */
export type WrapperContractEventParameters =
    | WrapperERC20TransferEvent
    | WrapperERC20ApprovalEvent
    | WrapperERC721TransferEvent
    | WrapperERC721ApprovalEvent
    | WrapperExchangeFillEvent
    | WrapperExchangeCancelUpToEvent
    | WrapperWethWithdrawalEvent
    | WrapperWethDepositEvent
    | ERC721ApprovalForAllEvent
    | ExchangeCancelEvent
    | WrapperERC1155TransferSingleEvent
    | WrapperERC1155TransferBatchEvent
    | ERC1155ApprovalForAllEvent;

/** @ignore */
export type ContractEventParameters =
    | ERC20TransferEvent
    | ERC20ApprovalEvent
    | ERC721TransferEvent
    | ERC721ApprovalEvent
    | ExchangeFillEvent
    | ExchangeCancelUpToEvent
    | WethWithdrawalEvent
    | WethDepositEvent
    | ERC721ApprovalForAllEvent
    | ExchangeCancelEvent
    | ERC1155TransferSingleEvent
    | ERC1155TransferBatchEvent
    | ERC1155ApprovalForAllEvent;

export interface ContractEvent {
    blockHash: string;
    txHash: string;
    txIndex: number;
    logIndex: number;
    isRemoved: boolean;
    address: string;
    kind: ContractEventKind;
    parameters: ContractEventParameters;
}

/**
 * The type for order events exposed by MeshWrapper.
 * @ignore
 */
export interface WrapperContractEvent {
    blockHash: string;
    txHash: string;
    txIndex: number;
    logIndex: number;
    isRemoved: boolean;
    address: string;
    kind: string;
    parameters: WrapperContractEventParameters;
}

export enum OrderEventEndState {
    Invalid = 'INVALID',
    Added = 'ADDED',
    Filled = 'FILLED',
    FullyFilled = 'FULLY_FILLED',
    Cancelled = 'CANCELLED',
    Expired = 'EXPIRED',
    Unexpired = 'UNEXPIRED',
    Unfunded = 'UNFUNDED',
    FillabilityIncreased = 'FILLABILITY_INCREASED',
    StoppedWatching = 'STOPPED_WATCHING',
}

/** @ignore */
export interface WrapperOrderEvent {
    timestamp: string;
    orderHash: string;
    signedOrder: WrapperSignedOrder;
    endState: OrderEventEndState;
    fillableTakerAssetAmount: string;
    contractEvents: WrapperContractEvent[];
}

/**
 * Order events are fired by Mesh whenever an order is added, canceled, expired,
 * or filled.
 */
export interface OrderEvent {
    timestampMs: number;
    orderHash: string;
    signedOrder: SignedOrder;
    endState: OrderEventEndState;
    fillableTakerAssetAmount: BigNumber;
    contractEvents: ContractEvent[];
}

/** @ignore */
export interface WrapperValidationResults {
    accepted: WrapperAcceptedOrderInfo[];
    rejected: WrapperRejectedOrderInfo[];
}

/** @ignore */
export interface WrapperAcceptedOrderInfo {
    orderHash: string;
    signedOrder: WrapperSignedOrder;
    fillableTakerAssetAmount: string;
    isNew: boolean;
}

/** @ignore */
export interface WrapperRejectedOrderInfo {
    orderHash: string;
    signedOrder: WrapperSignedOrder;
    kind: RejectedOrderKind;
    status: RejectedOrderStatus;
}

/**
 * Indicates which orders where accepted, which were rejected, and why.
 */
export interface ValidationResults {
    accepted: AcceptedOrderInfo[];
    rejected: RejectedOrderInfo[];
}

/**
 * Info for any orders that were accepted.
 */
export interface AcceptedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: BigNumber;
    isNew: boolean;
}

/**
 * Info for any orders that were rejected, including the reason they were
 * rejected.
 */
export interface RejectedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: RejectedOrderKind;
    status: RejectedOrderStatus;
}

/**
 * A set of categories for rejected orders.
 */
export enum RejectedOrderKind {
    ZeroExValidation = 'ZEROEX_VALIDATION',
    MeshError = 'MESH_ERROR',
    MeshValidation = 'MESH_VALIDATION',
}

/**
 * Provides more information about why an order was rejected.
 */
export interface RejectedOrderStatus {
    code: string;
    message: string;
}

/** @ignore */
export interface WrapperLatestBlock {
    number: string;
    hash: string;
}

export interface LatestBlock {
    number: BigNumber;
    hash: string;
}

/** @ignore */
export interface WrapperStats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock?: WrapperLatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: string; // string instead of BigNumber
    startOfCurrentUTCDay: string; // string instead of Date
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}

export interface Stats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock?: LatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: BigNumber;
    startOfCurrentUTCDay: Date;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}
// tslint:disable-next-line:max-file-line-count
