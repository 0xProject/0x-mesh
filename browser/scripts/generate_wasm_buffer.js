const fs = require('fs');
const base64 = require('base64-arraybuffer');

const file = fs.readFileSync('./wasm/main.wasm');

// HACK(albrow): We generate a TypeScript file that contains the WASM output
// encoded as a base64 string. This is the most reliable way to load WASM such
// that users just see a TypeScript/JavaScript package and without relying on a
// third-party server.
let outputContents =
    `import * as base64 from "base64-arraybuffer";\nexport const wasmBuffer = base64.decode("` +
    base64.encode(file.buffer) +
    `");`;

fs.writeFileSync('./ts/generated/wasm_buffer.ts', outputContents);
