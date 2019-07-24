"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
var __spread = (this && this.__spread) || function () {
    for (var ar = [], i = 0; i < arguments.length; i++) ar = ar.concat(__read(arguments[i]));
    return ar;
};
Object.defineProperty(exports, "__esModule", { value: true });
var assert_1 = require("@0x/assert");
var order_utils_1 = require("@0x/order-utils");
var utils_1 = require("@0x/utils");
var Web3Providers = require("@0x/web3-providers-fork");
var uuid_1 = require("uuid");
var WebSocket = require("websocket");
var CLOSE_REASON_NO_HEARTBEAT = 3001;
var CLOSE_DESCRIPTION_NO_HEARTBEAT = 'No heartbeat received';
var DEFAULT_RECONNECT_AFTER = 5000;
var DEFAULT_WS_OPTS = {
    clientConfig: {
        // For some reason fragmenting the payloads causes the connection to close
        // Source: https://github.com/theturtle32/WebSocket-Node/issues/359
        fragmentOutgoingMessages: false,
    },
    reconnectAfter: DEFAULT_RECONNECT_AFTER,
};
/**
 * This class includes all the functionality related to interacting with a Mesh JSON RPC
 * websocket endpoint.
 */
var WSClient = /** @class */ (function () {
    /**
     * Instantiates a new WSClient instance
     * @param   url               WS server endpoint
     * @param   wsOpts            WebSocket options
     * @return  An instance of WSClient
     */
    function WSClient(url, wsOpts) {
        this._subscriptionIdToMeshSpecificId = {};
        if (wsOpts !== undefined && wsOpts.reconnectAfter === undefined) {
            wsOpts.reconnectAfter = DEFAULT_RECONNECT_AFTER;
        }
        this._wsProvider = new Web3Providers.WebsocketProvider(url, wsOpts !== undefined ? wsOpts : DEFAULT_WS_OPTS);
        // Intentional fire-and-forget
        // tslint:disable-next-line:no-floating-promises
        this._startInternalLivenessCheckAsync();
    }
    WSClient._convertRawAcceptedOrders = function (rawAcceptedOrders) {
        var acceptedOrderInfos = [];
        rawAcceptedOrders.forEach(function (rawAcceptedOrderInfo) {
            var acceptedOrderInfo = {
                orderHash: rawAcceptedOrderInfo.orderHash,
                signedOrder: order_utils_1.orderParsingUtils.convertOrderStringFieldsToBigNumber(rawAcceptedOrderInfo.signedOrder),
                fillableTakerAssetAmount: new utils_1.BigNumber(rawAcceptedOrderInfo.fillableTakerAssetAmount),
            };
            acceptedOrderInfos.push(acceptedOrderInfo);
        });
        return acceptedOrderInfos;
    };
    /**
     * Adds an array of 0x signed orders to the Mesh node.
     * @param signedOrders signedOrders to add
     * @returns validation results
     */
    WSClient.prototype.addOrdersAsync = function (signedOrders) {
        return __awaiter(this, void 0, void 0, function () {
            var rawValidationResults, validationResults;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        assert_1.assert.isArray('signedOrders', signedOrders);
                        return [4 /*yield*/, this._wsProvider.send('mesh_addOrders', [
                                signedOrders,
                            ])];
                    case 1:
                        rawValidationResults = _a.sent();
                        validationResults = {
                            accepted: WSClient._convertRawAcceptedOrders(rawValidationResults.accepted),
                            rejected: [],
                        };
                        rawValidationResults.rejected.forEach(function (rawRejectedOrderInfo) {
                            var rejectedOrderInfo = {
                                orderHash: rawRejectedOrderInfo.orderHash,
                                signedOrder: order_utils_1.orderParsingUtils.convertOrderStringFieldsToBigNumber(rawRejectedOrderInfo.signedOrder),
                                kind: rawRejectedOrderInfo.kind,
                                status: rawRejectedOrderInfo.status,
                            };
                            validationResults.rejected.push(rejectedOrderInfo);
                        });
                        return [2 /*return*/, validationResults];
                }
            });
        });
    };
    /**
     * Get all 0x signed orders currently stored in the Mesh node
     * @param perPage number of signedOrders to fetch per paginated request
     * @returns all orders, their hash and their fillableTakerAssetAmount
     */
    WSClient.prototype.getOrdersAsync = function (perPage) {
        if (perPage === void 0) { perPage = 200; }
        return __awaiter(this, void 0, void 0, function () {
            var snapshotID, ordersInfosLen, i, rawAcceptedOrderInfos, page, getOrdersResponse, allOrdersInfos;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        snapshotID = '';
                        ordersInfosLen = 1;
                        i = 0;
                        rawAcceptedOrderInfos = [];
                        _a.label = 1;
                    case 1:
                        if (!(ordersInfosLen === perPage)) return [3 /*break*/, 3];
                        page = i;
                        return [4 /*yield*/, this._wsProvider.send('mesh_getOrders', [
                                page,
                                perPage,
                                snapshotID,
                            ])];
                    case 2:
                        getOrdersResponse = _a.sent();
                        snapshotID = getOrdersResponse.snapshotID;
                        rawAcceptedOrderInfos = __spread(rawAcceptedOrderInfos, getOrdersResponse.ordersInfos);
                        ordersInfosLen = getOrdersResponse.ordersInfos.length;
                        i++;
                        return [3 /*break*/, 1];
                    case 3:
                        allOrdersInfos = WSClient._convertRawAcceptedOrders(rawAcceptedOrderInfos);
                        return [2 /*return*/, allOrdersInfos];
                }
            });
        });
    };
    /**
     * Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a
     * subscriptionId that can be used to `unsubscribe()` from this subscription.
     * @param   cb   callback function where you'd like to get notified about order events
     * @return subscriptionId
     */
    WSClient.prototype.subscribeToOrdersAsync = function (cb) {
        return __awaiter(this, void 0, void 0, function () {
            var orderEventsSubscriptionId, id, orderEventsCallback;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        assert_1.assert.isFunction('cb', cb);
                        return [4 /*yield*/, this._wsProvider.subscribe('mesh_subscribe', 'orders', [])];
                    case 1:
                        orderEventsSubscriptionId = _a.sent();
                        id = uuid_1.v4();
                        this._subscriptionIdToMeshSpecificId[id] = orderEventsSubscriptionId;
                        orderEventsCallback = function (eventPayload) {
                            _this._subscriptionIdToMeshSpecificId[id] = eventPayload.subscription;
                            var rawOrderEvents = eventPayload.result;
                            var orderEvents = [];
                            rawOrderEvents.forEach(function (rawOrderEvent) {
                                var orderEvent = {
                                    orderHash: rawOrderEvent.orderHash,
                                    signedOrder: order_utils_1.orderParsingUtils.convertOrderStringFieldsToBigNumber(rawOrderEvent.signedOrder),
                                    kind: rawOrderEvent.kind,
                                    fillableTakerAssetAmount: new utils_1.BigNumber(rawOrderEvent.fillableTakerAssetAmount),
                                    txHashes: rawOrderEvent.txHashes,
                                };
                                orderEvents.push(orderEvent);
                            });
                            cb(orderEvents);
                        };
                        this._wsProvider.on(orderEventsSubscriptionId, orderEventsCallback);
                        return [2 /*return*/, id];
                }
            });
        });
    };
    /**
     * Subscribe to the 'heartbeat' topic and receive an ack from the Mesh every 5 seconds. This method
     * returns a subscriptionId that can be used to `unsubscribe()` from this subscription.
     * @param   cb   callback function where you'd like to get notified about heartbeats
     * @return subscriptionId
     */
    WSClient.prototype.subscribeToHeartbeatAsync = function (cb) {
        return __awaiter(this, void 0, void 0, function () {
            var heartbeatSubscriptionId, id, orderEventsCallback;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        assert_1.assert.isFunction('cb', cb);
                        return [4 /*yield*/, this._wsProvider.subscribe('mesh_subscribe', 'heartbeat', [])];
                    case 1:
                        heartbeatSubscriptionId = _a.sent();
                        id = uuid_1.v4();
                        this._subscriptionIdToMeshSpecificId[id] = heartbeatSubscriptionId;
                        orderEventsCallback = function (eventPayload) {
                            _this._subscriptionIdToMeshSpecificId[id] = eventPayload.subscription;
                            cb(eventPayload.result);
                        };
                        this._wsProvider.on(heartbeatSubscriptionId, orderEventsCallback);
                        return [2 /*return*/, id];
                }
            });
        });
    };
    /**
     * Unsubscribe from a subscription
     * @param subscriptionId identifier of the subscription to cancel
     */
    WSClient.prototype.unsubscribeAsync = function (subscriptionId) {
        return __awaiter(this, void 0, void 0, function () {
            var meshSubscriptionId;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        assert_1.assert.isString('subscriptionId', subscriptionId);
                        meshSubscriptionId = this._subscriptionIdToMeshSpecificId[subscriptionId];
                        return [4 /*yield*/, this._wsProvider.send('mesh_unsubscribe', [meshSubscriptionId])];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        });
    };
    /**
     * Get notified when the underlying WS connection closes normally. If it closes with an
     * error, WSClient automatically attempts to re-connect without emitting a `close` event.
     * @param cb callback to call when WS connection closes
     */
    WSClient.prototype.onClose = function (cb) {
        this._wsProvider.connection.addEventListener('close', function () {
            cb();
        });
    };
    /**
     * Get notified when a connection to the underlying WS connection is re-established
     * @param cb callback to call with the error when it occurs
     */
    WSClient.prototype.onReconnected = function (cb) {
        this._wsProvider.on('reconnected', function () {
            cb();
        });
    };
    /**
     * destroy unsubscribes all active subscriptions, closes the websocket connection
     * and stops the internal heartbeat connection liveness check.
     */
    WSClient.prototype.destroy = function () {
        clearInterval(this._heartbeatCheckIntervalId);
        this._wsProvider.clearSubscriptions('mesh_unsubscribe');
        this._wsProvider.removeAllListeners();
        this._wsProvider.disconnect(WebSocket.connection.CLOSE_REASON_NORMAL, 'Normal connection closure');
    };
    WSClient.prototype._startInternalLivenessCheckAsync = function () {
        return __awaiter(this, void 0, void 0, function () {
            var lastHeartbeatTimestampMs, err_1, oneSecondInMs;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        lastHeartbeatTimestampMs = new Date().getTime();
                        _a.label = 1;
                    case 1:
                        _a.trys.push([1, 3, , 4]);
                        return [4 /*yield*/, this.subscribeToHeartbeatAsync(function (ack) {
                                lastHeartbeatTimestampMs = new Date().getTime();
                            })];
                    case 2:
                        _a.sent();
                        return [3 /*break*/, 4];
                    case 3:
                        err_1 = _a.sent();
                        throw new Error('Failed to establish under-the-hood heartbeat subscription');
                    case 4:
                        oneSecondInMs = 1000;
                        this._heartbeatCheckIntervalId = setInterval(function () {
                            var twentySecondsInMs = 20000;
                            var haveTwentySecondsPastWithoutAHeartBeat = lastHeartbeatTimestampMs + twentySecondsInMs < new Date().getTime();
                            if (haveTwentySecondsPastWithoutAHeartBeat) {
                                // If connected, we haven't received a heartbeat in over 20 seconds, re-connect
                                if (_this._wsProvider.connected) {
                                    _this._wsProvider.disconnect(CLOSE_REASON_NO_HEARTBEAT, CLOSE_DESCRIPTION_NO_HEARTBEAT);
                                }
                                lastHeartbeatTimestampMs = new Date().getTime();
                            }
                        }, oneSecondInMs);
                        return [2 /*return*/];
                }
            });
        });
    };
    return WSClient;
}());
exports.WSClient = WSClient;
