import { SignedOrder } from '@0x/types';
import { BigNumber } from '@0x/utils';

export interface StatsResponse {
    stats: StringifiedStats;
}

export interface AddOrdersResponse {
    addOrders: StringifiedAddOrdersResults;
}

export interface OrderResponse {
    order: StringifiedOrderWithMetadata | null;
}

export interface OrdersResponse {
    orders: StringifiedOrderWithMetadata[];
}

export interface OrderEventResponse {
    orderEvents: StringifiedOrderEvent[];
}

export interface Stats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock: LatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: BigNumber;
    startOfCurrentUTCDay: Date;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}

export interface LatestBlock {
    number: BigNumber;
    hash: string;
}

export interface OrderWithMetadata extends SignedOrder {
    hash: string;
    fillableTakerAssetAmount: BigNumber;
}

export interface AddOrdersResults {
    // The set of orders that were accepted. Accepted orders will be watched and order events will be emitted if
    // their status changes.
    accepted: AcceptedOrderResult[];
    // The set of orders that were rejected, including the reason they were rejected. Rejected orders will not be
    // watched.
    rejected: RejectedOrderResult[];
}

export interface AcceptedOrderResult {
    // The order that was accepted, including metadata.
    order: OrderWithMetadata;
    // Whether or not the order is new. Set to true if this is the first time this Mesh node has accepted the order
    // and false otherwise.
    isNew: boolean;
}

export interface RejectedOrderResult {
    // The hash of the order. May be null if the hash could not be computed.
    hash?: string;
    // The order that was rejected.
    order: SignedOrder;
    // A machine-readable code indicating why the order was rejected. This code is designed to
    // be used by programs and applications and will never change without breaking backwards-compatibility.
    code: RejectedOrderCode;
    // A human-readable message indicating why the order was rejected. This message may change
    // in future releases and is not covered by backwards-compatibility guarantees.
    message: string;
}

export enum RejectedOrderCode {
    EthRpcRequestFailed = 'ETH_RPC_REQUEST_FAILED',
    OrderHasInvalidMakerAssetAmount = 'ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT',
    OrderHasInvalidTakerAssetAmount = 'ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT',
    OrderExpired = 'ORDER_EXPIRED',
    OrderFullyFilled = 'ORDER_FULLY_FILLED',
    OrderCancelled = 'ORDER_CANCELLED',
    OrderUnfunded = 'ORDER_UNFUNDED',
    OrderHasInvalidMakerAssetData = 'ORDER_HAS_INVALID_MAKER_ASSET_DATA',
    OrderHasInvalidMakerFeeAssetData = 'ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA',
    OrderHasInvalidTakerAssetData = 'ORDER_HAS_INVALID_TAKER_ASSET_DATA',
    OrderHasInvalidTakerFeeAssetData = 'ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA',
    OrderHasInvalidSignature = 'ORDER_HAS_INVALID_SIGNATURE',
    OrderMaxExpirationExceeded = 'ORDER_MAX_EXPIRATION_EXCEEDED',
    InternalError = 'INTERNAL_ERROR',
    MaxOrderSizeExceeded = 'MAX_ORDER_SIZE_EXCEEDED',
    OrderAlreadyStoredAndUnfillable = 'ORDER_ALREADY_STORED_AND_UNFILLABLE',
    OrderForIncorrectChain = 'ORDER_FOR_INCORRECT_CHAIN',
    IncorrectExchangeAddress = 'INCORRECT_EXCHANGE_ADDRESS',
    SenderAddressNotAllowed = 'SENDER_ADDRESS_NOT_ALLOWED',
    DatabaseFullOfOrders = 'DATABASE_FULL_OF_ORDERS',
}

export interface OrderEvent {
    timestampMs: number;
    order: OrderWithMetadata;
    endState: OrderEventEndState;
    contractEvents: ContractEvent[];
}

