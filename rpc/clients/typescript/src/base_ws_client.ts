import * as EventEmitter from 'eventemitter3';
import * as WebSocket from 'websocket';

const READY = 'ready';
const CONNECT = 'connect';
const ERROR = 'error';
const CLOSE = 'close';

const SOCKET_MESSAGE = 'socket_message';
const SOCKET_READY = 'socket_ready';
const SOCKET_CLOSE = 'socket_close';
const SOCKET_ERROR = 'socket_error';
const SOCKET_CONNECT = 'socket_connect';

/**
 * The BaseWSClient is a generic WS client with subscriptions support as defined in
 * https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB
 * It handles re-connecting closed connections due to errors and also re-subscribes
 * all active subscriptions
 */
export class BaseWSClient extends EventEmitter {
    private _timeoutIfExists: number|undefined;
    private _subscriptions: any;
    private _connection: any;
    private _host: string;
    private _reconnectionTimeoutMs: number;
    private _messageId: number;
    private static _validateJSONRPCResponse(response: any): boolean|Error {
        if (typeof response === 'object') {
            if (response.error) {
                if (response.error instanceof Error) {
                    return new Error(`Node error: ${response.error.message}`);
                }

                return new Error(`Node error: ${JSON.stringify(response.error)}`);
            }

            if (response.result === undefined) {
                return new Error('Validation error: Undefined JSON-RPC result');
            }

            return true;
        }

        return new Error('Validation error: Response should be of type Object');
    }
    // HACK(fabio): We could have used `WebSocket.connection` as the type for param `connection` but
    // the type definitions for the `websocket` package are very out-of-date and this would cause us
    // to use `as any` in most places we access it. Simply using `any` felt cleaner until the typings
    // are updated.
    constructor(connection: any, reconnectionTimeout: number = 5000, timeout: number|undefined) {
        super();
        this._messageId = 1;
        this._connection = connection;
        this._timeoutIfExists = timeout;
        this._subscriptions = {};
        this._registerEventListeners();
        this._host = this._connection.url;
        this._reconnectionTimeoutMs = reconnectionTimeout;
    }
    /**
     * Creates the JSON-RPC payload and sends it to the node.
     * @param method JSON-RPC method to call
     * @param parameters parameters to send to method call
     * @returns response to JSON-RPC call
     */
    protected async _sendAsync(method: string, parameters: any): Promise<any> {
        const response = await this._sendPayloadAsync(this._toPayload(method, parameters));
        const validationResult = BaseWSClient._validateJSONRPCResponse(response);
        if (validationResult instanceof Error) {
            throw validationResult;
        }
        return response.result;
    }
    /**
     * Will close the socket connection with a error code and reason.
     * Please have a look at https://developer.mozilla.org/de/docs/Web/API/WebSocket/close
     * for further information.
     * @param code WebSocket error code
     * @param reason error description
     */
    protected _disconnect(code: number | null = null, reason: string | null = null): void {
        this._connection.close(code, reason);
    }
    /**
     * Returns true if the socket connection state is OPEN
     * @returns whether we are connected
     */
    protected _isConnected(): boolean {
        return this._connection.readyState === this._connection.OPEN;
    }
    /**
     * Subscribes to a given subscriptionType
     * @param subscribeMethod JSON-RPC method name to use for subscription
     * @param subscriptionMethod Subscription namespace
     * @param parameters Additional parameters to subscribe call
     * @returns The subscriptionId of an error
     */
    protected _subscribeAsync(subscribeMethod: string, subscriptionMethod: string, parameters: any[]): Promise<any> {
        parameters.unshift(subscriptionMethod);

        return this._sendAsync(subscribeMethod, parameters)
            .then(subscriptionId => {
                this._subscriptions[subscriptionId] = {
                    id: subscriptionId,
                    subscribeMethod,
                    parameters,
                };

                return subscriptionId;
            })
            .catch(error => {
                throw new Error(`Provider error: ${error}`);
            });
    }
    /**
     * Unsubscribes the subscription by his id
     * @param subscriptionId subscription identifier corresponding to the subscription to cancel
     * @param unsubscribeMethod The JSON-RPC subscription method name
     *
     * @returns either a boolean of whether the subscription was cancelled, or an error
     */
    protected _unsubscribeAsync(subscriptionId: string, unsubscribeMethod: string): Promise<boolean|Error> {
        if (this._hasSubscription(subscriptionId)) {
            return this._sendAsync(unsubscribeMethod, [subscriptionId]).then(response => {
                if (response) {
                    this._removeAllListeners(this._getSubscriptionEvent(subscriptionId));

                    delete this._subscriptions[subscriptionId];
                }

                return response;
            });
        }

        return Promise.reject(new Error(`Provider error: Subscription with ID ${subscriptionId} does not exist.`));
    }
    /**
     * Clears all _subscriptions and listeners
     * @param unsubscribeMethod JSON-RPC unsubscribe method
     * @returns true if clearing subscription succeeds, otherwise an error
     */
    protected async _clearSubscriptionsAsync(unsubscribeMethod: string): Promise<boolean|Error> {
        const unsubscribePromises: Array<Promise<any>> = [];

        Object.keys(this._subscriptions).forEach(key => {
            this._removeAllListeners(key);
            unsubscribePromises.push(this._unsubscribeAsync(this._subscriptions[key].id, unsubscribeMethod));
        });

        return Promise.all(unsubscribePromises).then(results => {
            if (results.includes(false)) {
                throw new Error(`Could not unsubscribe all _subscriptions: ${JSON.stringify(results)}`);
            }

            return true;
        });
    }
    /**
     * Sends the JSON-RPC payload to the node.
     * @param payload JSON-RPC payload to send
     *
     * @returns the response received with the matching id specified in the payload
     */
    private async _sendPayloadAsync(payload: any): Promise<any> {
        return new Promise((resolve, reject) => {
            this.once('error', reject);

            if (!this._isConnecting()) {
                let timeout: any;

                if (this._connection.readyState !== this._connection.OPEN) {
                    this.removeListener('error', reject);

                    return reject(new Error('Connection error: Connection is not open on send()'));
                }

                try {
                    this._connection.send(JSON.stringify(payload));
                } catch (error) {
                    this.removeListener('error', reject);

                    return reject(error);
                }

                if (this._timeoutIfExists) {
                    timeout = setTimeout(() => {
                        reject(new Error('Connection error: Timeout exceeded'));
                    }, this._timeoutIfExists);
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
                this._sendPayloadAsync(payload)
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
    /**
     * Emits the ready event when the connection is established
     * @param event Event to emit on ready
     */
    private _onReady(event: any): void {
        (this as EventEmitter).emit(READY, event);
        (this as EventEmitter).emit(SOCKET_READY, event);
    }
    /**
     * This is the listener for the 'close' event of the current socket connection.
     * @param closeEvent close event
     */
    private _onClose(closeEvent: CloseEvent): void {
        if (closeEvent.code !== WebSocket.connection.CLOSE_REASON_NORMAL || closeEvent.wasClean === false) {
            this._reconnect();

            return;
        }
        (this as EventEmitter).emit(CLOSE, closeEvent);
        (this as EventEmitter).emit(SOCKET_CLOSE, closeEvent);
        this._removeAllSocketListeners();
        this._removeAllListeners();
    }
    /**
     * This is the listener for the 'error' event of the current socket connection.
     * @param event error event
     */
    private _onError(event: any): void {
        if (event.code === 'ECONNREFUSED') {
            this._reconnect();

            return;
        }
        (this as EventEmitter).emit(ERROR, event);
        (this as EventEmitter).emit(SOCKET_ERROR, event);
        this._removeAllSocketListeners();
    }
    /**
     * This is the listener for the 'message' events of the current socket connection.
     * @param messageEvent message event
     */
    private _onMessage(messageEvent: MessageEvent): void {
        const response = messageEvent.data;
        let event;

        let responseObject = response as any;
        if (typeof response !== 'object') {
            responseObject = JSON.parse(response);
        }

        if (Array.isArray(responseObject)) {
            event = responseObject[0].id;
        } else if (typeof responseObject.id === 'undefined') {
            event = this._getSubscriptionEvent(responseObject.params.subscription);
            responseObject = responseObject.params;
        } else {
            event = responseObject.id;
        }

        (this as EventEmitter).emit(SOCKET_MESSAGE, responseObject);
        (this as EventEmitter).emit(event, responseObject);
    }
    /**
     * Emits the connect event and checks if there are _subscriptions defined that should be resubscribed.
     */
    private async _onConnectAsync(): Promise<void> {
        const subscriptionKeys = Object.keys(this._subscriptions);

        if (subscriptionKeys.length > 0) {
            let subscriptionId;

            for (const key of subscriptionKeys) {
                subscriptionId = await this._subscribeAsync(
                    this._subscriptions[key].subscribeMethod,
                    this._subscriptions[key].parameters[0],
                    this._subscriptions[key].parameters.slice(1),
                );

                if (key !== subscriptionId) {
                    delete this._subscriptions[subscriptionId];
                }

                this._subscriptions[key].id = subscriptionId;
            }
        }

        (this as EventEmitter).emit(SOCKET_CONNECT);
        (this as EventEmitter).emit(CONNECT);
    }
    /**
     * Returns if the socket connection is in the connecting state.
     * @returns whether we are connecting
     */
    private _isConnecting(): boolean {
        return this._connection.readyState === this._connection.CONNECTING;
    }
    /**
     * Removes all listeners on the EventEmitter and the socket object.
     * @param event socket event name
     */
    private _removeAllListeners(event?: string | symbol | undefined): any {
        switch (event) {
            case SOCKET_MESSAGE:
                this._connection.removeEventListener('message', this._onMessage);
                break;
            case SOCKET_READY:
                this._connection.removeEventListener('open', this._onReady);
                break;
            case SOCKET_CLOSE:
                this._connection.removeEventListener('close', this._onClose);
                break;
            case SOCKET_ERROR:
                this._connection.removeEventListener('error', this._onError);
                break;
            case SOCKET_CONNECT:
                this._connection.removeEventListener('connect', this._onConnectAsync);
                break;
            default:
                // Noop
        }
        super.removeAllListeners(event);
    }
    /**
     * Registers all the required listeners.
     */
    private _registerEventListeners(): void {
        this._connection.addEventListener('message', this._onMessage.bind(this));
        this._connection.addEventListener('open', this._onReady.bind(this));
        this._connection.addEventListener('open', this._onConnectAsync.bind(this));
        this._connection.addEventListener('close', this._onClose.bind(this));
        this._connection.addEventListener('error', this._onError.bind(this));
    }
    /**
     * Removes all socket listeners
     */
    private _removeAllSocketListeners(): void {
        this._removeAllListeners(SOCKET_MESSAGE);
        this._removeAllListeners(SOCKET_READY);
        this._removeAllListeners(SOCKET_CLOSE);
        this._removeAllListeners(SOCKET_ERROR);
        this._removeAllListeners(SOCKET_CONNECT);
    }
    /**
     * Removes the listeners and reconnects to the socket.
     */
    private _reconnect(): void {
        setTimeout(() => {
            this._removeAllSocketListeners();

            let connection = [];

            if (this._connection.constructor.name === 'W3CWebSocket') {
                connection = new this._connection.constructor(
                    this._host,
                    this._connection._client.protocol,
                    null,
                    this._connection._client.headers,
                    this._connection._client.requestOptions,
                    this._connection._client.config,
                );
            } else {
                connection = new this._connection.constructor(this._host, this._connection.protocol || undefined);
            }

            this._connection = connection;
            this._registerEventListeners();
            // Emit a "reconnected" event only once the new connection is established
            this.once('connect', () => {
                this.emit('reconnected');
            });
        }, this._reconnectionTimeoutMs);
    }
    /**
     * Checks if the given subscription id exists
     * @param subscriptionId subscription ID to check existence for
     * @returns whether or not the subscription exists
     */
    private _hasSubscription(subscriptionId: string): boolean {
        return typeof this._getSubscriptionEvent(subscriptionId) !== 'undefined';
    }
    /**
     * Returns the event the subscription is listening for.
     * @param subscriptionId subscription ID
     * @returns subscription event name (e.g. "heartbeat")
     */
    private _getSubscriptionEvent(subscriptionId: string): string|undefined {
        if (this._subscriptions[subscriptionId]) {
            return subscriptionId;
        }

        let event: string|undefined;
        for (const key in this._subscriptions) {
            if (this._subscriptions[key].id === subscriptionId) {
                event = key;
            }
        }

        return event;
    }
    /**
     * Creates a valid json payload object
     * @param method JSON-RPC method to call
     * @param params parameters to supply in method call
     * @returns JSON-RPC payload
     */
    private _toPayload(method: string, params: any[]): any {
        if (!method) {
            throw new Error(`JSON-RPC method should be specified in payload: "${JSON.stringify(params)}"`);
        }

        const id = this._messageId;
        this._messageId++;

        return {
            jsonrpc: '2.0',
            id,
            method,
            params: params || [],
        };
    }
}
