import { assert } from '@0x/assert';
import { ObjectMap, SignedOrder } from '@0x/types';
import { BigNumber } from '@0x/utils';
import { v4 as uuid } from 'uuid';
import * as Web3Providers from 'web3-providers';
import * as WebSocket from 'websocket';

import {
    AcceptedOrderInfo,
    ContractEvent,
    ContractEventKind,
    ContractEventParameters,
    ERC1155ApprovalForAllEvent,
    ERC721ApprovalForAllEvent,
    ExchangeCancelEvent,
    GetOrdersResponse,
    GetStatsResponse,
    HeartbeatEventPayload,
    OrderEvent,
    OrderEventPayload,
    OrderInfo,
    RawAcceptedOrderInfo,
    RawGetOrdersResponse,
    RawOrderEvent,
    RawOrderInfo,
    RawValidationResults,
    RejectedOrderInfo,
    StringifiedContractEvent,
    StringifiedERC1155TransferBatchEvent,
    StringifiedERC1155TransferSingleEvent,
    StringifiedERC20ApprovalEvent,
    StringifiedERC20TransferEvent,
    StringifiedERC721ApprovalEvent,
    StringifiedERC721TransferEvent,
    StringifiedExchangeCancelUpToEvent,
    StringifiedExchangeFillEvent,
    StringifiedWethDepositEvent,
    StringifiedWethWithdrawalEvent,
    ValidationResults,
    WSOpts,
} from './types';

const CLOSE_REASON_NO_HEARTBEAT = 3001;
const CLOSE_DESCRIPTION_NO_HEARTBEAT = 'No heartbeat received';

const DEFAULT_RECONNECT_AFTER_MS = 5000;
const DEFAULT_RPC_REQUEST_TIMEOUT = 30000;
const DEFAULT_WS_OPTS = {
    clientConfig: {
        // For some reason fragmenting the payloads causes the connection to close
        // Source: https://github.com/theturtle32/WebSocket-Node/issues/359
        fragmentOutgoingMessages: false,
    },
    timeout: DEFAULT_RPC_REQUEST_TIMEOUT,
    reconnectDelay: DEFAULT_RECONNECT_AFTER_MS,
};

/**
 * This class includes all the functionality related to interacting with a Mesh JSON RPC
 * websocket endpoint.
 */
