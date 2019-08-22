import { SignedOrder } from '@0x/order-utils';
import { wasmBuffer } from './generated/wasm_buffer';
import { BigNumber } from '@0x/utils';

export { SignedOrder } from '@0x/order-utils';
export { BigNumber } from '@0x/utils';

// Side-effect import
// Side-effects include adding an `fs` and `Go` property on the global object.
import './wasm_exec';

// The Go code sets certain global values and this is our only way of
// interacting with it. Define those values and their types here.
declare global {
    // Defined in wasm_exec.ts
    class Go {
        run(instance: WebAssembly.Instance): void;
        importObject: any;
    }

    // Defined in ../go/main.go
    const zeroExMesh: ZeroExMesh;
}

export interface Config {
    ethereumRPCURL: string;
    ethereumNetworkID: number;
    useBootstrapList?: boolean;
    orderExpirationBufferSeconds?: number;
    blockPollingIntervalSeconds?: number;
    ethereumRPCMaxContentLength?: number;
}

interface ZeroExMesh {
    newWrapperAsync(config: Config): Promise<MeshWrapper>;
}

interface MeshWrapper {
    startAsync(): Promise<void>;
    onError(handler: (err: Error) => void): void;
    onOrderEvents(handler: (events: Array<MeshOrderEvent>) => void): void;
    addOrdersAsync(orders: Array<MeshSignedOrder>): Promise<MeshValidationResults>;
}

interface MeshSignedOrder {
    makerAddress: string;
    makerAssetData: string;
    makerAssetAmount: string;
    makerFee: string;
    takerAddress: string;
    takerAssetData: string;
    takerAssetAmount: string;
    takerFee: string;
    senderAddress: string;
    exchangeAddress: string;
    feeRecipientAddress: string;
    expirationTimeSeconds: string;
    salt: string;
    signature: string;
}

interface MeshOrderEvent {
    orderHash: string;
    signedOrder: MeshSignedOrder;
    kind: string;
    fillableTakerAssetAmount: string;
    txHashes: Array<string>;
}

export interface OrderEvent {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: string;
    fillableTakerAssetAmount: BigNumber;
    txHashes: Array<string>;
}

interface MeshValidationResults {
    accepted: Array<MeshAcceptedOrderInfo>;
    rejected: Array<MeshRejectedOrderInfo>;
}

interface MeshAcceptedOrderInfo {
    orderHash: string;
    signedOrder: MeshSignedOrder;
    fillableTakerAssetAmount: string;
    isNew: boolean;
}

interface MeshRejectedOrderInfo {
    orderHash: string;
    signedOrder: MeshSignedOrder;
    kind: RejectedOrderKind;
    status: RejectedOrderStatus;
}

export interface ValidationResults {
    accepted: Array<AcceptedOrderInfo>;
    rejected: Array<RejectedOrderInfo>;
}

export interface AcceptedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: BigNumber;
    isNew: boolean;
}

export interface RejectedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: RejectedOrderKind;
    status: RejectedOrderStatus;
}

export enum RejectedOrderKind {
    ZeroExValidation = 'ZEROEX_VALIDATION',
    MeshError = 'MESH_ERROR',
    MeshValidation = 'MESH_VALIDATION',
    CoordinatorError = 'COORDINATOR_ERROR',
}

export interface RejectedOrderStatus {
    code: string;
    message: string;
}

var isWasmLoaded = false;
const loadEventName = '0xmeshload';
window.addEventListener(loadEventName, () => {
    console.log('Wasm is done loading. Mesh is ready to go :D');
    isWasmLoaded = true;
});

const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject)
    .then(module => {
        go.run(module.instance);
    })
    .catch(err => {
        console.error('Could not load Wasm');
        console.error(err);
    });

export class Mesh {
    private _config: Config;
    private _wrapper?: MeshWrapper;
    private _errHandler?: (err: Error) => void;
    private _orderEventsHandler?: (events: Array<MeshOrderEvent>) => void;

    constructor(config: Config) {
        this._config = config;
    }

    private async _waitForLoadAsync(): Promise<void> {
        while (!isWasmLoaded) {
            await sleepAsync(100);
        }
    }

    onError(handler: (err: Error) => void) {
        this._errHandler = handler;
        if (this._wrapper != undefined) {
            this._wrapper.onError(this._errHandler);
        }
    }

