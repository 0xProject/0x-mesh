import { wasmBuffer } from './generated/wasm_buffer';

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
    setOrderEventsHandler(handler: (events: Array<OrderEvent>) => void): void;
    addOrdersAsync(orders: Array<SignedOrder>): Promise<ValidationResults>;
}

export interface OrderEvent {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: string;
    fillableTakerAssetAmount: string;
    txHashes: Array<string>;
}

// TODO(albrow): Use existing 0x types where possible instead of creating new
// ones.
export interface SignedOrder {
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

export interface ValidationResults {
    accepted: Array<AcceptedOrderInfo>;
    rejected: Array<RejectedOrderInfo>;
}

export interface AcceptedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: string;
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
    private _orderEventHandler?: (events: Array<OrderEvent>) => void;

    constructor(config: Config) {
        this._config = config;
    }

    private async _waitForLoadAsync(): Promise<void> {
        while (!isWasmLoaded) {
            await sleepAsync(100);
        }
    }

    setOrderEventsHandler(handler: (events: Array<OrderEvent>) => void) {
        this._orderEventHandler = handler;
        if (this._wrapper != undefined) {
            this._wrapper.setOrderEventsHandler(this._orderEventHandler);
        }
    }

    async startAsync(): Promise<void> {
        await this._waitForLoadAsync();
        this._wrapper = await zeroExMesh.newWrapperAsync(this._config);
        if (this._orderEventHandler != undefined) {
            this._wrapper.setOrderEventsHandler(this._orderEventHandler);
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
        return this._wrapper.addOrdersAsync(orders);
    }
}

function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
