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

// The Go code sets certain global values and this is our only way of
// interacting with it. Define those values and their types here.
declare global {
    // Defined in wasm_exec.ts
    class Go {
        public importObject: any;
        public run(instance: WebAssembly.Instance): void;
    }

    // Define variables that are defined in `browser/go/browser-test.go`
    // FIXME
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