    onOrderEvents(handler: (events: Array<OrderEvent>) => void) {
        this._orderEventsHandler = orderEventsHandlerToMeshOrderEventsHandler(handler);
        if (this._wrapper != undefined) {
            this._wrapper.onOrderEvents(this._orderEventsHandler);
        }
    }

    async startAsync(): Promise<void> {
        await this._waitForLoadAsync();
        this._wrapper = await zeroExMesh.newWrapperAsync(this._config);
        if (this._orderEventsHandler != undefined) {
            this._wrapper.onOrderEvents(this._orderEventsHandler);
        }
        if (this._errHandler != undefined) {
            this._wrapper.onError(this._errHandler);
        }
        return this._wrapper.startAsync();
    }

    async addOrdersAsync(orders: Array<SignedOrder>): Promise<ValidationResults> {
        await this._waitForLoadAsync();
        if (this._wrapper == undefined) {
            // If this is called after startAsync, this._wrapper is always
            // defined. This check is here just in case and satisfies the
            // compiler.
            return Promise.reject(new Error('Mesh is still loading. Try again soon.'));
        }
        const meshOrders = orders.map(signedOrderToMeshSignedOrder);
        const meshResults = await this._wrapper.addOrdersAsync(meshOrders);
        return meshValidationResultsToValidationResults(meshResults);
    }
}

function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function meshSignedOrderToSignedOrder(meshSignedOrder: MeshSignedOrder): SignedOrder {
    return {
        ...meshSignedOrder,
        makerFee: new BigNumber(meshSignedOrder.makerFee),
        takerFee: new BigNumber(meshSignedOrder.takerFee),
        makerAssetAmount: new BigNumber(meshSignedOrder.makerAssetAmount),
        takerAssetAmount: new BigNumber(meshSignedOrder.takerAssetAmount),
        salt: new BigNumber(meshSignedOrder.salt),
        expirationTimeSeconds: new BigNumber(meshSignedOrder.expirationTimeSeconds),
    };
}

function signedOrderToMeshSignedOrder(signedOrder: SignedOrder): MeshSignedOrder {
    return {
        ...signedOrder,
        makerFee: signedOrder.makerFee.toString(),
        takerFee: signedOrder.takerFee.toString(),
        makerAssetAmount: signedOrder.makerAssetAmount.toString(),
        takerAssetAmount: signedOrder.takerAssetAmount.toString(),
        salt: signedOrder.salt.toString(),
        expirationTimeSeconds: signedOrder.expirationTimeSeconds.toString(),
    };
}

function meshOrderEventToOrderEvent(meshOrderEvent: MeshOrderEvent): OrderEvent {
    return {
        ...meshOrderEvent,
        signedOrder: meshSignedOrderToSignedOrder(meshOrderEvent.signedOrder),
        fillableTakerAssetAmount: new BigNumber(meshOrderEvent.fillableTakerAssetAmount),
    };
}

function orderEventsHandlerToMeshOrderEventsHandler(
    orderEventsHandler: (events: Array<OrderEvent>) => void,
): (events: Array<MeshOrderEvent>) => void {
    return (meshOrderEvents: Array<MeshOrderEvent>) => {
        const orderEvents = meshOrderEvents.map(meshOrderEventToOrderEvent);
        orderEventsHandler(orderEvents);
    };
}

function meshValidationResultsToValidationResults(meshValidationResults: MeshValidationResults): ValidationResults {
    return {
        accepted: meshValidationResults.accepted.map(meshAcceptedOrderInfoToAcceptedOrderInfo),
        rejected: meshValidationResults.rejected.map(meshRejectedOrderInfoToRejectedOrderInfo),
    };
}

function meshAcceptedOrderInfoToAcceptedOrderInfo(meshAcceptedOrderInfo: MeshAcceptedOrderInfo): AcceptedOrderInfo {
    return {
        ...meshAcceptedOrderInfo,
        signedOrder: meshSignedOrderToSignedOrder(meshAcceptedOrderInfo.signedOrder),
        fillableTakerAssetAmount: new BigNumber(meshAcceptedOrderInfo.fillableTakerAssetAmount),
    };
}

function meshRejectedOrderInfoToRejectedOrderInfo(meshRejectedOrderInfo: MeshRejectedOrderInfo): RejectedOrderInfo {
    return {
        ...meshRejectedOrderInfo,
        signedOrder: meshSignedOrderToSignedOrder(meshRejectedOrderInfo.signedOrder),
    };
}
