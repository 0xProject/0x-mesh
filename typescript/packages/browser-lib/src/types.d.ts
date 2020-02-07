import { SignedOrder } from '@0x/order-utils';
import { BigNumber } from '@0x/utils';
import { SupportedProvider, ZeroExProvider } from 'ethereum-types';
export { SignedOrder } from '@0x/order-utils';
export { BigNumber } from '@0x/utils';
export { SupportedProvider } from 'ethereum-types';
export interface WrapperGetOrdersResponse {
    snapshotID: string;
    snapshotTimestamp: string;
    ordersInfos: WrapperOrderInfo[];
}
export interface GetOrdersResponse {
    snapshotID: string;
    snapshotTimestamp: number;
    ordersInfos: OrderInfo[];
}
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
    const?: any;
    type?: string | string[];
    format?: string;
    allOf?: JsonSchema[];
    anyOf?: JsonSchema[];
    oneOf?: JsonSchema[];
    not?: JsonSchema;
}
/**
 * A set of configuration options for Mesh.
 */
export interface Config {
    verbosity?: Verbosity;
    ethereumRPCURL?: string;
    ethereumChainID: number;
    useBootstrapList?: boolean;
    bootstrapList?: string[];
    blockPollingIntervalSeconds?: number;
    ethereumRPCMaxContentLength?: number;
    enableEthereumRPCRateLimiting?: boolean;
    ethereumRPCMaxRequestsPer24HrUTC?: number;
    ethereumRPCMaxRequestsPerSecond?: number;
    customContractAddresses?: ContractAddresses;
    maxOrdersInStorage?: number;
    customOrderFilter?: JsonSchema;
    web3Provider?: SupportedProvider;
}
export interface ContractAddresses {
    exchange: string;
    devUtils: string;
    erc20Proxy: string;
    erc721Proxy: string;
    erc1155Proxy: string;
    coordinator?: string;
    coordinatorRegistry?: string;
    weth9?: string;
    zrxToken?: string;
}
export declare enum Verbosity {
    Panic = 0,
    Fatal = 1,
    Error = 2,
    Warn = 3,
    Info = 4,
    Debug = 5,
    Trace = 6
}
export interface ZeroExMesh {
    newWrapperAsync(config: WrapperConfig): Promise<MeshWrapper>;
}
export interface MeshWrapper {
    startAsync(): Promise<void>;
    onError(handler: (err: Error) => void): void;
    onOrderEvents(handler: (events: WrapperOrderEvent[]) => void): void;
    getStatsAsync(): Promise<WrapperStats>;
    getOrdersForPageAsync(page: number, perPage: number, snapshotID?: string): Promise<WrapperGetOrdersResponse>;
    addOrdersAsync(orders: WrapperSignedOrder[], pinned: boolean): Promise<WrapperValidationResults>;
}
export interface WrapperConfig {
    verbosity?: number;
    ethereumRPCURL?: string;
    ethereumChainID: number;
    useBootstrapList?: boolean;
    bootstrapList?: string;
    blockPollingIntervalSeconds?: number;
    ethereumRPCMaxContentLength?: number;
    ethereumRPCMaxRequestsPer24HrUTC?: number;
    ethereumRPCMaxRequestsPerSecond?: number;
    enableEthereumRPCRateLimiting?: boolean;
    customContractAddresses?: string;
    maxOrdersInStorage?: number;
    customOrderFilter?: string;
    web3Provider?: ZeroExProvider;
}
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
export interface WrapperExchangeCancelUpToEvent {
    makerAddress: string;
    orderSenderAddress: string;
    orderEpoch: string;
}
export interface WethWithdrawalEvent {
    owner: string;
    value: BigNumber;
}
export interface WrapperWethWithdrawalEvent {
    owner: string;
    value: string;
}
export interface WethDepositEvent {
    owner: string;
    value: BigNumber;
}
export interface WrapperWethDepositEvent {
    owner: string;
    value: string;
}
export declare enum ContractEventKind {
    ERC20TransferEvent = "ERC20TransferEvent",
    ERC20ApprovalEvent = "ERC20ApprovalEvent",
    ERC721TransferEvent = "ERC721TransferEvent",
    ERC721ApprovalEvent = "ERC721ApprovalEvent",
    ERC721ApprovalForAllEvent = "ERC721ApprovalForAllEvent",
    ERC1155ApprovalForAllEvent = "ERC1155ApprovalForAllEvent",
    ERC1155TransferSingleEvent = "ERC1155TransferSingleEvent",
    ERC1155TransferBatchEvent = "ERC1155TransferBatchEvent",
    ExchangeFillEvent = "ExchangeFillEvent",
    ExchangeCancelEvent = "ExchangeCancelEvent",
    ExchangeCancelUpToEvent = "ExchangeCancelUpToEvent",
    WethDepositEvent = "WethDepositEvent",
    WethWithdrawalEvent = "WethWithdrawalEvent"
}
export declare type WrapperContractEventParameters = WrapperERC20TransferEvent | WrapperERC20ApprovalEvent | WrapperERC721TransferEvent | WrapperERC721ApprovalEvent | WrapperExchangeFillEvent | WrapperExchangeCancelUpToEvent | WrapperWethWithdrawalEvent | WrapperWethDepositEvent | ERC721ApprovalForAllEvent | ExchangeCancelEvent | WrapperERC1155TransferSingleEvent | WrapperERC1155TransferBatchEvent | ERC1155ApprovalForAllEvent;
export declare type ContractEventParameters = ERC20TransferEvent | ERC20ApprovalEvent | ERC721TransferEvent | ERC721ApprovalEvent | ExchangeFillEvent | ExchangeCancelUpToEvent | WethWithdrawalEvent | WethDepositEvent | ERC721ApprovalForAllEvent | ExchangeCancelEvent | ERC1155TransferSingleEvent | ERC1155TransferBatchEvent | ERC1155ApprovalForAllEvent;
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
export declare enum OrderEventEndState {
    Invalid = "INVALID",
    Added = "ADDED",
    Filled = "FILLED",
    FullyFilled = "FULLY_FILLED",
    Cancelled = "CANCELLED",
    Expired = "EXPIRED",
    Unexpired = "UNEXPIRED",
    Unfunded = "UNFUNDED",
    FillabilityIncreased = "FILLABILITY_INCREASED",
    StoppedWatching = "STOPPED_WATCHING"
}
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
export interface WrapperValidationResults {
    accepted: WrapperAcceptedOrderInfo[];
    rejected: WrapperRejectedOrderInfo[];
}
export interface WrapperAcceptedOrderInfo {
    orderHash: string;
    signedOrder: WrapperSignedOrder;
    fillableTakerAssetAmount: string;
    isNew: boolean;
}
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
export declare enum RejectedOrderKind {
    ZeroExValidation = "ZEROEX_VALIDATION",
    MeshError = "MESH_ERROR",
    MeshValidation = "MESH_VALIDATION",
    CoordinatorError = "COORDINATOR_ERROR"
}
/**
 * Provides more information about why an order was rejected.
 */
export interface RejectedOrderStatus {
    code: string;
    message: string;
}
export interface LatestBlock {
    number: number;
    hash: string;
}
export interface WrapperStats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    peerID: string;
    ethereumChainID: number;
    latestBlock: LatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: string;
    startOfCurrentUTCDay: string;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}
export interface Stats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
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
//# sourceMappingURL=types.d.ts.map