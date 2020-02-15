export * from './mesh';

/**
 * Loads the Wasm module that is provided by fetching a url.
 * @param url The URL to query for the Wasm binary.
 */
export function loadMeshStreamingWithURL(url: string): void {
    loadMeshStreaming(fetch(url));
}

/**
 * Loads the Wasm module that is provided by a response.
 * @param response The Wasm response that supplies the Wasm binary.
 */
export function loadMeshStreaming(response: Response | Promise<Response>): void {
    const go = new Go();
    WebAssembly.instantiateStreaming(response, go.importObject)
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
}
