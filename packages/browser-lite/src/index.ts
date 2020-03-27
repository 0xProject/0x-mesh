export * from './mesh';

// If needed, add a polyfill for instantiateStreaming
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp: any, importObject: any) => {
        const source = await (await resp).arrayBuffer();
        return WebAssembly.instantiate(source, importObject);
    };
}

/**
 * Loads the Wasm module that is provided by fetching a url.
 * @param url The URL to query for the Wasm binary.
 */
export async function loadMeshStreamingWithURLAsync(url: string): Promise<void> {
    return loadMeshStreamingAsync(fetch(url));
}

/**
 * Loads the Wasm module that is provided by a response.
 * @param response The Wasm response that supplies the Wasm binary.
 */
export async function loadMeshStreamingAsync(response: Response | Promise<Response>): Promise<void> {
    const go = new Go();

    const module = await WebAssembly.instantiateStreaming(response, go.importObject);
    // NOTE(jalextowle): Wrapping the `go.run(module.instance)` statement in `setImmediate`
    // prevents the statement from blocking when `await` is used with this load function.
    setImmediate(() => {
        go.run(module.instance);
    });
}
