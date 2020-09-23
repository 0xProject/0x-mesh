import { ChildProcessWithoutNullStreams, spawn } from 'child_process';
import * as jsonstream from 'jsonstream';
import { join } from 'path';
import * as rimraf from 'rimraf';

import { MeshGraphQLClient } from '../../src';

const dataDir = '/tmp/mesh-graphql-integration-testing/data/';

async function buildBinaryAsync(): Promise<void> {
    const cwd = join(__dirname, '../../../../../').normalize();
    const build = spawn('make', ['mesh'], { cwd });
    await new Promise<void>((resolve, reject) => {
        build.on('close', (code) => {
            code === 0 ? resolve() : reject(new Error('Failed to build 0x-mesh'));
        });
        build.on('error', (error) => {
            reject(error);
        });
    });
}

async function cleanupAsync(): Promise<void> {
    await new Promise<void>((resolve, reject) => {
        rimraf(dataDir, (err) => {
            if (err != null) {
                reject(err);
            }
            resolve();
        });
    });
}

export interface MeshDeployment {
    client: MeshGraphQLClient;
    mesh: MeshHarness;
    peerID: string;
}

// The amount of time to wait after seeing the "starting GraphQL server" log message
// before attempting to connect to the GraphQL server.
const serverStartWaitTimeMs = 100;

/**
 * Start a GraphQL client connected to a GraphQL server that is ready for use.
 * @return A mesh deployment including a GraphQL client, mesh manager, and the
 *         peer ID of the mesh process that is running in the mesh manager.
 */
export async function startServerAndClientAsync(): Promise<MeshDeployment> {
    await cleanupAsync();
    await buildBinaryAsync();

    const mesh = new MeshHarness();
    const log = await mesh.waitForPatternAsync(/starting GraphQL server/);
    const peerID = JSON.parse(log.toString()).myPeerID;
    await sleepAsync(serverStartWaitTimeMs);
    const client = new MeshGraphQLClient({
        httpUrl: `http://localhost:${mesh._graphQLServerPort}/graphql`,
        webSocketUrl: `ws://localhost:${mesh._graphQLServerPort}/graphql`,
    });
    return {
        client,
        mesh,
        peerID,
    };
}

export class MeshHarness {
    public static readonly DEFAULT_TIMEOUT = 1000;
    protected static _serverPort = 64321;

    public readonly _graphQLServerPort: number;
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
            const stream = jsonstream.parse(true);
            stream.on('data', (data) => {
                // Note(albrow): Uncomment this if you need to see the output from the server.
                // console.log(data);
                const dataString = JSON.stringify(data);
                if (pattern.test(dataString)) {
                    resolve(dataString);
                }
            });
            this._mesh.stderr.pipe(stream);
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
        this._graphQLServerPort = MeshHarness._serverPort++;
        env.DATA_DIR = dataDir;
        env.ETHEREUM_RPC_URL = 'http://localhost:8545';
        env.ETHEREUM_CHAIN_ID = '1337';
        env.VERBOSITY = '5';
        env.USE_BOOTSTRAP_LIST = false;
        env.ENABLE_GRAPHQL_SERVER = true;
        env.GRAPHQL_SERVER_ADDR = `localhost:${this._graphQLServerPort}`;
        this._mesh = spawn('mesh', [], { env });
        this._mesh.stderr.on('error', (error) => {
            throw new Error(`${error.name} - ${error.message}`);
        });
    }
}

async function sleepAsync(ms: number): Promise<NodeJS.Timer> {
    return new Promise<NodeJS.Timer>((resolve) => setTimeout(resolve, ms));
}
