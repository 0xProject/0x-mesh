import { getContractAddressesForChainOrThrow } from '@0x/contract-addresses';
import { SignedOrder } from '@0x/order-utils';
import * as ajv from 'ajv';
import * as BrowserFS from 'browserfs';

import './wasm_exec';

export { SignedOrder } from '@0x/order-utils';
export { BigNumber } from '@0x/utils';
export { SupportedProvider } from 'ethereum-types';

import {
    AcceptedOrderInfo,
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
    ValidationResults,
    Verbosity,
    WethDepositEvent,
    WethWithdrawalEvent,
    WrapperOrderEvent,
    ZeroExMesh,
} from './types';
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

// We use the global willLoadBrowserFS variable to signal that we are going to
// initialize BrowserFS.
(window as any).willLoadBrowserFS = true;

BrowserFS.configure(
    {
        fs: 'IndexedDB',
        options: {
            storeName: '0x-mesh-db',
        },
    },
    e => {
        if (e) {
            throw e;
        }
        // We use the global browserFS variable as a handle for Go/Wasm code to
        // call into the BrowserFS API. Setting this variable also indicates
        // that BrowserFS has finished loading.
        (window as any).browserFS = BrowserFS.BFSRequire('fs');
    },
);

// The interval (in milliseconds) to check whether Wasm is done loading.
const wasmLoadCheckIntervalMs = 100;

// We use a global variable to track whether the Wasm code has finished loading.
let isWasmLoaded = false;
const loadEventName = '0xmeshload';
window.addEventListener(loadEventName, () => {
    isWasmLoaded = true;
});

const addressSchema = {
    $id: 'http://example.com/address',
    type: 'string',
    pattern: '^0x[0-9a-fA-F]{40}$',
};
const wholeNumberSchema = {
    $id: 'http://example.com/wholeNumber',
    anyOf: [{ type: 'string', pattern: '^\\d+$' }, { type: 'integer' }],
};
const hexSchema = { $id: 'http://example.com/hex', type: 'string', pattern: '^0x(([0-9a-fA-F][0-9a-fA-F])+)?$' };
const orderSchema = {
    $id: 'http://example.com/order',
    properties: {
        makerAddress: { $ref: 'http://example.com/address' },
        takerAddress: { $ref: 'http://example.com/address' },
        makerFee: { $ref: 'http://example.com/wholeNumber' },
        takerFee: { $ref: 'http://example.com/wholeNumber' },
        senderAddress: { $ref: 'http://example.com/address' },
        makerAssetAmount: { $ref: 'http://example.com/wholeNumber' },
        takerAssetAmount: { $ref: 'http://example.com/wholeNumber' },
        makerAssetData: { $ref: 'http://example.com/hex' },
        takerAssetData: { $ref: 'http://example.com/hex' },
        makerFeeAssetData: { $ref: 'http://example.com/hex' },
        takerFeeAssetData: { $ref: 'http://example.com/hex' },
        salt: { $ref: 'http://example.com/wholeNumber' },
        feeRecipientAddress: { $ref: 'http://example.com/address' },
        expirationTimeSeconds: { $ref: 'http://example.com/wholeNumber' },
        exchangeAddress: { $ref: 'http://example.com/exchangeAddress' },
        chainId: { $ref: 'http://example.com/chainId' },
    },
    required: [
        'makerAddress',
        'takerAddress',
        'makerFee',
        'takerFee',
        'senderAddress',
        'makerAssetAmount',
        'takerAssetAmount',
        'makerAssetData',
        'takerAssetData',
        'makerFeeAssetData',
        'takerFeeAssetData',
        'salt',
        'feeRecipientAddress',
        'expirationTimeSeconds',
        'exchangeAddress',
        'chainId',
    ],
    type: 'object',
};
const signedOrderSchema = {
    $id: 'http://example.com/signedOrder',
    allOf: [
        { $ref: 'http://example.com/order' },
        { properties: { signature: { $ref: 'http://example.com/hex' } }, required: ['signature'] },
    ],
};
const rootOrderSchema = {
    $id: 'http://example.com/rootOrder',
    allOf: [{ $ref: 'http://example.com/customOrder' }, { $ref: 'http://example.com/signedOrder' }],
};
const rootOrderMessageSchema = {
    $id: 'http://example.com/rootOrderMessage',
    properties: {
        messageType: { type: 'string', pattern: 'order' },
        order: { $ref: 'http://example.com/rootOrder' },
        topics: { type: 'array', minItems: 1, items: { type: 'string' } },
    },
    required: ['messageType', 'order', 'topics'],
};

(window as any).schemaValidator = {};

const schemaValidator = (window as any).schemaValidator;

/**
 * The main class for this package. Has methods for receiving order events and
 * sending orders through the 0x Mesh network.
 */
// tslint:disable-next-line max-classes-per-file
export class Mesh {
    private readonly _config: Config;
    private _wrapper?: MeshWrapper;
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

        const chainIdSchema = {
            $id: 'http://example.com/chainId',
            const: this._config.ethereumChainID,
        };

        // TODO(jalextowle): These addresses are already fully lowercased in Mesh. Verify
        // that these schemas are not missing anything.
        let exchangeAddressSchema: { $id: string; enum: [string] };
        if (this._config.customContractAddresses && this._config.customContractAddresses.exchange) {
            exchangeAddressSchema = {
                $id: 'http://example.com/exchangeAddress',
                enum: [this._config.customContractAddresses.exchange],
            };
        } else {
            const contractAddresses = getContractAddressesForChainOrThrow(this._config.ethereumChainID);
            exchangeAddressSchema = {
                $id: 'http://example.com/exchangeAddress',
                enum: [contractAddresses.exchange],
            };
        }

