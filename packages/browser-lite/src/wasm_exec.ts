// Copyright 2018 The Go Authors. All rights reserved.
// Modified work copyright 2020 ZeroEx, Inc.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO_LICENSE file.

/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
/* tslint:disable */
(() => {
    // Map multiple JavaScript environments to a single common API,
    // preferring web standards over Node.js API.
    //
    // Environments considered:
    // - Browsers
    // - Node.js
    // - Electron
    // - Parcel

    if (typeof global !== 'undefined') {
        // global already exists
    } else if (typeof window !== 'undefined') {
        (window as any).global = window;
    } else if (typeof self !== 'undefined') {
        (self as any).global = self;
    } else {
        throw new Error('cannot export Go (neither global, window nor self is defined)');
    }

    // NOTE(albrow): Since we know this file is only for inclusion in browser
    // environments, we can skip some Node.js-related logic. This means
    // commenting out anything that involves `require`.

    // if (!(global as any).require && typeof require !== 'undefined') {
    //     (global as any).require = require;
    // }

    // if (!(global as any).fs && (global as any).require) {
    //     (global as any).fs = require('fs');
    // }

    if (!(global as any).fs) {
        let outputBuf = '';
        (global as any).fs = {
            constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused
            writeSync(fd: any, buf: any) {
                outputBuf += decoder.decode(buf);
                const nl = outputBuf.lastIndexOf('\n');
                if (nl != -1) {
                    console.log(outputBuf.substr(0, nl));
                    outputBuf = outputBuf.substr(nl + 1);
                }
                return buf.length;
            },
            write(
                fd: any,
                buf: string | any[],
                offset: number,
                length: any,
                position: null,
                callback: (arg0: null, arg1: any) => void,
            ) {
                if (offset !== 0 || length !== buf.length || position !== null) {
                    throw new Error('not implemented');
                }
                const n = this.writeSync(fd, buf);
                callback(null, n);
            },
            open(path: any, flags: any, mode: any, callback: (arg0: Error) => void) {
                const err: any = new Error('not implemented');
                err.code = 'ENOSYS';
                callback(err);
            },
            read(fd: any, buffer: any, offset: any, length: any, position: any, callback: (arg0: Error) => void) {
                const err: any = new Error('not implemented');
                err.code = 'ENOSYS';
                callback(err);
            },
            fsync(fd: any, callback: (arg0: null) => void) {
                callback(null);
            },
        };
    }

    // if (!(global as any).crypto) {
    //     const nodeCrypto = require('crypto');
    //     (global as any).crypto = {
    //         getRandomValues(b: any) {
    //             nodeCrypto.randomFillSync(b);
    //         },
    //     };
    // }

    if (!(global as any).performance) {
        (global as any).performance = {
            now() {
                const [sec, nsec] = process.hrtime();
                return sec * 1000 + nsec / 1000000;
            },
        };
    }

    // if (!(global as any).TextEncoder) {
    //     (global as any).TextEncoder = require('util').TextEncoder;
    // }

    // if (!(global as any).TextDecoder) {
    //     (global as any).TextDecoder = require('util').TextDecoder;
    // }

    // End of polyfills for common API.

    const encoder = new (TextEncoder as any)('utf-8');
    const decoder = new TextDecoder('utf-8');

    (global as any).Go = class {
        constructor() {
            (this as any).argv = ['js'];
            (this as any).env = {};
            (this as any).exit = (code: number) => {
                if (code !== 0) {
                    console.warn('exit code:', code);
                }
            };
            (this as any)._exitPromise = new Promise(resolve => {
                (this as any)._resolveExitPromise = resolve;
            });
            (this as any)._pendingEvent = null;
            (this as any)._scheduledTimeouts = new Map();
            (this as any)._nextCallbackTimeoutID = 1;

            const mem = () => {
                // The buffer may change when requesting more memory.
                return new DataView((this as any)._inst.exports.mem.buffer);
            };

            const setInt64 = (addr: number, v: number) => {
                mem().setUint32(addr + 0, v, true);
                mem().setUint32(addr + 4, Math.floor(v / 4294967296), true);
            };

            const getInt64 = (addr: number) => {
                const low = mem().getUint32(addr + 0, true);
                const high = mem().getInt32(addr + 4, true);
                return low + high * 4294967296;
            };

            const loadValue = (addr: number) => {
                const f = mem().getFloat64(addr, true);
                if (f === 0) {
                    return undefined;
                }
                if (!isNaN(f)) {
                    return f;
                }

                const id = mem().getUint32(addr, true);
                return (this as any)._values[id];
            };

            const storeValue = (addr: number, v: string | number | Uint8Array | boolean) => {
                const nanHead = 0x7ff80000;

                if (typeof v === 'number') {
                    if (isNaN(v)) {
                        mem().setUint32(addr + 4, nanHead, true);
                        mem().setUint32(addr, 0, true);
                        return;
                    }
                    if (v === 0) {
                        mem().setUint32(addr + 4, nanHead, true);
                        mem().setUint32(addr, 1, true);
                        return;
                    }
                    mem().setFloat64(addr, v, true);
                    return;
                }

                switch (v) {
                    case undefined:
                        mem().setFloat64(addr, 0, true);
                        return;
                    case null:
                        mem().setUint32(addr + 4, nanHead, true);
                        mem().setUint32(addr, 2, true);
                        return;
                    case true:
                        mem().setUint32(addr + 4, nanHead, true);
                        mem().setUint32(addr, 3, true);
                        return;
                    case false:
                        mem().setUint32(addr + 4, nanHead, true);
                        mem().setUint32(addr, 4, true);
                        return;
                }

                let ref = (this as any)._refs.get(v);
                if (ref === undefined) {
                    ref = (this as any)._values.length;
                    (this as any)._values.push(v);
                    (this as any)._refs.set(v, ref);
                }
                let typeFlag = 0;
                switch (typeof v) {
                    case 'string':
                        typeFlag = 1;
                        break;
                    case 'symbol':
                        typeFlag = 2;
                        break;
                    case 'function':
                        typeFlag = 3;
                        break;
                }
                mem().setUint32(addr + 4, nanHead | typeFlag, true);
                mem().setUint32(addr, ref, true);
            };

            const loadSlice = (addr: number) => {
                const array = getInt64(addr + 0);
                const len = getInt64(addr + 8);
                return new Uint8Array((this as any)._inst.exports.mem.buffer, array, len);
            };

            const loadSliceOfValues = (addr: number) => {
                const array = getInt64(addr + 0);
                const len = getInt64(addr + 8);
                const a = new Array(len);
                for (let i = 0; i < len; i++) {
                    a[i] = loadValue(array + i * 8);
                }
                return a;
            };

            const loadString = (addr: number) => {
                const saddr = getInt64(addr + 0);
                const len = getInt64(addr + 8);
                return decoder.decode(new DataView((this as any)._inst.exports.mem.buffer, saddr, len));
            };

            const timeOrigin = Date.now() - performance.now();
            (this as any).importObject = {
                go: {
                    // Go's SP does not change as long as no Go code is running. Some operations (e.g. calls, getters and setters)
                    // may synchronously trigger a Go event handler. This makes Go code get executed in the middle of the imported
                    // function. A goroutine can switch to a new stack if the current stack is too small (see morestack function).
                    // This changes the SP, thus we have to update the SP used by the imported function.

                    // func wasmExit(code int32)
                    'runtime.wasmExit': (sp: number) => {
                        const code = mem().getInt32(sp + 8, true);
                        (this as any).exited = true;
                        delete (this as any)._inst;
                        delete (this as any)._values;
                        delete (this as any)._refs;
                        (this as any).exit(code);
                    },

                    // func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)
                    'runtime.wasmWrite': (sp: number) => {
                        const fd = getInt64(sp + 8);
                        const p = getInt64(sp + 16);
                        const n = mem().getInt32(sp + 24, true);
                        (global as any).fs.writeSync(fd, new Uint8Array((this as any)._inst.exports.mem.buffer, p, n));
                    },

                    // func nanotime() int64
                    'runtime.nanotime': (sp: number) => {
                        setInt64(sp + 8, (timeOrigin + performance.now()) * 1000000);
                    },

                    // func walltime() (sec int64, nsec int32)
                    'runtime.walltime': (sp: number) => {
                        const msec = new Date().getTime();
                        setInt64(sp + 8, msec / 1000);
                        mem().setInt32(sp + 16, (msec % 1000) * 1000000, true);
                    },

                    // func scheduleTimeoutEvent(delay int64) int32
                    'runtime.scheduleTimeoutEvent': (sp: number) => {
                        const id = (this as any)._nextCallbackTimeoutID;
                        (this as any)._nextCallbackTimeoutID++;
                        (this as any)._scheduledTimeouts.set(
                            id,
                            setTimeout(
                                () => {
                                    (this as any)._resume();
                                    while ((this as any)._scheduledTimeouts.has(id)) {
                                        // for some reason Go failed to register the timeout event, log and try again
                                        // (temporary workaround for https://github.com/golang/go/issues/28975)
                                        console.warn('scheduleTimeoutEvent: missed timeout event');
                                        (this as any)._resume();
                                    }
                                },
                                getInt64(sp + 8) + 1, // setTimeout has been seen to fire up to 1 millisecond early
                            ),
                        );
                        mem().setInt32(sp + 16, id, true);
                    },

                    // func clearTimeoutEvent(id int32)
                    'runtime.clearTimeoutEvent': (sp: number) => {
                        const id = mem().getInt32(sp + 8, true);
                        clearTimeout((this as any)._scheduledTimeouts.get(id));
                        (this as any)._scheduledTimeouts.delete(id);
                    },

                    // func getRandomData(r []byte)
                    'runtime.getRandomData': (sp: number) => {
                        crypto.getRandomValues(loadSlice(sp + 8));
                    },

                    // func stringVal(value string) ref
                    'syscall/js.stringVal': (sp: number) => {
                        storeValue(sp + 24, loadString(sp + 8));
                    },

                    // func valueGet(v ref, p string) ref
                    'syscall/js.valueGet': (sp: number) => {
                        const result = Reflect.get(loadValue(sp + 8), loadString(sp + 16));
                        sp = (this as any)._inst.exports.getsp(); // see comment above
                        storeValue(sp + 32, result);
                    },

                    // func valueSet(v ref, p string, x ref)
                    'syscall/js.valueSet': (sp: number) => {
                        Reflect.set(loadValue(sp + 8), loadString(sp + 16), loadValue(sp + 32));
                    },

                    // func valueIndex(v ref, i int) ref
                    'syscall/js.valueIndex': (sp: number) => {
                        storeValue(sp + 24, Reflect.get(loadValue(sp + 8), getInt64(sp + 16)));
                    },

                    // valueSetIndex(v ref, i int, x ref)
                    'syscall/js.valueSetIndex': (sp: number) => {
                        Reflect.set(loadValue(sp + 8), getInt64(sp + 16), loadValue(sp + 24));
                    },

                    // func valueCall(v ref, m string, args []ref) (ref, bool)
                    'syscall/js.valueCall': (sp: number) => {
                        try {
                            const v = loadValue(sp + 8);
                            const m = Reflect.get(v, loadString(sp + 16));
                            const args = loadSliceOfValues(sp + 32);
                            const result = Reflect.apply(m, v, args);
                            sp = (this as any)._inst.exports.getsp(); // see comment above
                            storeValue(sp + 56, result);
                            mem().setUint8(sp + 64, 1);
                        } catch (err) {
                            storeValue(sp + 56, err);
                            mem().setUint8(sp + 64, 0);
                        }
                    },

                    // func valueInvoke(v ref, args []ref) (ref, bool)
                    'syscall/js.valueInvoke': (sp: number) => {
                        try {
                            const v = loadValue(sp + 8);
                            const args = loadSliceOfValues(sp + 16);
                            const result = Reflect.apply(v, undefined, args);
                            sp = (this as any)._inst.exports.getsp(); // see comment above
                            storeValue(sp + 40, result);
                            mem().setUint8(sp + 48, 1);
                        } catch (err) {
                            storeValue(sp + 40, err);
                            mem().setUint8(sp + 48, 0);
                        }
                    },

                    // func valueNew(v ref, args []ref) (ref, bool)
                    'syscall/js.valueNew': (sp: number) => {
                        try {
                            const v = loadValue(sp + 8);
                            const args = loadSliceOfValues(sp + 16);
                            const result = Reflect.construct(v, args);
                            sp = (this as any)._inst.exports.getsp(); // see comment above
                            storeValue(sp + 40, result);
                            mem().setUint8(sp + 48, 1);
                        } catch (err) {
                            storeValue(sp + 40, err);
                            mem().setUint8(sp + 48, 0);
                        }
                    },

                    // func valueLength(v ref) int
                    'syscall/js.valueLength': (sp: number) => {
                        setInt64(sp + 16, parseInt(loadValue(sp + 8).length));
                    },

                    // valuePrepareString(v ref) (ref, int)
                    'syscall/js.valuePrepareString': (sp: number) => {
                        const str = encoder.encode(String(loadValue(sp + 8)));
                        storeValue(sp + 16, str);
                        setInt64(sp + 24, str.length);
                    },

                    // valueLoadString(v ref, b []byte)
                    'syscall/js.valueLoadString': (sp: number) => {
                        const str = loadValue(sp + 8);
                        loadSlice(sp + 16).set(str);
                    },

                    // func valueInstanceOf(v ref, t ref) bool
                    'syscall/js.valueInstanceOf': (sp: number) => {
                        mem().setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16) ? 0 : 1);
                    },

                    // func copyBytesToGo(dst []byte, src ref) (int, bool)
                    'syscall/js.copyBytesToGo': (sp: number) => {
                        const dst = loadSlice(sp + 8);
                        const src = loadValue(sp + 32);
                        if (!(src instanceof Uint8Array)) {
                            mem().setUint8(sp + 48, 0);
                            return;
                        }
                        const toCopy = src.subarray(0, dst.length);
                        dst.set(toCopy);
                        setInt64(sp + 40, toCopy.length);
                        mem().setUint8(sp + 48, 1);
                    },

                    // func copyBytesToJS(dst ref, src []byte) (int, bool)
                    'syscall/js.copyBytesToJS': (sp: number) => {
                        const dst = loadValue(sp + 8);
                        const src = loadSlice(sp + 16);
                        if (!(dst instanceof Uint8Array)) {
                            mem().setUint8(sp + 48, 0);
                            return;
                        }
                        const toCopy = src.subarray(0, dst.length);
                        dst.set(toCopy);
                        setInt64(sp + 40, toCopy.length);
                        mem().setUint8(sp + 48, 1);
                    },

                    debug: (value: any) => {
                        console.log(value);
                    },
                },
            };
        }

        async run(instance: any) {
            (this as any)._inst = instance;
            (this as any)._values = [
                // TODO: garbage collection
                NaN,
                0,
                null,
                true,
                false,
                global,
                this,
            ];
            (this as any)._refs = new Map();
            (this as any).exited = false;

            const mem = new DataView((this as any)._inst.exports.mem.buffer);

            // Pass command line arguments and environment variables to WebAssembly by writing them to the linear memory.
            let offset = 4096;

            const strPtr = (str: string) => {
                const ptr = offset;
                const bytes = encoder.encode(str + '\0');
                new Uint8Array(mem.buffer, offset, bytes.length).set(bytes);
                offset += bytes.length;
                if (offset % 8 !== 0) {
                    offset += 8 - (offset % 8);
                }
                return ptr;
            };

            const argc = (this as any).argv.length;

            const argvPtrs = [];
            (this as any).argv.forEach((arg: any) => {
                argvPtrs.push(strPtr(arg));
            });

            const keys = Object.keys((this as any).env).sort();
            argvPtrs.push(keys.length);
            keys.forEach(key => {
                argvPtrs.push(strPtr(`${key}=${(this as any).env[key]}`));
            });

            const argv = offset;
            argvPtrs.forEach(ptr => {
                mem.setUint32(offset, ptr, true);
                mem.setUint32(offset + 4, 0, true);
                offset += 8;
            });

            (this as any)._inst.exports.run(argc, argv);
            if ((this as any).exited) {
                (this as any)._resolveExitPromise();
            }
            await (this as any)._exitPromise;
        }

        _resume() {
            if ((this as any).exited) {
                throw new Error('Go program has already exited');
            }
            (this as any)._inst.exports.resume();
            if ((this as any).exited) {
                (this as any)._resolveExitPromise();
            }
        }

        _makeFuncWrapper(id: any) {
            const go = this;
            return function(this: any) {
                const event = { id: id, this: this, args: arguments };
                (go as any)._pendingEvent = event;
                go._resume();
                return (event as any).result;
            };
        }
    };

    // if (
    //     (global as any).require &&
    //     (global as any).require.main === module &&
    //     global.process &&
    //     global.process.versions &&
    //     !(global.process.versions as any).electron
    // ) {
    //     if (process.argv.length < 3) {
    //         console.error('usage: go_js_wasm_exec [wasm binary] [arguments]');
    //         process.exit(1);
    //     }

    //     const go: any = new Go();
    //     go.argv = process.argv.slice(2);
    //     go.env = Object.assign({ TMPDIR: require('os').tmpdir() }, process.env);
    //     go.exit = process.exit;
    //     WebAssembly.instantiate((global as any).fs.readFileSync(process.argv[2]), go.importObject)
    //         .then(result => {
    //             process.on('exit', code => {
    //                 // Node.js exits if no event handler is pending
    //                 if (code === 0 && !go.exited) {
    //                     // deadlock, make Go print error and stack traces
    //                     go._pendingEvent = { id: 0 };
    //                     go._resume();
    //                 }
    //             });
    //             return go.run(result.instance);
    //         })
    //         .catch(err => {
    //             console.error(err);
    //             process.exit(1);
    //         });
    // }
})();
