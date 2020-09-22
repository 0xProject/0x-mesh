import { SignedOrder } from '@0x/order-utils';

import { createDatabase } from './database';
import { createSchemaValidator } from './schema_validator';
import {
    AcceptedOrderInfo,
    BigNumber,
    Config,
    ContractAddresses,
    ContractEvent,
    ERC1155ApprovalForAllEvent,
    ERC1155TransferBatchEvent,
    ERC1155TransferSingleEvent,
    ERC20ApprovalEvent,
    ERC20TransferEvent,
    ERC721ApprovalEvent,
    ERC721ApprovalForAllEvent,
    ERC721TransferEvent,
    ExchangeCancelEvent,
    ExchangeCancelUpToEvent,
    ExchangeFillEvent,
    GetOrdersResponse,
    JsonSchema,
    LatestBlock,
    MeshWrapper,
    OrderEvent,
    OrderEventEndState,
    OrderInfo,
    RejectedOrderInfo,
    RejectedOrderKind,
    RejectedOrderStatus,
    Stats,
    SupportedProvider,
    ValidationResults,
    Verbosity,
    WethDepositEvent,
    WethWithdrawalEvent,
    WrapperOrderEvent,
    ZeroExMesh,
} from './types';
import './wasm_exec';
import {
    configToWrapperConfig,
    orderEventsHandlerToWrapperOrderEventsHandler,
    signedOrderToWrapperSignedOrder,
    wrapperGetOrdersResponseToGetOrdersResponse,
    wrapperStatsToStats,
    wrapperValidationResultsToValidationResults,
} from './wrapper_conversion';

export {
    AcceptedOrderInfo,
    BigNumber,
    Config,
    ContractAddresses,
    ContractEvent,
    ERC1155ApprovalForAllEvent,
    ERC1155TransferSingleEvent,
    ERC1155TransferBatchEvent,
    ERC20ApprovalEvent,
    ERC20TransferEvent,
    ERC721ApprovalEvent,
    ERC721ApprovalForAllEvent,
    ERC721TransferEvent,
    ExchangeCancelEvent,
    ExchangeCancelUpToEvent,
    ExchangeFillEvent,
    GetOrdersResponse,
    LatestBlock,
    JsonSchema,
    OrderEvent,
    OrderEventEndState,
    OrderInfo,
    RejectedOrderInfo,
    RejectedOrderKind,
    RejectedOrderStatus,
    SignedOrder,
    SupportedProvider,
    Stats,
    ValidationResults,
    Verbosity,
    WethDepositEvent,
    WethWithdrawalEvent,
};

// The Go code sets certain global values and this is our only way of
// interacting with it. Define those values and their types here.
declare global {
    // Defined in wasm_exec.ts
    class Go {
        public importObject: any;
        public run(instance: WebAssembly.Instance): void;
    }

    // Defined in ../go/main.go
    const zeroExMesh: ZeroExMesh;
}

/**
 * Sets the required global variables the Mesh needs to access from Go land.
 * This includes the `db` and `orderfilter` packages.
 *
 * @ignore
 */
export function _setGlobals(): void {
    (window as any).__mesh_createSchemaValidator__ = createSchemaValidator;
    (window as any).__mesh_dexie_newDatabase__ = createDatabase;
}

// We immediately want to set the required globals.
_setGlobals();

// The interval (in milliseconds) to check whether Wasm is done loading.
const wasmLoadCheckIntervalMs = 100;

// We use a global variable to track whether the Wasm code has finished loading.
let isWasmLoaded = false;
const loadEventName = '0xmeshload';
window.addEventListener(loadEventName, () => {
    isWasmLoaded = true;
});

/**
 * The main class for this package. Has methods for receiving order events and
 * sending orders through the 0x Mesh network.
 */
// tslint:disable-next-line max-classes-per-file
export class Mesh {
    public wrapper?: MeshWrapper;
    private readonly _config: Config;
    private _errHandler?: (err: Error) => void;
    private _orderEventsHandler?: (events: WrapperOrderEvent[]) => void;

    /**
     * Instantiates a new Mesh instance.
     *
     * @param   config               Configuration options for Mesh
     * @return  An instance of Mesh
     */
    constructor(config: Config) {
        this._config = config;
    }

    /**
     * Registers a handler which will be called in the event of a critical
     * error. Note that the handler will not be called for non-critical errors.
     * In order to ensure no errors are missed, this should be called before
     * startAsync.
     *
     * @param   handler               The handler to be called.
     */
    public onError(handler: (err: Error) => void): void {
        this._errHandler = handler;
        this.wrapper?.onError(this._errHandler);
    }

