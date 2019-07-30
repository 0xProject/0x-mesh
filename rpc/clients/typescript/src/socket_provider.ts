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

// tslint:disable:prefer-function-over-method async-suffix
export class SocketProvider extends EventEmitter {
    public connection: any;
    public timeoutIfExists: number|undefined;
    public subscriptions: any;
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
    /**
     * isConnected checks if the ws connection is connected
     * @returns Returns true if the socket is connected
     */
    public isConnected(): boolean {
        throw new Error('NOT IMPLEMENTED');
    }
    constructor(connection: WebSocket.connection, timeout: number|undefined) {
        super();
        this._messageId = 1;
        this.connection = connection;
        this.timeoutIfExists = timeout;
        this.subscriptions = {};
        this.registerEventListeners();
    }
    /**
     * Method for checking subscriptions support of a internal provider
     * @returns whether we support subscriptionsx
     */
    public supportsSubscriptions(): boolean {
        return true;
    }
    /**
     * Registers all the required listeners.
     */
    public registerEventListeners(): void {
        throw new Error('NOT_IMPLEMENTED');
    }
    /**
     * Registers all the required listeners.
     */
    public async sendPayloadAsync(_payload: any): Promise<any> {
        throw new Error('NOT_IMPLEMENTED');
    }
    /**
     * Removes all socket listeners
     */
    public removeAllSocketListeners(): void {
        (this as EventEmitter).removeAllListeners(SOCKET_MESSAGE);
        (this as EventEmitter).removeAllListeners(SOCKET_READY);
        (this as EventEmitter).removeAllListeners(SOCKET_CLOSE);
        (this as EventEmitter).removeAllListeners(SOCKET_ERROR);
        (this as EventEmitter).removeAllListeners(SOCKET_CONNECT);
    }
    /**
     * Closes the socket connection.
     * @param code WebSocket error code
     * @param reason error description
     */
    public disconnect(_code: number, _reason: string): void {
        throw new Error('NOT_IMPLEMENTED');
    }
    /**
     * Creates the JSON-RPC payload and sends it to the node.
     * @param method JSON-RPC method to call
     * @param parameters parameters to send to method call
     * @returns response to JSON-RPC call
     */
    public async sendAsync(method: string, parameters: any): Promise<any> {
        const response = await this.sendPayloadAsync(this.toPayload(method, parameters));
        const validationResult = SocketProvider._validateJSONRPCResponse(response);

        if (validationResult instanceof Error) {
            throw validationResult;
        }

        return response.result;
    }
    /**
     * Emits the ready event when the connection is established
     * @param event Event to emit on ready
     */
    public onReady(event: any): void {
        (this as EventEmitter).emit(READY, event);
        (this as EventEmitter).emit(SOCKET_READY, event);
    }

