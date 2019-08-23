// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO_LICENSE file.

/* tslint:disable */
(() => {
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
    // environments, we can skip some Node.js-related logic.

    // // Map web browser API and Node.js API to a single common API (preferring web standards over Node.js API).
    // const isNodeJS = global.process && global.process.title === "node";
    // if (isNodeJS) {
    // 	global.require = require;
    // 	global.fs = require("fs");

    // 	const nodeCrypto = require("crypto");
    // 	global.crypto = {
    // 		getRandomValues(b) {
    // 			nodeCrypto.randomFillSync(b);
    // 		},
    // 	};

    // 	global.performance = {
    // 		now() {
    // 			const [sec, nsec] = process.hrtime();
    // 			return sec * 1000 + nsec / 1000000;
    // 		},
    // 	};

    // 	const util = require("util");
    // 	global.TextEncoder = util.TextEncoder;
    // 	global.TextDecoder = util.TextDecoder;
    // } else {
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
        write(fd: any, buf: any, offset: any, length: any, position: any, callback: any) {
            if (offset !== 0 || length !== buf.length || position !== null) {
                throw new Error('not implemented');
            }
            const n = this.writeSync(fd, buf);
            callback(null, n);
        },
        open(path: any, flags: any, mode: any, callback: any) {
            const err = new Error('not implemented');
            (err as any).code = 'ENOSYS';
            callback(err);
        },
        read(fd: any, buffer: any, offset: any, length: any, position: any, callback: any) {
            const err = new Error('not implemented');
            (err as any).code = 'ENOSYS';
            callback(err);
        },
        fsync(fd: any, callback: any) {
            callback(null);
        },
    };
    // }

    const encoder = new (TextEncoder as any)('utf-8');
    const decoder = new TextDecoder('utf-8');

    (global as any).Go = class {
        argv: any;
        env: any;
        exit: any;
        _callbackTimeouts: any;
        _nextCallbackTimeoutID: any;
        _inst: any;
        _values: any;
        _refs: any;
        importObject: any;
        exited: any;
        _callbackShutdown: any;
        _exitPromise: any;
        _resolveExitPromise: any;
        _pendingEvent: any;
        _scheduledTimeouts: any;
        constructor() {
            this.argv = ['js'];
            this.env = {};
            this.exit = (code: any) => {
                if (code !== 0) {
                    console.warn('exit code:', code);
                }
            };
            this._exitPromise = new Promise(resolve => {
                this._resolveExitPromise = resolve;
            });
            this._pendingEvent = null;
            this._scheduledTimeouts = new Map();
            this._nextCallbackTimeoutID = 1;

            const mem = () => {
                // The buffer may change when requesting more memory.
                return new DataView(this._inst.exports.mem.buffer);
            };

            const setInt64 = (addr: any, v: any) => {
                mem().setUint32(addr + 0, v, true);
                mem().setUint32(addr + 4, Math.floor(v / 4294967296), true);
            };

            const getInt64 = (addr: any) => {
                const low = mem().getUint32(addr + 0, true);
                const high = mem().getInt32(addr + 4, true);
                return low + high * 4294967296;
            };

            const loadValue = (addr: any) => {
                const f = mem().getFloat64(addr, true);
                if (f === 0) {
                    return undefined;
                }
                if (!isNaN(f)) {
                    return f;
                }

                const id = mem().getUint32(addr, true);
                return this._values[id];
            };

            const storeValue = (addr: number, v: any) => {
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

                let ref = this._refs.get(v);
                if (ref === undefined) {
                    ref = this._values.length;
                    this._values.push(v);
                    this._refs.set(v, ref);
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
                return new Uint8Array(this._inst.exports.mem.buffer, array, len);
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
                return decoder.decode(new DataView(this._inst.exports.mem.buffer, saddr, len));
            };

            const timeOrigin = Date.now() - performance.now();
            this.importObject = {
                go: {
                    // Go's SP does not change as long as no Go code is running. Some operations (e.g. calls, getters and setters)
                    // may synchronously trigger a Go event handler. This makes Go code get executed in the middle of the imported
                    // function. A goroutine can switch to a new stack if the current stack is too small (see morestack function).
                    // This changes the SP, thus we have to update the SP used by the imported function.

                    // func wasmExit(code int32)
                    'runtime.wasmExit': (sp: number) => {
                        const code = mem().getInt32(sp + 8, true);
                        this.exited = true;
                        delete this._inst;
                        delete this._values;
                        delete this._refs;
                        this.exit(code);
                    },

                    // func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)
                    'runtime.wasmWrite': (sp: number) => {
                        const fd = getInt64(sp + 8);
                        const p = getInt64(sp + 16);
                        const n = mem().getInt32(sp + 24, true);
                        (global as any).fs.writeSync(fd, new Uint8Array(this._inst.exports.mem.buffer, p, n));
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
                        const id = this._nextCallbackTimeoutID;
                        this._nextCallbackTimeoutID++;
                        this._scheduledTimeouts.set(
                            id,
                            setTimeout(
                                () => {
                                    this._resume();
                                },
                                getInt64(sp + 8) + 1, // setTimeout has been seen to fire up to 1 millisecond early
                            ),
                        );
                        mem().setInt32(sp + 16, id, true);
                    },

                    // func clearTimeoutEvent(id int32)
                    'runtime.clearTimeoutEvent': (sp: number) => {
                        const id = mem().getInt32(sp + 8, true);
                        clearTimeout(this._scheduledTimeouts.get(id));
                        this._scheduledTimeouts.delete(id);
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
                        sp = this._inst.exports.getsp(); // see comment above
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
                            sp = this._inst.exports.getsp(); // see comment above
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
                            sp = this._inst.exports.getsp(); // see comment above
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
                            sp = this._inst.exports.getsp(); // see comment above
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
                        (mem() as any).setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16));
                    },

                    debug: (value: any) => {
                        console.log(value);
                    },
                },
            };
        }

        async run(instance: any) {
            this._inst = instance;
            this._values = [
                // TODO: garbage collection
                NaN,
                0,
                null,
                true,
                false,
                global,
                this._inst.exports.mem,
                this,
            ];
            this._refs = new Map();
            this.exited = false;

            const mem = new DataView(this._inst.exports.mem.buffer);

            // Pass command line arguments and environment variables to WebAssembly by writing them to the linear memory.
            let offset = 4096;

            const strPtr = (str: string) => {
                let ptr = offset;
                new Uint8Array(mem.buffer, offset, str.length + 1).set(encoder.encode(str + '\0'));
                offset += str.length + (8 - (str.length % 8));
                return ptr;
            };

            const argc = this.argv.length;

            const argvPtrs = [];
            this.argv.forEach((arg: any) => {
                argvPtrs.push(strPtr(arg));
            });

            const keys = Object.keys(this.env).sort();
            argvPtrs.push(keys.length);
            keys.forEach(key => {
                argvPtrs.push(strPtr(`${key}=${this.env[key]}`));
            });

            const argv = offset;
            argvPtrs.forEach(ptr => {
                mem.setUint32(offset, ptr, true);
                mem.setUint32(offset + 4, 0, true);
                offset += 8;
            });

            this._inst.exports.run(argc, argv);
            if (this.exited) {
                this._resolveExitPromise();
            }
            await this._exitPromise;
        }

        _resume() {
            if (this.exited) {
                throw new Error('Go program has already exited');
            }
            this._inst.exports.resume();
            if (this.exited) {
                this._resolveExitPromise();
            }
        }

        _makeFuncWrapper(id: any) {
            const go = this;
            return function() {
                const event = { id: id, this: go, args: arguments };
                go._pendingEvent = event;
                go._resume();
                return (event as any).result;
            };
        }
    };

    // if (isNodeJS) {
    // 	if (process.argv.length < 3) {
    // 		process.stderr.write("usage: go_js_wasm_exec [wasm binary] [arguments]\n");
    // 		process.exit(1);
    // 	}

    // 	const go = new Go();
    // 	go.argv = process.argv.slice(2);
    // 	go.env = Object.assign({ TMPDIR: require("os").tmpdir() }, process.env);
    // 	go.exit = process.exit;
    // 	WebAssembly.instantiate(fs.readFileSync(process.argv[2]), go.importObject).then((result) => {
    // 		process.on("exit", (code) => { // Node.js exits if no event handler is pending
    // 			if (code === 0 && !go.exited) {
    // 				// deadlock, make Go print error and stack traces
    // 				go._pendingEvent = { id: 0 };
    // 				go._resume();
    // 			}
    // 		});
    // 		return go.run(result.instance);
    // 	}).catch((err) => {
    // 		throw err;
    // 	});
    // }
})();
