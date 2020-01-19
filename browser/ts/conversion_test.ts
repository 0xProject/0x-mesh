import { logUtils } from '@0x/utils';
import * as BrowserFS from 'browserfs';

import {
    configToWrapperConfig,
    orderEventsHandlerToWrapperOrderEventsHandler,
    signedOrderToWrapperSignedOrder,
    wrapperAcceptedOrderInfoToAcceptedOrderInfo,
    wrapperContractEventsToContractEvents,
    wrapperOrderEventToOrderEvent,
    wrapperRejectedOrderInfoToRejectedOrderInfo,
    wrapperSignedOrderToSignedOrder,
    wrapperValidationResultsToValidationResults,
} from '../ts/encoding';
import { wasmBuffer } from '../ts/generated/test_wasm_buffer';
import '../ts/wasm_exec';

interface TestCase {}

// The Go code sets certain global values and this is our only way of
// interacting with it. Define those values and their types here.
declare global {
    // Defined in wasm_exec.ts
    class Go {
        public importObject: any;
        public run(instance: WebAssembly.Instance): void;
    }

    // Define variables that are defined in `browser/go/conversion-test.go`
    const testCases: TestCase;
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
const loadEventName = '0xmeshtest';
window.addEventListener(loadEventName, () => {
    isWasmLoaded = true;
});

// Start compiling the WebAssembly as soon as the script is loaded. This lets
// us initialize as quickly as possible.
const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject)
    .then(module => {
        go.run(module.instance);
    })
    .catch(err => {
        // tslint:disable-next-line no-console
        console.error('Could not load Wasm');
        // tslint:disable-next-line no-console
        console.error(err);
        // If the Wasm bytecode didn't compile, Mesh won't work. We have no
        // choice but to throw an error.
        setImmediate(() => {
            throw err;
        });
    });

/*********************** Tests ***********************/
// tslint:disable:no-console

(async () => {
    await waitForLoadAsync();

    if (!testCases) {
        console.log('test cases not initialized');
        console.error('test cases not initialized');
    }

    console.log('test cases initialized');
})();

// tslint:enable:no-console
/*********************** Utils ***********************/

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