        const AJV = new ajv({
            schemas: [
                {
                    ...this._config.customOrderFilter,
                    $id: 'http://example.com/customOrder',
                },
                addressSchema,
                wholeNumberSchema,
                hexSchema,
                chainIdSchema,
                exchangeAddressSchema,
                orderSchema,
                signedOrderSchema,
                rootOrderSchema,
                rootOrderMessageSchema,
            ],
        });
        // tslint:disable:no-non-null-assertion
        const orderValidate = AJV.getSchema('http://example.com/rootOrder');
        if (orderValidate === undefined) {
            throw new Error('Cannot find "/rootOrder" schema in AJV');
        }
        schemaValidator.orderValidator = (input: string) => {
            const result: any = { success: false, errors: [] };
            try {
                result.success = orderValidate(JSON.parse(input));
                if (orderValidate.errors) {
                    result.errors = orderValidate.errors.map(error => error.message);
                }
            } catch (error) {
                result.fatal = JSON.stringify(error);
            }
            return result;
        };

        const messageValidate = AJV.getSchema('http://example.com/rootOrderMessage');
        if (messageValidate === undefined) {
            throw new Error('Cannot find "rootOrderMessage" schema in AJV');
        }
        schemaValidator.messageValidator = (input: string) => {
            const result: any = { success: false, errors: [] };
            try {
                result.success = messageValidate(JSON.parse(input));
                if (messageValidate.errors) {
                    result.errors = messageValidate.errors.map(error => error.message);
                }
            } catch (error) {
                result.fatal = JSON.stringify(error);
            }
            return result;
        };
        // tslint:enable:no-non-null-assertion
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
        if (this._wrapper !== undefined) {
            this._wrapper.onError(this._errHandler);
        }
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
        if (this._wrapper !== undefined) {
            this._wrapper.onOrderEvents(this._orderEventsHandler);
        }
    }

    /**
     * Starts the Mesh node in the background. Mesh will automatically find
     * peers in the network and begin receiving orders from them.
     */
    public async startAsync(): Promise<void> {
        await waitForLoadAsync();
        this._wrapper = await zeroExMesh.newWrapperAsync(configToWrapperConfig(this._config));
        if (this._orderEventsHandler !== undefined) {
            this._wrapper.onOrderEvents(this._orderEventsHandler);
        }
        if (this._errHandler !== undefined) {
            this._wrapper.onError(this._errHandler);
        }
        return this._wrapper.startAsync();
    }

    /**
     * Returns various stats about Mesh, including the total number of orders
     * and the number of peers Mesh is connected to.
     */
    public async getStatsAsync(): Promise<Stats> {
        await waitForLoadAsync();
        if (this._wrapper === undefined) {
            // If this is called after startAsync, this._wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }
        const wrapperStats = await this._wrapper.getStatsAsync();
        return wrapperStatsToStats(wrapperStats);
    }

    /**
     * Get all 0x signed orders currently stored in the Mesh node
     * @param perPage number of signedOrders to fetch per paginated request
     * @returns the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts
     */
    public async getOrdersAsync(perPage: number = 200): Promise<GetOrdersResponse> {
        await waitForLoadAsync();
        if (this._wrapper === undefined) {
            // If this is called after startAsync, this._wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }

        let snapshotID = ''; // New snapshot

        // TODO(albrow): De-dupe this code with the method by the same name
        // in the TypeScript RPC client.
        let page = 0;
        let getOrdersResponse = await this.getOrdersForPageAsync(page, perPage, snapshotID);
        snapshotID = getOrdersResponse.snapshotID;
        let ordersInfos = getOrdersResponse.ordersInfos;

        let allOrderInfos: OrderInfo[] = [];

        do {
            allOrderInfos = [...allOrderInfos, ...ordersInfos];
            page++;
            getOrdersResponse = await this.getOrdersForPageAsync(page, perPage, snapshotID);
            ordersInfos = getOrdersResponse.ordersInfos;
        } while (ordersInfos.length > 0);

        getOrdersResponse = {
            snapshotID,
            snapshotTimestamp: getOrdersResponse.snapshotTimestamp,
            ordersInfos: allOrderInfos,
        };
        return getOrdersResponse;
    }

    /**
     * Get page of 0x signed orders stored on the Mesh node at the specified snapshot
     * @param page Page index at which to retrieve orders
     * @param perPage Number of signedOrders to fetch per paginated request
     * @param snapshotID The DB snapshot at which to fetch orders. If omitted, a new snapshot is created
     * @returns the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts
     */
    public async getOrdersForPageAsync(page: number, perPage: number, snapshotID?: string): Promise<GetOrdersResponse> {
        await waitForLoadAsync();
        if (this._wrapper === undefined) {
            // If this is called after startAsync, this._wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }

        const wrapperOrderResponse = await this._wrapper.getOrdersForPageAsync(page, perPage, snapshotID);
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
        if (this._wrapper === undefined) {
            // If this is called after startAsync, this._wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }
        const meshOrders = orders.map(signedOrderToWrapperSignedOrder);
        const meshResults = await this._wrapper.addOrdersAsync(meshOrders, pinned);
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
