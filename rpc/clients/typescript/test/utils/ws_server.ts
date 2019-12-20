import {ChildProcessWithoutNullStreams, spawn} from 'child_process';
import {join} from 'path';

import {WSClient} from '../../src';

async function buildBinaryAsync(): Promise<void> {
    const cwd = join(__dirname, '../../../../../../').normalize();
    const build = spawn('make', ['mesh'], {cwd});
    await new Promise<void>((resolve, reject) => {
        build.on('close', code => {
            code === 0 ? resolve() : reject();
        });
    });
}

async function cleanupAsync(): Promise<void> {
    const cleanup = spawn('rm', ['-r', './0x_mesh']);
    await new Promise<void>((resolve, reject) => {
        cleanup.on('close', code => {
            // NOTE(jalextowle): In the event that the 0x_mesh files
            // are not in this directory, the "cleanup" command will
            // fail. We want to allow this so that the rest of the
            // program can execute.
            code === 0 || code === 1 ? resolve() : reject();
        });
    });
}

export interface MeshDeployment {
    client: WSClient;
    mesh: MeshManager;
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

    const mesh = new MeshManager();
    const log = await mesh.waitForPatternAsync(/started RPC server/);
    const peerID = JSON.parse(log.toString()).myPeerID;
    const client = new WSClient(`http://localhost:${mesh.port}`);
    return {
        client,
        mesh,
        peerID,
    };
}

export class MeshManager {
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
                    resolve(data.toString());
                }
            });
            setTimeout(reject, timeout || MeshManager.DEFAULT_TIMEOUT);
        });
    }

    /**
     * Kill the mesh process of this mesh instance.
     */
    public stopMesh(): void {
        this._killed = true;
        this._mesh.kill();
    }

    public constructor() {
        const env = Object.create(process.env);
        this.port = MeshManager._serverPort++;
        env.ETHEREUM_RPC_URL = 'http://localhost:8545';
        env.ETHEREUM_CHAIN_ID = '1337';
        env.VERBOSITY = '5';
        env.RPC_ADDR = `localhost:${this.port}`;
        this._mesh = spawn('mesh', [], {env});
    }
}
