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
}

export interface OrderEvent {
    orderHash: string;
    signedOrder: any;
    kind: string;
    fillableTakerAssetAmount: string;
    txHashes: Array<string>;
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

    async startAsync(): Promise<void> {
        await this._waitForLoadAsync();
        this._wrapper = await zeroExMesh.newWrapperAsync(this._config);
        if (this._orderEventHandler != undefined) {
            this._wrapper.setOrderEventsHandler(this._orderEventHandler);
        }
        return this._wrapper.startAsync();
    }

    setOrderEventsHandler(handler: (events: Array<OrderEvent>) => void) {
        this._orderEventHandler = handler;
        if (this._wrapper != undefined) {
            this._wrapper.setOrderEventsHandler(this._orderEventHandler);
        }
    }
}

function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
