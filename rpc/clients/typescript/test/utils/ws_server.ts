import {logUtils} from '@0x/utils';
import {ChildProcessWithoutNullStreams, spawn} from 'child_process';

export const SERVER_PORT = 64321;

// 1 second
const DEFAULT_TIMEOUT = 1000;

/**
 * Setup a mesh node with an associated RPC server.
 */
/* tslint:disable:custom-no-magic-numbers */
export async function setupMeshNodeAsync(): Promise<ChildProcessWithoutNullStreams> {
    // Spawn a child process running the mesh executable with an rpc server.
    const env = Object.create(process.env);
    env.ETHEREUM_RPC_URL = 'http://localhost:8545';
    env.ETHEREUM_CHAIN_ID = 1337;
    env.VERBOSITY = 5;
    env.RPC_ADDR = 'localhost:'.concat(SERVER_PORT.toString());
    return spawn('mesh', [], {env});
}
/* tslint:enable:custom-no-magic-numbers */

/**
 * Given a child process running a mesh node, wait for a log on stderr that
 * matches the given regex pattern.
 * @param mesh The mesh process to scrape for logs.
 * @param pattern The regex pattern to use when testing incoming logs.
 * @param timeout An optional timeout parameter to schedule an end to waiting on the logs.
 */
export async function waitForPatternLogAsync(
    mesh: ChildProcessWithoutNullStreams,
    pattern: RegExp,
    timeout?: number,
): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        mesh.stderr.on('data', async data => {
            if (pattern.test(data.toString())) {
                resolve(data.toString());
            }
        });
        setTimeout(reject, timeout || DEFAULT_TIMEOUT);
    });
}