    /**
     * Emits the error event and removes all listeners.
     * @param error The error event received
     */
    public onError(error: any): void {
        (this as EventEmitter).emit(ERROR, error);
        (this as EventEmitter).emit(SOCKET_ERROR, error);
        this.removeAllSocketListeners();
    }
    /**
     * Emits the close event and removes all listeners.
     * @param error Error received on close (if any)
     */
    public onClose(error: CloseEvent|null = null): void {
        (this as EventEmitter).emit(CLOSE, error);
        (this as EventEmitter).emit(SOCKET_CLOSE, error);
        this.removeAllSocketListeners();
        (this as EventEmitter).removeAllListeners();
    }
    /**
     * Emits the connect event and checks if there are subscriptions defined that should be resubscribed.
     */
    public async onConnectAsync(): Promise<void> {
        const subscriptionKeys = Object.keys(this.subscriptions);

        if (subscriptionKeys.length > 0) {
            let subscriptionId;

            for (const key of subscriptionKeys) {
                subscriptionId = await this.subscribeAsync(
                    this.subscriptions[key].subscribeMethod,
                    this.subscriptions[key].parameters[0],
                    this.subscriptions[key].parameters.slice(1),
                );

                if (key !== subscriptionId) {
                    delete this.subscriptions[subscriptionId];
                }

                this.subscriptions[key].id = subscriptionId;
            }
        }

        (this as EventEmitter).emit(SOCKET_CONNECT);
        (this as EventEmitter).emit(CONNECT);
    }
    /**
     * This is the listener for the 'message' events of the current socket connection.
     * @param response message returned via the socket connection
     */
    public onMessage(response: string|object): void {
        let event;

        let responseObject = response as any;
        if (typeof response !== 'object') {
            responseObject = JSON.parse(response);
        }

        if (Array.isArray(responseObject)) {
            event = responseObject[0].id;
        } else if (typeof responseObject.id === 'undefined') {
            event = this.getSubscriptionEvent(responseObject.params.subscription);
            responseObject = responseObject.params;
        } else {
            event = responseObject.id;
        }

        (this as EventEmitter).emit(SOCKET_MESSAGE, responseObject);
        (this as EventEmitter).emit(event, responseObject);
    }
    /**
     * Resets the providers, clears all callbacks
     */
    public reset(): void {
        (this as EventEmitter).removeAllListeners();
        this.registerEventListeners();
    }
    /**
     * Subscribes to a given subscriptionType
     * @param subscribeMethod JSON-RPC method name to use for subscription
     * @param subscriptionMethod Subscription namespace
     * @param parameters Additional parameters to subscribe call
     * @returns The subscriptionId of an error
     */
    public subscribeAsync(subscribeMethod: string, subscriptionMethod: string, parameters: any[]): Promise<any> {
        parameters.unshift(subscriptionMethod);

        return this.sendAsync(subscribeMethod, parameters)
            .then(subscriptionId => {
                this.subscriptions[subscriptionId] = {
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
    public unsubscribeAsync(subscriptionId: string, unsubscribeMethod: string): Promise<boolean|Error> {
        if (this.hasSubscription(subscriptionId)) {
            return this.sendAsync(unsubscribeMethod, [subscriptionId]).then(response => {
                if (response) {
                    (this as EventEmitter).removeAllListeners(this.getSubscriptionEvent(subscriptionId));

                    delete this.subscriptions[subscriptionId];
                }

                return response;
            });
        }

        return Promise.reject(new Error(`Provider error: Subscription with ID ${subscriptionId} does not exist.`));
    }
    /**
     * Clears all subscriptions and listeners
     * @param unsubscribeMethod JSON-RPC unsubscribe method
     * @returns true if clearing subscription succeeds, otherwise an error
     */
    public async clearSubscriptionsAsync(unsubscribeMethod: string): Promise<boolean|Error> {
        const unsubscribePromises: Array<Promise<any>> = [];

        Object.keys(this.subscriptions).forEach(key => {
            (this as EventEmitter).removeAllListeners(key);
            unsubscribePromises.push(this.unsubscribeAsync(this.subscriptions[key].id, unsubscribeMethod));
        });

        return Promise.all(unsubscribePromises).then(results => {
            if (results.includes(false)) {
                throw new Error(`Could not unsubscribe all subscriptions: ${JSON.stringify(results)}`);
            }

            return true;
        });
    }
    /**
     * Checks if the given subscription id exists
     * @param subscriptionId subscription ID to check existence for
     * @returns whether or not the subscription exists
     */
    public hasSubscription(subscriptionId: string): boolean {
        return typeof this.getSubscriptionEvent(subscriptionId) !== 'undefined';
    }
    /**
     * Returns the event the subscription is listening for.
     * @param subscriptionId subscription ID
     * @returns subscription event name (e.g. "heartbeat")
     */
    public getSubscriptionEvent(subscriptionId: string): string|undefined {
        if (this.subscriptions[subscriptionId]) {
            return subscriptionId;
        }

        let event: string|undefined;
        for (const key in this.subscriptions) {
            if (this.subscriptions[key].id === subscriptionId) {
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
    public toPayload(method: string, params: any[]): any {
        if (!method) {
            throw new Error(`JSONRPC method should be specified for params: "${JSON.stringify(params)}"!`);
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
// tslint:enable:prefer-function-over-method