export interface ContractEvent {
    blockHash: string;
    txHash: string;
    txIndex: number;
    logIndex: number;
    isRemoved: boolean;
    address: string;
    // TODO(albrow): Use an enum type for kind?
    kind: string;
    // TODO(albrow): Use a union type for parameters?
    parameters: any;
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

export enum OrderEventEndState {
    // The order was successfully validated and added to the Mesh node. The order is now being watched and any changes to
    // the fillability will result in subsequent order events.
    Added = 'ADDED',
    // The order was filled for a partial amount. The order is still fillable up to the fillableTakerAssetAmount.
    Filled = 'FILLED',
    // The order was fully filled and its remaining fillableTakerAssetAmount is 0. The order is no longer fillable.
    FullyFilled = 'FULLY_FILLED',
    // The order was cancelled and is no longer fillable.
    Cancelled = 'CANCELLED',
    // The order expired and is no longer fillable.
    Expired = 'EXPIRED',
    // The order was previously expired, but due to a block re-org it is no longer considered expired (should be rare).
    Unexpired = 'UNEXPIRED',
    // The order has become unfunded and is no longer fillable. This can happen if the maker makes a transfer or changes their allowance.
    Unfunded = 'UNFUNDED',
    // The fillability of the order has increased. This can happen if a previously processed fill event gets reverted due to a block re-org,
    // or if a maker makes a transfer or changes their allowance.
    FillabilityIncreased = 'FILLABILITY_INCREASED',
    // The order is potentially still valid but was removed for a different reason (e.g.
    // the database is full or the peer that sent the order was misbehaving). The order will no longer be watched
    // and no further events for this order will be emitted. In some cases, the order may be re-added in the
    // future.
    StoppedWatching = 'STOPPED_WATCHING',
}

export type OrderField = Extract<keyof OrderWithMetadata, string>;

export enum SortDirection {
    Asc = 'ASC',
    Desc = 'DESC',
}

export enum FilterKind {
    Equal = 'EQUAL',
    NotEqual = 'NOT_EQUAL',
    Greater = 'GREATER',
    GreaterOrEqual = 'GREATER_OR_EQUAL',
    Less = 'LESS',
    LessOrEqual = 'LESS_OR_EQUAL',
}

export interface OrderSort {
    field: OrderField;
    direction: SortDirection;
}

export interface OrderFilter {
    field: OrderField;
    kind: FilterKind;
    value: OrderWithMetadata[OrderField];
}

export interface OrderQuery {
    filters?: OrderFilter[];
    sort?: OrderSort[];
    limit?: number;
}

export interface StringifiedLatestBlock {
    number: string;
    hash: string;
}

export interface StringifiedStats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock: StringifiedLatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: string;
    startOfCurrentUTCDay: string;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}

export interface StringifiedSignedOrder {
    chainId: string;
    exchangeAddress: string;
    makerAddress: string;
    takerAddress: string;
    feeRecipientAddress: string;
    senderAddress: string;
    makerAssetAmount: string;
    takerAssetAmount: string;
    makerFee: string;
    takerFee: string;
    expirationTimeSeconds: string;
    salt: string;
    makerAssetData: string;
    takerAssetData: string;
    makerFeeAssetData: string;
    takerFeeAssetData: string;
    signature: string;
}

export interface StringifiedOrderWithMetadata extends StringifiedSignedOrder {
    hash: string;
    fillableTakerAssetAmount: string;
}

export interface StringifiedAddOrdersResults {
    accepted: StringifiedAcceptedOrderResult[];
    rejected: StringifiedRejectedOrderResult[];
}

export interface StringifiedAcceptedOrderResult {
    order: StringifiedOrderWithMetadata;
    isNew: boolean;
}

export interface StringifiedRejectedOrderResult {
    hash?: string;
    order: StringifiedSignedOrder;
    code: RejectedOrderCode;
    message: string;
}

export interface StringifiedOrderEvent {
    timestamp: string;
    order: StringifiedOrderWithMetadata;
    endState: OrderEventEndState;
    fillableTakerAssetAmount: BigNumber;
    contractEvents: ContractEvent[];
}

/**
 * Converts StringifiedStats to Stats
 */
export function fromStringifiedStats(stats: StringifiedStats): Stats {
    return {
        ...stats,
        latestBlock: fromStringifiedLatestBlock(stats.latestBlock),
        maxExpirationTime: new BigNumber(stats.maxExpirationTime),
        startOfCurrentUTCDay: new Date(stats.startOfCurrentUTCDay),
    };
}

