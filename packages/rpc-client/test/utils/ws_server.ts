import { ChildProcessWithoutNullStreams, spawn } from 'child_process';
import * as jsonstream from 'jsonstream';
import { join } from 'path';
import * as rimraf from 'rimraf';

import { WSClient } from '../../src';

const DEFAULT_TIMEOUT = 3000;

async function buildBinaryAsync(): Promise<void> {
    const cwd = join(__dirname, '../../../../../').normalize();
    const build = spawn('make', ['mesh'], { cwd });
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
    const log = await waitForMeshLogsAsync(mesh);
    const peerID = JSON.parse(log.toString()).myPeerID;
    const client = new WSClient(`ws://localhost:${mesh.wsPort}`);
    return {
        client,
        mesh,
        peerID,
    };
}

// Wait for the `core.App was started` and the `started WS RPC server` logs.
async function waitForMeshLogsAsync(harness: MeshHarness): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        let peerIdLog = '';
        const stream = jsonstream.parse(true);
        let didSeeWSLog = false;
        let didAppStartedLog = false;
        stream.on('data', data => {
            const dataString = JSON.stringify(data);
            if (/started WS RPC server/.test(dataString)) {
                didSeeWSLog = true;
                peerIdLog = dataString;
            } else if (/core.App was started/.test(dataString)) {
                didAppStartedLog = true;
            }

            if (didSeeWSLog && didAppStartedLog) {
                resolve(peerIdLog);
            }
        });
        harness.mesh.stderr.pipe(stream);
        setTimeout(reject, DEFAULT_TIMEOUT);
    });
}

export class MeshHarness {
    protected static _serverPort = 64321;

    public readonly wsPort: number;
    public readonly mesh: ChildProcessWithoutNullStreams;
    private _killed = false;

    /**
     * Returns a value indicating whether the harness's mesh node has been stopped.
     */
    public wasStopped(): boolean {
        return this._killed;
    }

    /**
     * Kill the mesh process of this mesh instance.
     */
    public stopMesh(): void {
        this._killed = true;
        this.mesh.kill('SIGKILL');
    }

    public constructor() {
        const env = Object.create(process.env);
        this.wsPort = MeshHarness._serverPort++;
        env.ETHEREUM_RPC_URL = 'http://localhost:8545';
        env.ETHEREUM_CHAIN_ID = '1337';
        env.VERBOSITY = '5';
        env.WS_RPC_ADDR = `localhost:${this.wsPort}`;
        this.mesh = spawn('mesh', [], { env });
        this.mesh.stderr.on('error', error => {
            throw new Error(`${error.name} - ${error.message}`);
        });
    }
}
