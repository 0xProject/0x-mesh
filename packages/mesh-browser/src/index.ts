export * from '@0x/mesh-browser-lite/lib/mesh';

import { wasmBuffer } from './generated/wasm_buffer';

// Start compiling the WebAssembly as soon as the script is loaded. This lets
// us initialize as quickly as possible.
const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject)
    .then((module) => {
        go.run(module.instance);
    })
    .catch((err) => {
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