export class WSClient {
    private readonly _subscriptionIdToMeshSpecificId: ObjectMap<string>;
    private _heartbeatCheckIntervalId: number | undefined;
    private readonly _wsProvider: Web3Providers.WebsocketProvider;
    private static _convertRawAcceptedOrderInfos(rawAcceptedOrderInfos: RawAcceptedOrderInfo[]): AcceptedOrderInfo[] {
        const acceptedOrderInfos: AcceptedOrderInfo[] = [];
        rawAcceptedOrderInfos.forEach(rawAcceptedOrderInfo => {
            const acceptedOrderInfo: AcceptedOrderInfo = {
                orderHash: rawAcceptedOrderInfo.orderHash,
                signedOrder: WSClient._convertOrderStringFieldsToBigNumber(rawAcceptedOrderInfo.signedOrder),
                fillableTakerAssetAmount: new BigNumber(rawAcceptedOrderInfo.fillableTakerAssetAmount),
                isNew: rawAcceptedOrderInfo.isNew,
            };
            acceptedOrderInfos.push(acceptedOrderInfo);
        });
        return acceptedOrderInfos;
    }
    private static _convertRawOrderInfos(rawOrderInfos: RawOrderInfo[]): OrderInfo[] {
        const orderInfos: OrderInfo[] = [];
        rawOrderInfos.forEach(rawOrderInfo => {
            const orderInfo: OrderInfo = {
                orderHash: rawOrderInfo.orderHash,
                signedOrder: WSClient._convertOrderStringFieldsToBigNumber(rawOrderInfo.signedOrder),
                fillableTakerAssetAmount: new BigNumber(rawOrderInfo.fillableTakerAssetAmount),
            };
            orderInfos.push(orderInfo);
        });
        return orderInfos;
    }
    private static _convertStringsFieldsToBigNumbers(obj: any, fields: string[]): any {
        const result = { ...obj };
        fields.forEach(field => {
            if (result[field] === undefined) {
                throw new Error(`Could not find field '${field}' while converting string fields to BigNumber.`);
            }
            result[field] = new BigNumber(result[field]);
        });
        return result;
    }
    private static _convertOrderStringFieldsToBigNumber(order: any): any {
        return WSClient._convertStringsFieldsToBigNumbers(order, [
            'makerAssetAmount',
            'takerAssetAmount',
            'makerFee',
            'takerFee',
            'expirationTimeSeconds',
            'salt',
        ]);
    }
    private static _convertRawGetOrdersResponse(rawGetOrdersResponse: RawGetOrdersResponse): GetOrdersResponse {
        return {
            snapshotID: rawGetOrdersResponse.snapshotID,
            // tslint:disable-next-line:custom-no-magic-numbers
            snapshotTimestamp: Math.round(new Date(rawGetOrdersResponse.snapshotTimestamp).getTime() / 1000),
            ordersInfos: WSClient._convertRawOrderInfos(rawGetOrdersResponse.ordersInfos),
        };
    }
    private static _convertStringifiedContractEvents(rawContractEvents: StringifiedContractEvent[]): ContractEvent[] {
        const contractEvents: ContractEvent[] = [];
        if (rawContractEvents === null) {
            return contractEvents;
        }
        rawContractEvents.forEach(rawContractEvent => {
            const kind = rawContractEvent.kind as ContractEventKind;
            const rawParameters = rawContractEvent.parameters;
            let parameters: ContractEventParameters;
            switch (kind) {
                case ContractEventKind.ERC20TransferEvent:
                    const erc20TransferEvent = rawParameters as StringifiedERC20TransferEvent;
                    parameters = {
                        from: erc20TransferEvent.from,
                        to: erc20TransferEvent.to,
                        value: new BigNumber(erc20TransferEvent.value),
                    };
                    break;
                case ContractEventKind.ERC20ApprovalEvent:
                    const erc20ApprovalEvent = rawParameters as StringifiedERC20ApprovalEvent;
                    parameters = {
                        owner: erc20ApprovalEvent.owner,
                        spender: erc20ApprovalEvent.spender,
                        value: new BigNumber(erc20ApprovalEvent.value),
                    };
                    break;
                case ContractEventKind.ERC721TransferEvent:
                    const erc721TransferEvent = rawParameters as StringifiedERC721TransferEvent;
                    parameters = {
                        from: erc721TransferEvent.from,
                        to: erc721TransferEvent.to,
                        tokenId: new BigNumber(erc721TransferEvent.tokenId),
                    };
                    break;
                case ContractEventKind.ERC721ApprovalEvent:
                    const erc721ApprovalEvent = rawParameters as StringifiedERC721ApprovalEvent;
                    parameters = {
                        owner: erc721ApprovalEvent.owner,
                        approved: erc721ApprovalEvent.approved,
                        tokenId: new BigNumber(erc721ApprovalEvent.tokenId),
                    };
                    break;
                case ContractEventKind.ERC721ApprovalForAllEvent:
                    parameters = rawParameters as ERC721ApprovalForAllEvent;
                    break;
                case ContractEventKind.ERC1155ApprovalForAllEvent:
                    parameters = rawParameters as ERC1155ApprovalForAllEvent;
                    break;
                case ContractEventKind.ERC1155TransferSingleEvent:
                    const erc1155TransferSingleEvent = rawParameters as StringifiedERC1155TransferSingleEvent;
                    parameters = {
                        operator: erc1155TransferSingleEvent.operator,
                        from: erc1155TransferSingleEvent.from,
                        to: erc1155TransferSingleEvent.to,
                        id: new BigNumber(erc1155TransferSingleEvent.id),
                        value: new BigNumber(erc1155TransferSingleEvent.value),
                    };
                    break;
                case ContractEventKind.ERC1155TransferBatchEvent:
                    const erc1155TransferBatchEvent = rawParameters as StringifiedERC1155TransferBatchEvent;
                    const ids: BigNumber[] = [];
                    erc1155TransferBatchEvent.ids.forEach(id => {
                        ids.push(new BigNumber(id));
                    });
                    const values: BigNumber[] = [];
                    erc1155TransferBatchEvent.values.forEach(value => {
                        values.push(new BigNumber(value));
                    });
                    parameters = {
                        operator: erc1155TransferBatchEvent.operator,
                        from: erc1155TransferBatchEvent.from,
                        to: erc1155TransferBatchEvent.to,
                        ids,
                        values,
                    };
                    break;
                case ContractEventKind.ExchangeFillEvent:
                    const exchangeFillEvent = rawParameters as StringifiedExchangeFillEvent;
                    parameters = {
                        makerAddress: exchangeFillEvent.makerAddress,
                        takerAddress: exchangeFillEvent.takerAddress,
                        senderAddress: exchangeFillEvent.senderAddress,
                        feeRecipientAddress: exchangeFillEvent.feeRecipientAddress,
                        makerAssetFilledAmount: new BigNumber(exchangeFillEvent.makerAssetFilledAmount),
                        takerAssetFilledAmount: new BigNumber(exchangeFillEvent.takerAssetFilledAmount),
                        makerFeePaid: new BigNumber(exchangeFillEvent.makerFeePaid),
                        takerFeePaid: new BigNumber(exchangeFillEvent.takerFeePaid),
                        orderHash: exchangeFillEvent.orderHash,
                        makerAssetData: exchangeFillEvent.makerAssetData,
                        takerAssetData: exchangeFillEvent.takerAssetData,
                    };
                    break;
                case ContractEventKind.ExchangeCancelEvent:
                    parameters = rawParameters as ExchangeCancelEvent;
                    break;
                case ContractEventKind.ExchangeCancelUpToEvent:
                    const exchangeCancelUpToEvent = rawParameters as StringifiedExchangeCancelUpToEvent;
                    parameters = {
                        makerAddress: exchangeCancelUpToEvent.makerAddress,
                        senderAddress: exchangeCancelUpToEvent.senderAddress,
                        orderEpoch: new BigNumber(exchangeCancelUpToEvent.orderEpoch),
                    };
                    break;
                case ContractEventKind.WethDepositEvent:
                    const wethDepositEvent = rawParameters as StringifiedWethDepositEvent;
                    parameters = {
                        owner: wethDepositEvent.owner,
                        value: new BigNumber(wethDepositEvent.value),
                    };
                    break;
                case ContractEventKind.WethWithdrawalEvent:
                    const wethWithdrawalEvent = rawParameters as StringifiedWethWithdrawalEvent;
                    parameters = {
                        owner: wethWithdrawalEvent.owner,
                        value: new BigNumber(wethWithdrawalEvent.value),
                    };
                    break;
                default:
                    throw new Error(`Unrecognized ContractEventKind: ${kind}`);
            }
            const contractEvent: ContractEvent = {
                blockHash: rawContractEvent.blockHash,
                txHash: rawContractEvent.txHash,
                txIndex: rawContractEvent.txIndex,
                logIndex: rawContractEvent.logIndex,
                isRemoved: rawContractEvent.isRemoved,
                address: rawContractEvent.address,
                kind,
                parameters,
            };
            contractEvents.push(contractEvent);
        });
        return contractEvents;
    }
    /**
     * Instantiates a new WSClient instance
     * @param   url               WS server endpoint
     * @param   wsOpts            WebSocket options
     * @return  An instance of WSClient
     */
    constructor(url: string, wsOpts?: WSOpts) {
        this._subscriptionIdToMeshSpecificId = {};
        if (wsOpts !== undefined && wsOpts.reconnectDelay === undefined) {
            wsOpts.reconnectDelay = DEFAULT_RECONNECT_AFTER_MS;
        }
        this._wsProvider = new Web3Providers.WebsocketProvider(
            url,
            wsOpts !== undefined ? (wsOpts as any) : DEFAULT_WS_OPTS,
        );
        // Intentional fire-and-forget
        // tslint:disable-next-line:no-floating-promises
        this._startInternalLivenessCheckAsync();
    }
    /**
     * Adds an array of 0x signed orders to the Mesh node.
     * @param signedOrders signedOrders to add
     * @param pinned       Whether or not the orders should be pinned. Pinned
     * orders will not be affected by any DDoS prevention or incentive
     * mechanisms and will always stay in storage until they are no longer
     * fillable.
     * @returns validation results
     */
    public async addOrdersAsync(signedOrders: SignedOrder[], pinned: boolean = true): Promise<ValidationResults> {
        assert.isArray('signedOrders', signedOrders);
        const rawValidationResults: RawValidationResults = await this._wsProvider.send('mesh_addOrders', [
            signedOrders,
            { pinned },
        ]);
        const validationResults: ValidationResults = {
            accepted: WSClient._convertRawAcceptedOrderInfos(rawValidationResults.accepted),
            rejected: [],
        };
        rawValidationResults.rejected.forEach(rawRejectedOrderInfo => {
            const rejectedOrderInfo: RejectedOrderInfo = {
                orderHash: rawRejectedOrderInfo.orderHash,
                signedOrder: WSClient._convertOrderStringFieldsToBigNumber(rawRejectedOrderInfo.signedOrder),
                kind: rawRejectedOrderInfo.kind,
                status: rawRejectedOrderInfo.status,
            };
            validationResults.rejected.push(rejectedOrderInfo);
        });
        return validationResults;
    }
    public async getStatsAsync(): Promise<GetStatsResponse> {
        const stats = await this._wsProvider.send('mesh_getStats', []);
        return stats;
    }
    /**
     * Get all 0x signed orders currently stored in the Mesh node
     * @param perPage number of signedOrders to fetch per paginated request
     * @returns the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts
     */
    public async getOrdersAsync(perPage: number = 200): Promise<GetOrdersResponse> {
        let snapshotID = ''; // New snapshot

        let page = 0;
        let getOrdersResponse = await this.getOrdersForPageAsync(page, perPage, snapshotID);
        snapshotID = getOrdersResponse.snapshotID;
        let ordersInfos = getOrdersResponse.ordersInfos;

        let allOrderInfos: OrderInfo[] = [];

        do {
            allOrderInfos = [...allOrderInfos, ...ordersInfos];
            page++;
            getOrdersResponse = await this.getOrdersForPageAsync(page, perPage, snapshotID);
            ordersInfos = getOrdersResponse.ordersInfos;
        } while (ordersInfos.length > 0);

        getOrdersResponse = {
            snapshotID,
            snapshotTimestamp: getOrdersResponse.snapshotTimestamp,
            ordersInfos: allOrderInfos,
        };
        return getOrdersResponse;
    }
    /**
     * Get page of 0x signed orders stored on the Mesh node at the specified snapshot
     * @param page Page index at which to retrieve orders
     * @param perPage number of signedOrders to fetch per paginated request
     * @param snapshotID The DB snapshot at which to fetch orders. If omitted, a new snapshot is created
     * @returns the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts
     */
    public async getOrdersForPageAsync(
        page: number,
        perPage: number = 200,
        snapshotID?: string,
    ): Promise<GetOrdersResponse> {
        const finalSnapshotID = snapshotID === undefined ? '' : snapshotID;

        const rawGetOrdersResponse: RawGetOrdersResponse = await this._wsProvider.send('mesh_getOrders', [
            page,
            perPage,
            finalSnapshotID,
        ]);
        const getOrdersResponse = WSClient._convertRawGetOrdersResponse(rawGetOrdersResponse);
        return getOrdersResponse;
    }
    /**
     * Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a
     * subscriptionId that can be used to `unsubscribe()` from this subscription.
     * @param   cb   callback function where you'd like to get notified about order events
     * @return subscriptionId
     */
    public async subscribeToOrdersAsync(cb: (orderEvents: OrderEvent[]) => void): Promise<string> {
        assert.isFunction('cb', cb);
        const orderEventsSubscriptionId = await this._wsProvider.subscribe('mesh_subscribe', 'orders', []);
        const id = uuid();
        this._subscriptionIdToMeshSpecificId[id] = orderEventsSubscriptionId;

        const orderEventsCallback = (eventPayload: OrderEventPayload) => {
            this._subscriptionIdToMeshSpecificId[id] = eventPayload.subscription;
            const rawOrderEvents: RawOrderEvent[] = eventPayload.result;
            const orderEvents: OrderEvent[] = [];
            rawOrderEvents.forEach(rawOrderEvent => {
                const orderEvent = {
                    timestampMs: new Date(rawOrderEvent.timestamp).getTime(),
                    orderHash: rawOrderEvent.orderHash,
                    signedOrder: WSClient._convertOrderStringFieldsToBigNumber(rawOrderEvent.signedOrder),
                    endState: rawOrderEvent.endState,
                    fillableTakerAssetAmount: new BigNumber(rawOrderEvent.fillableTakerAssetAmount),
                    contractEvents: WSClient._convertStringifiedContractEvents(rawOrderEvent.contractEvents),
                };
                orderEvents.push(orderEvent);
            });
            cb(orderEvents);
        };
        this._wsProvider.on(orderEventsSubscriptionId, orderEventsCallback as any);
        return id;
    }
    /**
     * Unsubscribe from a subscription
     * @param subscriptionId identifier of the subscription to cancel
     */
    public async unsubscribeAsync(subscriptionId: string): Promise<void> {
        assert.isString('subscriptionId', subscriptionId);
        const meshSubscriptionId = this._subscriptionIdToMeshSpecificId[subscriptionId];
        await this._wsProvider.send('mesh_unsubscribe', [meshSubscriptionId]);
    }
    /**
     * Get notified when the underlying WS connection closes normally. If it closes with an
     * error, WSClient automatically attempts to re-connect without emitting a `close` event.
     * @param cb callback to call when WS connection closes
     */
    public onClose(cb: () => void): void {
        (this._wsProvider as any).connection.addEventListener('close', () => {
            cb();
        });
    }
    /**
     * Get notified when a connection to the underlying WS connection is re-established
     * @param cb callback to call with the error when it occurs
     */
    public onReconnected(cb: () => void): void {
        this._wsProvider.on('reconnected', () => {
            cb();
        });
    }
    /**
     * destroy unsubscribes all active subscriptions, closes the websocket connection
     * and stops the internal heartbeat connection liveness check.
     */
    public destroy(): void {
        clearInterval(this._heartbeatCheckIntervalId);
        this._wsProvider.clearSubscriptions('mesh_unsubscribe');
        this._wsProvider.disconnect(WebSocket.connection.CLOSE_REASON_NORMAL, 'Normal connection closure');
    }
    /**
     * Subscribe to the 'heartbeat' topic and receive an ack from the Mesh every 5 seconds. This method
     * returns a subscriptionId that can be used to `unsubscribe()` from this subscription.
     * @param   cb   callback function where you'd like to get notified about heartbeats
     * @return subscriptionId
     */
    private async _subscribeToHeartbeatAsync(cb: (ack: string) => void): Promise<string> {
        assert.isFunction('cb', cb);
        const heartbeatSubscriptionId = await this._wsProvider.subscribe('mesh_subscribe', 'heartbeat', []);
        const id = uuid();
        this._subscriptionIdToMeshSpecificId[id] = heartbeatSubscriptionId;

        const orderEventsCallback = (eventPayload: HeartbeatEventPayload) => {
            this._subscriptionIdToMeshSpecificId[id] = eventPayload.subscription;
            cb(eventPayload.result);
        };
        this._wsProvider.on(heartbeatSubscriptionId, orderEventsCallback as any);
        return id;
    }
    private async _startInternalLivenessCheckAsync(): Promise<void> {
        let lastHeartbeatTimestampMs = new Date().getTime();
        try {
            await this._subscribeToHeartbeatAsync((ack: string) => {
                lastHeartbeatTimestampMs = new Date().getTime();
            });
        } catch (err) {
            throw new Error('Failed to establish under-the-hood heartbeat subscription');
        }
        const oneSecondInMs = 1000;
        this._heartbeatCheckIntervalId = setInterval(() => {
            const twentySecondsInMs = 20000;
            // tslint:disable-next-line:boolean-naming
            const haveTwentySecondsPastWithoutAHeartBeat =
                lastHeartbeatTimestampMs + twentySecondsInMs < new Date().getTime();
            if (haveTwentySecondsPastWithoutAHeartBeat) {
                // If connected, we haven't received a heartbeat in over 20 seconds, re-connect
                if (this._wsProvider.connected) {
                    this._wsProvider.disconnect(CLOSE_REASON_NO_HEARTBEAT, CLOSE_DESCRIPTION_NO_HEARTBEAT);
                }
                lastHeartbeatTimestampMs = new Date().getTime();
            }
        }, oneSecondInMs) as any;
    }
}
