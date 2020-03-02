export * from './mesh';

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
    return new Promise((resolve, reject) => {
        const go = new Go();
        WebAssembly.instantiateStreaming(response, go.importObject)
            .then(module => {
                resolve(go.run(module.instance));
            })
            .catch(reject);
    });
}
