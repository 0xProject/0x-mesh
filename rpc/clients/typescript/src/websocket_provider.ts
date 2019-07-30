import * as WebSocket from 'websocket';

import { SocketProvider } from './socket_provider';

const SOCKET_MESSAGE = 'socket_message';
const SOCKET_READY = 'socket_ready';
const SOCKET_CLOSE = 'socket_close';
const SOCKET_ERROR = 'socket_error';
const SOCKET_CONNECT = 'socket_connect';

export class WebSocketProvider extends SocketProvider {
    private _host: string;
    private _reconnectionTimeoutMs: number;
    constructor(connection: any, reconnectionTimeout: number = 5000, timeout: number|undefined) {
        super(connection, timeout);
        // HACK(fabio): @types/websocket out-of-date
        this._host = this.connection.url;
        this._reconnectionTimeoutMs = reconnectionTimeout;
    }
    /**
     * This is the listener for the 'message' events of the current socket connection.
     * @param messageEvent message event
     */
    public onMessage(messageEvent: MessageEvent): void {
        super.onMessage(messageEvent.data);
    }
    /**
     * This is the listener for the 'error' event of the current socket connection.
     * @param event error event
     */
    public onError(event: any): void {
        if (event.code === 'ECONNREFUSED') {
            this.reconnect();

            return;
        }
        super.onError(event);
    }
    /**
     * This is the listener for the 'close' event of the current socket connection.
     * @method onClose
     * @param closeEvent close event
     */
    public onClose(closeEvent: CloseEvent): void {
        if (closeEvent.code !== WebSocket.connection.CLOSE_REASON_NORMAL || closeEvent.wasClean === false) {
            this.reconnect();

            return;
        }
        super.onClose();
    }
    /**
     * Removes the listeners and reconnects to the socket.
     */
    public reconnect(): void {
        setTimeout(() => {
            this.removeAllSocketListeners();

            let connection = [];

            if (this.connection.constructor.name === 'W3CWebSocket') {
                connection = new this.connection.constructor(
                    this._host,
                    this.connection._client.protocol,
                    null,
                    this.connection._client.headers,
                    this.connection._client.requestOptions,
                    this.connection._client.config,
                );
            } else {
                connection = new this.connection.constructor(this._host, this.connection.protocol || undefined);
            }

            this.connection = connection;
            this.registerEventListeners();
            // Emit a "reconnected" event once new connection established
            this.once('connect', () => {
                this.emit('reconnected');
            });
        }, this._reconnectionTimeoutMs);
    }
    /**
     * Will close the socket connection with a error code and reason.
     * Please have a look at https://developer.mozilla.org/de/docs/Web/API/WebSocket/close
     * for further information.
     * @param code WebSocket error code
     * @param reason error description
     */
    public disconnect(code: number | null = null, reason: string | null = null): void {
        this.connection.close(code, reason);
    }
    /**
     * Registers all the required listeners.
     */
    public registerEventListeners(): void {
        this.connection.addEventListener('message', this.onMessage.bind(this));
        this.connection.addEventListener('open', this.onReady.bind(this));
        this.connection.addEventListener('open', this.onConnectAsync.bind(this));
        this.connection.addEventListener('close', this.onClose.bind(this));
        this.connection.addEventListener('error', this.onError.bind(this));
    }
    /**
     * Removes all listeners on the EventEmitter and the socket object.
     * @param event socket event name
     */
    public removeAllListeners(event?: string | symbol | undefined): any {
        switch (event) {
            case SOCKET_MESSAGE:
                this.connection.removeEventListener('message', this.onMessage);
                break;
            case SOCKET_READY:
                this.connection.removeEventListener('open', this.onReady);
                break;
            case SOCKET_CLOSE:
                this.connection.removeEventListener('close', this.onClose);
                break;
            case SOCKET_ERROR:
                this.connection.removeEventListener('error', this.onError);
                break;
            case SOCKET_CONNECT:
                this.connection.removeEventListener('connect', this.onConnectAsync);
                break;
            default:
                // Noop
        }
        super.removeAllListeners(event);
    }
    /**
     * Returns true if the socket connection state is OPEN
     * @returns whether we are connected
     */
    public isConnected(): boolean {
        return this.connection.readyState === this.connection.OPEN;
    }
    /**
     * Returns if the socket connection is in the connecting state.
     * @returns whether we are connecting
     */
    public isConnecting(): boolean {
        return this.connection.readyState === this.connection.CONNECTING;
    }
    /**
     * Sends the JSON-RPC payload to the node.
     * @param payload JSON-RPC payload to send
     *
     * @returns the response received with the matching id specified in the payload
     */
    // tslint:disable-next-line:async-suffix
    public async sendPayloadAsync(payload: any): Promise<any> {
        return new Promise((resolve, reject) => {
            this.once('error', reject);

            if (!this.isConnecting()) {
                let timeout: any;

                if (this.connection.readyState !== this.connection.OPEN) {
                    this.removeListener('error', reject);

                    return reject(new Error('Connection error: Connection is not open on send()'));
                }

                try {
                    this.connection.send(JSON.stringify(payload));
                } catch (error) {
                    this.removeListener('error', reject);

                    return reject(error);
                }

                if (this.timeoutIfExists) {
                    timeout = setTimeout(() => {
                        reject(new Error('Connection error: Timeout exceeded'));
                    }, this.timeoutIfExists);
                }

                const id = Array.isArray(payload) ? payload[0].id : payload.id;
                this.once(id, response => {
                    if (timeout) {
                        clearTimeout(timeout);
                    }

                    this.removeListener('error', reject);

                    return resolve(response);
                });

                return;
            }

            this.once('connect', () => {
                this.sendPayloadAsync(payload)
                    .then(response => {
                        this.removeListener('error', reject);

                        return resolve(response);
                    })
                    .catch(error => {
                        this.removeListener('error', reject);

                        return reject(error);
                    });
            });
        });
    }
}