/**
 * Converts StringifiedLatestBlock to LatestBlock
 */
export function fromStringifiedLatestBlock(latestBlock: StringifiedLatestBlock): LatestBlock {
    return {
        ...latestBlock,
        number: new BigNumber(latestBlock.number),
    };
}

/**
 * Converts SignedOrder to StringifiedSignedOrder
 */
export function toStringifiedSignedOrder(order: SignedOrder): StringifiedSignedOrder {
    return {
        ...order,
        chainId: order.chainId.toString(),
        makerAssetAmount: order.makerAssetAmount.toString(),
        takerAssetAmount: order.takerAssetAmount.toString(),
        makerFee: order.makerFee.toString(),
        takerFee: order.takerFee.toString(),
        expirationTimeSeconds: order.expirationTimeSeconds.toString(),
        salt: order.salt.toString(),
    };
}

/**
 * Converts StringifiedOrderWithMetadata to OrderWithMetadata
 */
export function fromStringifiedOrderWithMetadata(order: StringifiedOrderWithMetadata): OrderWithMetadata {
    return {
        ...order,
        // tslint:disable-next-line: custom-no-magic-numbers
        chainId: Number.parseInt(order.chainId, 10),
        makerAssetAmount: new BigNumber(order.makerAssetAmount),
        takerAssetAmount: new BigNumber(order.takerAssetAmount),
        makerFee: new BigNumber(order.makerFee),
        takerFee: new BigNumber(order.takerFee),
        expirationTimeSeconds: new BigNumber(order.expirationTimeSeconds),
        salt: new BigNumber(order.salt),
        fillableTakerAssetAmount: new BigNumber(order.fillableTakerAssetAmount),
    };
}

/**
 * Converts StringifiedSignedOrder to SignedOrder
 */
export function fromStringifiedSignedOrder(order: StringifiedSignedOrder): SignedOrder {
    return {
        ...order,
        // tslint:disable-next-line: custom-no-magic-numbers
        chainId: Number.parseInt(order.chainId, 10),
        makerAssetAmount: new BigNumber(order.makerAssetAmount),
        takerAssetAmount: new BigNumber(order.takerAssetAmount),
        makerFee: new BigNumber(order.makerFee),
        takerFee: new BigNumber(order.takerFee),
        expirationTimeSeconds: new BigNumber(order.expirationTimeSeconds),
        salt: new BigNumber(order.salt),
    };
}

/**
 * Converts StringifiedAddOrdersResults to AddOrdersResults
 */
export function fromStringifiedAddOrdersResults(results: StringifiedAddOrdersResults): AddOrdersResults {
    return {
        accepted: results.accepted.map(fromStringifiedAcceptedOrderResult),
        rejected: results.rejected.map(fromStringifiedRejectedOrderResult),
    };
}

/**
 * Converts StringifiedAcceptedOrderResult to AcceptedOrderResult
 */
export function fromStringifiedAcceptedOrderResult(
    acceptedResult: StringifiedAcceptedOrderResult,
): AcceptedOrderResult {
    return {
        ...acceptedResult,
        order: fromStringifiedOrderWithMetadata(acceptedResult.order),
    };
}

/**
 * Converts StringifiedRejectedOrderResult to RejectedOrderResult
 */
export function fromStringifiedRejectedOrderResult(
    rejectedResult: StringifiedRejectedOrderResult,
): RejectedOrderResult {
    return {
        ...rejectedResult,
        order: fromStringifiedSignedOrder(rejectedResult.order),
    };
}

/**
 * Converts StringifiedOrderEvent to OrderEvent
 */
export function fromStringifiedOrderEvent(event: StringifiedOrderEvent): OrderEvent {
    return {
        ...event,
        timestampMs: new Date(event.timestamp).getUTCMilliseconds(),
        order: fromStringifiedOrderWithMetadata(event.order),
    };
}

/**
 * Converts filter.value the the appropriate JSON/GraphQL type (e.g. BigNumber gets converted to string).
 */
export function convertFilterValue(filter: OrderFilter): OrderFilter {
    return {
        ...filter,
        value: BigNumber.isBigNumber(filter.value) ? filter.value.toString() : filter.value,
    };
}