    /**
     * Registers a handler which will be called for any incoming order events.
     * Order events are fired whenver an order is added, canceled, expired, or
     * filled. In order to ensure no events are missed, this should be called
     * before startAsync.
     *
     * @param   handler                The handler to be called.
     */
    public onOrderEvents(handler: (events: OrderEvent[]) => void): void {
        this._orderEventsHandler = orderEventsHandlerToWrapperOrderEventsHandler(handler);
        this.wrapper?.onOrderEvents(this._orderEventsHandler);
    }

    /**
     * Starts the Mesh node in the background. Mesh will automatically find
     * peers in the network and begin receiving orders from them.
     */
    public async startAsync(): Promise<void> {
        await waitForLoadAsync();
        this.wrapper = await zeroExMesh.newWrapperAsync(configToWrapperConfig(this._config));
        if (this._orderEventsHandler !== undefined) {
            this.wrapper.onOrderEvents(this._orderEventsHandler);
        }
        if (this._errHandler !== undefined) {
            this.wrapper.onError(this._errHandler);
        }
        return this.wrapper.startAsync();
    }

    /**
     * Returns various stats about Mesh, including the total number of orders
     * and the number of peers Mesh is connected to.
     */
    public async getStatsAsync(): Promise<Stats> {
        await waitForLoadAsync();
        if (this.wrapper === undefined) {
            // If this is called after startAsync, this.wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }
        const wrapperStats = await this.wrapper.getStatsAsync();
        return wrapperStatsToStats(wrapperStats);
    }

    /**
     * Get all 0x signed orders currently stored in the Mesh node
     * @param perPage number of signedOrders to fetch per paginated request
     * @returns the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts
     */
    public async getOrdersAsync(perPage: number = 200): Promise<GetOrdersResponse> {
        await waitForLoadAsync();
        if (this.wrapper === undefined) {
            // If this is called after startAsync, this.wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }

        let getOrdersResponse = await this.getOrdersForPageAsync(perPage);
        let ordersInfos = getOrdersResponse.ordersInfos;
        let allOrderInfos: OrderInfo[] = [];

        while (ordersInfos.length > 0) {
            allOrderInfos = [...allOrderInfos, ...ordersInfos];
            const minOrderHash = ordersInfos[ordersInfos.length - 1].orderHash;
            getOrdersResponse = await this.getOrdersForPageAsync(perPage, minOrderHash);
            ordersInfos = getOrdersResponse.ordersInfos;
        }

        getOrdersResponse = {
            timestamp: getOrdersResponse.timestamp,
            ordersInfos: allOrderInfos,
        };
        return getOrdersResponse;
    }

    /**
     * Get page of 0x signed orders stored on the Mesh node at the specified snapshot
     * @param perPage Number of signedOrders to fetch per paginated request
     * @param minOrderHash The minimum order hash for the returned orders. Should be set based on the last hash from the previous response.
     * @returns Up to perPage orders with hash greater than minOrderHash, including order hashes and fillableTakerAssetAmounts
     */
    public async getOrdersForPageAsync(perPage: number, minOrderHash?: string): Promise<GetOrdersResponse> {
        await waitForLoadAsync();
        if (this.wrapper === undefined) {
            // If this is called after startAsync, this.wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }

        const wrapperOrderResponse = await this.wrapper.getOrdersForPageAsync(perPage, minOrderHash);
        return wrapperGetOrdersResponseToGetOrdersResponse(wrapperOrderResponse);
    }

    /**
     * Validates and adds the given orders to Mesh. If an order is successfully
     * added, Mesh will share it with any peers in the network and start
     * watching it for changes (e.g. filled, canceled, expired). The returned
     * promise will only be rejected if there was an error validating or adding
     * the order; it will not be rejected for any invalid orders (check
     * results.rejected instead).
     *
     * @param   orders      An array of orders to add.
     * @param   pinned      Whether or not the orders should be pinned. Pinned
     * orders will not be affected by any DDoS prevention or incentive
     * mechanisms and will always stay in storage until they are no longer
     * fillable.
     * @returns Validation results for the given orders, indicating which orders
     * were accepted and which were rejected.
     */
    public async addOrdersAsync(orders: SignedOrder[], pinned: boolean = true): Promise<ValidationResults> {
        await waitForLoadAsync();
        if (this.wrapper === undefined) {
            // If this is called after startAsync, this.wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }
        const meshOrders = orders.map(signedOrderToWrapperSignedOrder);
        const meshResults = await this.wrapper.addOrdersAsync(meshOrders, pinned);
        return wrapperValidationResultsToValidationResults(meshResults);
    }
}

async function waitForLoadAsync(): Promise<void> {
    // Note: this approach is not CPU efficient but it avoids race
    // conditions and has the advantage of returning instantaneously if the
    // Wasm code has already loaded.
    while (!isWasmLoaded) {
        await sleepAsync(wasmLoadCheckIntervalMs);
    }
}

async function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
// tslint:disable-next-line:max-file-line-count
