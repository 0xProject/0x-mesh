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
var _this = this;
Object.defineProperty(exports, "__esModule", { value: true });
var dev_utils_1 = require("@0x/dev-utils");
var chai = require("chai");
require("mocha");
var WebSocket = require("websocket");
var index_1 = require("../src/index");
var chai_setup_1 = require("./utils/chai_setup");
var mock_ws_server_1 = require("./utils/mock_ws_server");
chai_setup_1.chaiSetup.configure();
var expect = chai.expect;
describe('WSClient', function () {
    afterEach(function () {
        mock_ws_server_1.stopServer();
    });
    describe('#getOrdersAsync', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('properly makes multiple paginated requests under-the-hood and returns all signedOrders', function () { return __awaiter(_this, void 0, void 0, function () {
                var wsServer, client, perPage, orderInfos;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                        case 1:
                            wsServer = _a.sent();
                            wsServer.on('connect', (function (connection) {
                                var requestNum = 0;
                                connection.on('message', (function (message) {
                                    var jsonRpcRequest = JSON.parse(message.utf8Data);
                                    var responses = [
                                        // Heartbeat subscription (under-the-hood)
                                        "\n                            {\n                                \"id\": \"" + jsonRpcRequest.id + "\",\n                                \"jsonrpc\": \"2.0\",\n                                \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                            }\n                        ",
                                        // First paginated request
                                        "\n                            {\n                                \"id\": \"" + jsonRpcRequest.id + "\",\n                                \"jsonrpc\": \"2.0\",\n                                \"result\": {\n                                    \"snapshotID\": \"123\",\n                                    \"ordersInfos\": [\n                                        {\n                                            \"orderHash\": \"0xa0fcb54919f0b3823aa14b3f511146f6ac087ab333a70f9b24bbb1ba657a4250\",\n                                            \"signedOrder\": {\n                                                \"makerAddress\": \"0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149\",\n                                                \"makerAssetData\": \"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498\",\n                                                \"makerAssetAmount\": \"1000000000000000000\",\n                                                \"makerFee\": \"0\",\n                                                \"takerAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"takerAssetData\": \"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2\",\n                                                \"takerAssetAmount\": \"10000000000000000000000\",\n                                                \"takerFee\": \"0\",\n                                                \"senderAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"exchangeAddress\": \"0x080bf510FCbF18b91105470639e9561022937712\",\n                                                \"feeRecipientAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"expirationTimeSeconds\": \"1586340602\",\n                                                \"salt\": \"41253767178111694375645046549067933145709740457131351457334397888365956743955\",\n                                                \"signature\": \"0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02\"\n                                            },\n                                            \"fillableTakerAssetAmount\": \"10000000000000000000000\"\n                                        }\n                                    ]\n                                }\n                            }\n                            ",
                                        // Second paginated request
                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": {\n                                        \"snapshotID\": \"123\",\n                                        \"ordersInfos\": []\n                                    }\n                                }\n                            ",
                                    ];
                                    connection.sendUTF(responses[requestNum]);
                                    requestNum++;
                                }));
                            }));
                            client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                            perPage = 1;
                            return [4 /*yield*/, client.getOrdersAsync(perPage)];
                        case 2:
                            orderInfos = _a.sent();
                            expect(orderInfos).to.have.length(1);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.makerAssetAmount)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.takerAssetAmount)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.makerFee)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.takerFee)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.salt)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(orderInfos[0].signedOrder.expirationTimeSeconds)).to.equal(true);
                            client.destroy();
                            return [2 /*return*/];
                    }
                });
            }); });
            return [2 /*return*/];
        });
    }); });
    describe('#addOrdersAsync', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('sends a mesh_addOrders request and converts all numerical fields to BigNumbers in returned signedOrders', function () { return __awaiter(_this, void 0, void 0, function () {
                var wsServer, signedOrders, client, validationResults;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                        case 1:
                            wsServer = _a.sent();
                            wsServer.on('connect', (function (connection) {
                                var requestNum = 0;
                                connection.on('message', (function (message) {
                                    var jsonRpcRequest = JSON.parse(message.utf8Data);
                                    var responses = [
                                        // Heartbeat subscription (under-the-hood)
                                        "\n                            {\n                                \"id\": \"" + jsonRpcRequest.id + "\",\n                                \"jsonrpc\": \"2.0\",\n                                \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                            }\n                        ",
                                        // mesh_addOrders response
                                        "\n                            {\n                                \"id\": \"" + jsonRpcRequest.id + "\",\n                                \"jsonrpc\": \"2.0\",\n                                \"result\": {\n                                    \"accepted\": [\n                                        {\n                                            \"orderHash\": \"0xa0fcb54919f0b3823aa14b3f511146f6ac087ab333a70f9b24bbb1ba657a4250\",\n                                            \"signedOrder\": {\n                                                \"makerAddress\": \"0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149\",\n                                                \"makerAssetData\": \"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498\",\n                                                \"makerAssetAmount\": \"1000000000000000000\",\n                                                \"makerFee\": \"0\",\n                                                \"takerAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"takerAssetData\": \"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2\",\n                                                \"takerAssetAmount\": \"10000000000000000000000\",\n                                                \"takerFee\": \"0\",\n                                                \"senderAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"exchangeAddress\": \"0x080bf510FCbF18b91105470639e9561022937712\",\n                                                \"feeRecipientAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                \"expirationTimeSeconds\": \"1586340602\",\n                                                \"salt\": \"41253767178111694375645046549067933145709740457131351457334397888365956743955\",\n                                                \"signature\": \"0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02\"\n                                            },\n                                            \"fillableTakerAssetAmount\": \"10000000000000000000000\"\n                                        }\n                                    ],\n                                    \"rejected\": []\n                                }\n                            }\n                        ",
                                    ];
                                    connection.sendUTF(responses[requestNum]);
                                    requestNum++;
                                }));
                            }));
                            signedOrders = [
                                {
                                    makerAddress: '0xa3ece5d5b6319fa785efc10d3112769a46c6e149',
                                    takerAddress: '0x0000000000000000000000000000000000000000',
                                    makerAssetAmount: new index_1.BigNumber('1000000000000000000'),
                                    takerAssetAmount: new index_1.BigNumber('10000000000000000000000'),
                                    expirationTimeSeconds: new index_1.BigNumber('1586340602'),
                                    makerFee: new index_1.BigNumber('0'),
                                    takerFee: new index_1.BigNumber('0'),
                                    feeRecipientAddress: '0x0000000000000000000000000000000000000000',
                                    senderAddress: '0x0000000000000000000000000000000000000000',
                                    salt: new index_1.BigNumber('41253767178111694375645046549067933145709740457131351457334397888365956743955'),
                                    makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
                                    takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
                                    exchangeAddress: '0x080bf510fcbf18b91105470639e9561022937712',
                                    signature: '0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02',
                                },
                            ];
                            client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                            return [4 /*yield*/, client.addOrdersAsync(signedOrders)];
                        case 2:
                            validationResults = _a.sent();
                            expect(typeof validationResults === 'object').to.equal(true);
                            expect(validationResults.accepted).to.have.length(1);
                            expect(validationResults.rejected).to.have.length(0);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.makerAssetAmount)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.takerAssetAmount)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.makerFee)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.takerFee)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.salt)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.expirationTimeSeconds)).to.equal(true);
                            expect(index_1.BigNumber.isBigNumber(validationResults.accepted[0].fillableTakerAssetAmount)).to.equal(true);
                            client.destroy();
                            return [2 /*return*/];
                    }
                });
            }); });
            return [2 /*return*/];
        });
    }); });
    describe('#subscribeToHeartbeatAsync', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('should receive subscription updates', function (done) {
                (function () { return __awaiter(_this, void 0, void 0, function () {
                    var wsServer, client, expectToBeCalledOnce, callback;
                    var _this = this;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                            case 1:
                                wsServer = _a.sent();
                                wsServer.on('connect', (function (connection) {
                                    var requestNum = 0;
                                    connection.on('message', (function (message) { return __awaiter(_this, void 0, void 0, function () {
                                        var jsonRpcRequest, responses, eventResponse;
                                        return __generator(this, function (_a) {
                                            switch (_a.label) {
                                                case 0:
                                                    jsonRpcRequest = JSON.parse(message.utf8Data);
                                                    responses = [
                                                        // Heartbeat subscription (under-the-hood)
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                                }\n                            ",
                                                        // Requested heartbeat subscription
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                                }\n                            ",
                                                        // Response to unsubscribe
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\":\"2.0\",\n                                    \"result\":true,\n                                }\n                            ",
                                                    ];
                                                    connection.sendUTF(responses[requestNum]);
                                                    requestNum++;
                                                    if (!(requestNum === 2)) return [3 /*break*/, 2];
                                                    // tslint:disable-next-line:custom-no-magic-numbers
                                                    return [4 /*yield*/, sleepAsync(100)];
                                                case 1:
                                                    // tslint:disable-next-line:custom-no-magic-numbers
                                                    _a.sent();
                                                    eventResponse = "\n                                {\n                                    \"jsonrpc\":\"2.0\",\n                                    \"method\":\"mesh_subscription\",\n                                    \"params\": {\n                                        \"subscription\":\"0xab1a3e8af590364c09d0fa6a12103ada\",\n                                        \"result\":\"tick\"\n                                    }\n                                }\n                            ";
                                                    connection.sendUTF(eventResponse);
                                                    _a.label = 2;
                                                case 2: return [2 /*return*/];
                                            }
                                        });
                                    }); }));
                                }));
                                client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                                expectToBeCalledOnce = true;
                                callback = dev_utils_1.callbackErrorReporter.reportNoErrorCallbackErrors(done, expectToBeCalledOnce)(function (ack) { return __awaiter(_this, void 0, void 0, function () {
                                    return __generator(this, function (_a) {
                                        expect(ack).to.be.equal('tick');
                                        client.destroy();
                                        return [2 /*return*/];
                                    });
                                }); });
                                return [4 /*yield*/, client.subscribeToHeartbeatAsync(callback)];
                            case 2:
                                _a.sent();
                                return [2 /*return*/];
                        }
                    });
                }); })().catch(done);
            });
            return [2 /*return*/];
        });
    }); });
    describe('#subscribeToOrdersAsync', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('should receive subscription updates', function (done) {
                (function () { return __awaiter(_this, void 0, void 0, function () {
                    var wsServer, client, expectToBeCalledOnce, callback;
                    var _this = this;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                            case 1:
                                wsServer = _a.sent();
                                wsServer.on('connect', (function (connection) {
                                    var requestNum = 0;
                                    connection.on('message', (function (message) { return __awaiter(_this, void 0, void 0, function () {
                                        var jsonRpcRequest, responses, eventResponse;
                                        return __generator(this, function (_a) {
                                            switch (_a.label) {
                                                case 0:
                                                    jsonRpcRequest = JSON.parse(message.utf8Data);
                                                    responses = [
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                                }\n                            ",
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xc2ba3e8af590364c09d0fa6a12103adb\"\n                                }\n                            ",
                                                        // Response to unsubscribe
                                                        "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\":\"2.0\",\n                                    \"result\":true,\n                                }\n                            ",
                                                    ];
                                                    connection.sendUTF(responses[requestNum]);
                                                    requestNum++;
                                                    if (!(requestNum === 2)) return [3 /*break*/, 2];
                                                    // tslint:disable-next-line:custom-no-magic-numbers
                                                    return [4 /*yield*/, sleepAsync(100)];
                                                case 1:
                                                    // tslint:disable-next-line:custom-no-magic-numbers
                                                    _a.sent();
                                                    eventResponse = "\n                                {\n                                    \"jsonrpc\":\"2.0\",\n                                    \"method\":\"mesh_subscription\",\n                                    \"params\": {\n                                        \"subscription\":\"0xc2ba3e8af590364c09d0fa6a12103adb\",\n                                        \"result\": [\n                                            {\n                                                \"orderHash\": \"0x96e6eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ecc13fd4\",\n                                                \"signedOrder\": {\n                                                    \"makerAddress\": \"0x50f84bbee6fb250d6f49e854fa280445369d64d9\",\n                                                    \"makerAssetData\": \"0xf47261b00000000000000000000000000f5d2fb29fb7d3cfee444a200298f468908cc942\",\n                                                    \"makerAssetAmount\": \"4424020538752105500000\",\n                                                    \"makerFee\": \"0\",\n                                                    \"takerAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                    \"takerAssetData\": \"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2\",\n                                                    \"takerAssetAmount\": \"1000000000000000061\",\n                                                    \"takerFee\": \"0\",\n                                                    \"senderAddress\": \"0x0000000000000000000000000000000000000000\",\n                                                    \"exchangeAddress\": \"0x4f833a24e1f95d70f028921e27040ca56e09ab0b\",\n                                                    \"feeRecipientAddress\": \"0xa258b39954cef5cb142fd567a46cddb31a670124\",\n                                                    \"expirationTimeSeconds\": \"1559422407\",\n                                                    \"salt\": \"1559422141994\",\n                                                    \"signature\": \"0x1cf16c2f3a210965b5e17f51b57b869ba4ddda33df92b0017b4d8da9dacd3152b122a73844eaf50ccde29a42950239ba36a525ed7f1698a8a5e1896cf7d651aed203\"\n                                                },\n                                                \"kind\": \"CANCELLED\",\n                                                \"fillableTakerAssetAmount\": 0,\n                                                \"txHash\": \"0x9e6830a7044b39e107f410e4f765995fd0d3d69d5c3b3582a1701b9d68167560\"\n                                            }\n                                        ]\n                                    }\n                                }\n                            ";
                                                    connection.sendUTF(eventResponse);
                                                    _a.label = 2;
                                                case 2: return [2 /*return*/];
                                            }
                                        });
                                    }); }));
                                }));
                                client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                                expectToBeCalledOnce = true;
                                callback = dev_utils_1.callbackErrorReporter.reportNoErrorCallbackErrors(done, expectToBeCalledOnce)(function (orderEvents) { return __awaiter(_this, void 0, void 0, function () {
                                    return __generator(this, function (_a) {
                                        expect(orderEvents).to.have.length(1);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.makerAssetAmount)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.takerAssetAmount)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.makerFee)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.takerFee)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.salt)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].signedOrder.expirationTimeSeconds)).to.equal(true);
                                        expect(index_1.BigNumber.isBigNumber(orderEvents[0].fillableTakerAssetAmount)).to.equal(true);
                                        client.destroy();
                                        return [2 /*return*/];
                                    });
                                }); });
                                return [4 /*yield*/, client.subscribeToOrdersAsync(callback)];
                            case 2:
                                _a.sent();
                                return [2 /*return*/];
                        }
                    });
                }); })().catch(done);
            });
            return [2 /*return*/];
        });
    }); });
    describe('#onClose', function () {
        it('should trigger when connection is closed', function (done) {
            (function () { return __awaiter(_this, void 0, void 0, function () {
                var wsServer, client;
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                        case 1:
                            wsServer = _a.sent();
                            wsServer.on('connect', function (connection) { return __awaiter(_this, void 0, void 0, function () {
                                return __generator(this, function (_a) {
                                    switch (_a.label) {
                                        case 0: 
                                        // tslint:disable-next-line:custom-no-magic-numbers
                                        return [4 /*yield*/, sleepAsync(100)];
                                        case 1:
                                            // tslint:disable-next-line:custom-no-magic-numbers
                                            _a.sent();
                                            connection.close();
                                            return [2 /*return*/];
                                    }
                                });
                            }); });
                            client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                            client.onClose(function () {
                                client.destroy();
                                done();
                            });
                            return [2 /*return*/];
                    }
                });
            }); })();
        });
    });
    describe('#onReconnected', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('should trigger the callback when reconnected', function (done) {
                (function () { return __awaiter(_this, void 0, void 0, function () {
                    var wsServer, client;
                    var _this = this;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                            case 1:
                                wsServer = _a.sent();
                                wsServer.on('connect', function (connection) { return __awaiter(_this, void 0, void 0, function () {
                                    var requestNum;
                                    var _this = this;
                                    return __generator(this, function (_a) {
                                        requestNum = 0;
                                        connection.on('message', (function (message) { return __awaiter(_this, void 0, void 0, function () {
                                            var jsonRpcRequest, responses, reasonCode, description;
                                            return __generator(this, function (_a) {
                                                switch (_a.label) {
                                                    case 0:
                                                        jsonRpcRequest = JSON.parse(message.utf8Data);
                                                        responses = [
                                                            "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                                }\n                            ",
                                                        ];
                                                        connection.sendUTF(responses[requestNum]);
                                                        requestNum++;
                                                        if (!(requestNum === 1)) return [3 /*break*/, 2];
                                                        // tslint:disable-next-line:custom-no-magic-numbers
                                                        return [4 /*yield*/, sleepAsync(100)];
                                                    case 1:
                                                        // tslint:disable-next-line:custom-no-magic-numbers
                                                        _a.sent();
                                                        reasonCode = WebSocket.connection.CLOSE_REASON_PROTOCOL_ERROR;
                                                        description = WebSocket.connection.CLOSE_DESCRIPTIONS[reasonCode];
                                                        connection.drop(reasonCode, description);
                                                        return [2 /*return*/];
                                                    case 2: return [2 /*return*/];
                                                }
                                            });
                                        }); }));
                                        return [2 /*return*/];
                                    });
                                }); });
                                client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT, { reconnectAfter: 100 });
                                client.onReconnected(function () {
                                    client.destroy();
                                    done();
                                });
                                return [2 /*return*/];
                        }
                    });
                }); })();
            });
            return [2 /*return*/];
        });
    }); });
    describe('#destroy', function () { return __awaiter(_this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            it('should unsubscribe and trigger onClose when close() is called', function (done) {
                (function () { return __awaiter(_this, void 0, void 0, function () {
                    var wsServer, numMessages, client;
                    var _this = this;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, mock_ws_server_1.setupServerAsync()];
                            case 1:
                                wsServer = _a.sent();
                                numMessages = 0;
                                wsServer.on('connect', (function (connection) {
                                    connection.on('message', (function (message) { return __awaiter(_this, void 0, void 0, function () {
                                        var jsonRpcRequest, response;
                                        return __generator(this, function (_a) {
                                            jsonRpcRequest = JSON.parse(message.utf8Data);
                                            if (numMessages === 0) {
                                                expect(jsonRpcRequest.method).to.be.equal('mesh_subscribe');
                                                response = "\n                                {\n                                    \"id\": \"" + jsonRpcRequest.id + "\",\n                                    \"jsonrpc\": \"2.0\",\n                                    \"result\": \"0xab1a3e8af590364c09d0fa6a12103ada\"\n                                }\n                            ";
                                                connection.sendUTF(response);
                                                numMessages++;
                                                return [2 /*return*/];
                                            }
                                            numMessages++;
                                            expect(jsonRpcRequest.method).to.be.equal('mesh_unsubscribe');
                                            return [2 /*return*/];
                                        });
                                    }); }));
                                }));
                                client = new index_1.WSClient("ws://localhost:" + mock_ws_server_1.SERVER_PORT);
                                client.onClose(function () {
                                    expect(numMessages).to.be.equal(2);
                                    done();
                                });
                                // We need to add a sleep here so that we leave time for the client
                                // to get connected before destroying it.
                                // tslint:disable-next-line:custom-no-magic-numbers
                                return [4 /*yield*/, sleepAsync(100)];
                            case 2:
                                // We need to add a sleep here so that we leave time for the client
                                // to get connected before destroying it.
                                // tslint:disable-next-line:custom-no-magic-numbers
                                _a.sent();
                                client.destroy();
                                return [2 /*return*/];
                        }
                    });
                }); })().catch(done);
            });
            return [2 /*return*/];
        });
    }); });
});
function sleepAsync(ms) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            return [2 /*return*/, new Promise(function (resolve) { return setTimeout(resolve, ms); })];
        });
    });
}
