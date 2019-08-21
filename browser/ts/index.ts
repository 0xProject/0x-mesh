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
}

var isWasmLoaded = false;
const loadEventName = '0xmeshload';
window.addEventListener(loadEventName, () => {
    console.log('WASM is done loading. Mesh is ready to go :D');
    isWasmLoaded = true;
});

const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject)
    .then(module => {
        go.run(module.instance);
    })
    .catch(err => {
        console.error('Could not load WASM');
        console.error(err);
    });

export class Mesh {
    private _config: Config;
    private _wrapper?: MeshWrapper;

    constructor(config: Config) {
        this._config = config;
    }

    async waitForLoadAsync(): Promise<void> {
        while (!isWasmLoaded) {
            await sleepAsync(100);
        }
    }

    async startAsync(): Promise<void> {
        await this.waitForLoadAsync();
        this._wrapper = await zeroExMesh.newWrapperAsync(this._config);
        return this._wrapper.startAsync();
    }
}

function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
