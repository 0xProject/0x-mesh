import {ChildProcessWithoutNullStreams, spawn} from 'child_process';
import {join} from 'path';
import * as rimraf from 'rimraf';

import {WSClient} from '../../src';

async function buildBinaryAsync(): Promise<void> {
    const cwd = join(__dirname, '../../../../../../').normalize();
    const build = spawn('make', ['mesh'], {cwd});
    await new Promise<void>((resolve, reject) => {
        build.on('close', code => {
            code === 0 ? resolve() : reject(new Error('Failed to build 0x-mesh'));
        });
        build.on('error', error => {
            reject(error);
        });
    });
}

async function cleanupAsync(): Promise<void> {
    await new Promise<void>((resolve, reject) => {
        rimraf('./0x_mesh', err => {
            if (err != null) {
                reject(err);
            }
            resolve();
        });
    });
}

export interface MeshDeployment {
    client: WSClient;
    mesh: MeshHarness;
    peerID: string;
}

/**
 * Start a RPC client connected to a RPC server that is ready for use.
 * @return A mesh deployment including a RPC client, mesh manager, and the
 *         peer ID of the mesh process that is running in the mesh manager.
 */
export async function startServerAndClientAsync(): Promise<MeshDeployment> {
    await cleanupAsync();
    await buildBinaryAsync();

    const mesh = new MeshHarness();
    const startingPattern = /started RPC server/;
    const log = await mesh.waitForPatternAsync(startingPattern);
    const peerID = JSON.parse(log.toString()).myPeerID;
    const client = new WSClient(`ws://localhost:${mesh.port}`);
    return {
        client,
        mesh,
        peerID,
    };
}

export class MeshHarness {
    public static readonly DEFAULT_TIMEOUT = 1000;
    protected static _serverPort = 64321;

    public readonly port: number;
    private readonly _mesh: ChildProcessWithoutNullStreams;
    private _killed = false;

    /**
     * Wait for a log on the mesh process's stderr that matches the given regex pattern.
     * @param pattern The regex pattern to use when testing incoming logs.
     * @param timeout An optional timeout parameter to schedule an end to waiting on the logs.
     */
    public async waitForPatternAsync(pattern: RegExp, timeout?: number): Promise<string> {
        if (this._killed) {
            throw new Error('mesh instance has already been killed');
        }
        return new Promise<string>((resolve, reject) => {
            this._mesh.stderr.on('data', async data => {
                if (pattern.test(data.toString())) {
                    // Since chunks can contain more than one log, process the chunks until a
                    // chunk is found that contains the pattern.
                    let log;
                    const chunks = data.toString().split('\n');
                    for (const chunk of chunks) {
                        if (pattern.test(chunk)) {
                            log = chunk;
                            break;
                        }
                    }
                    if (!log) {
                        throw new Error('Incorrect log found');
                    }
                    resolve(log);
                }
            });
            setTimeout(reject, timeout || MeshHarness.DEFAULT_TIMEOUT);
        });
    }

    /**
     * Kill the mesh process of this mesh instance.
     */
    public stopMesh(): void {
        this._killed = true;
        this._mesh.kill('SIGKILL');
    }

    public constructor() {
        const env = Object.create(process.env);
        this.port = MeshHarness._serverPort++;
        env.ETHEREUM_RPC_URL = 'http://localhost:8545';
        env.ETHEREUM_CHAIN_ID = '1337';
        env.VERBOSITY = '5';
        env.RPC_ADDR = `localhost:${this.port}`;
        this._mesh = spawn('mesh', [], {env});
        this._mesh.stderr.on('error', error => {
            throw new Error(`${error.name} - ${error.message}`);
        });
    }
}
