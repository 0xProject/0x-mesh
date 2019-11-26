import { callbackErrorReporter } from '@0x/dev-utils';
import { DoneCallback } from '@0x/types';
import * as chai from 'chai';
import 'mocha';
import * as WebSocket from 'websocket';

import { BigNumber, OrderEvent, WSClient } from '../src/index';
import { WSMessage } from '../src/types';

import { chaiSetup } from './utils/chai_setup';
import { SERVER_PORT, setupServerAsync, stopServer } from './utils/mock_ws_server';

chaiSetup.configure();
const expect = chai.expect;

describe('WSClient', () => {
    afterEach(() => {
        stopServer();
    });
    describe('#addOrdersAsync', async () => {
        it('sends a mesh_addOrders request and converts all numerical fields to BigNumbers in returned signedOrders', async () => {
            const wsServer = await setupServerAsync();
            wsServer.on('connect', ((connection: WebSocket.connection) => {
                let requestNum = 0;
                connection.on('message', ((message: WSMessage) => {
                    const jsonRpcRequest = JSON.parse(message.utf8Data);
                    const snapshotID = '123';
                    const snapshotTimestamp = '2009-11-10T23:00:00Z';
                    const responses = [
                        // Heartbeat subscription (under-the-hood)
                        `
                            {
                                "id": "${jsonRpcRequest.id}",
                                "jsonrpc": "2.0",
                                "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                            }
                        `,
                        // mesh_addOrders response
                        `
                            {
                                "id": "${jsonRpcRequest.id}",
                                "jsonrpc": "2.0",
                                "result": {
                                    "accepted": [
                                        {
                                            "orderHash": "0xa0fcb54919f0b3823aa14b3f511146f6ac087ab333a70f9b24bbb1ba657a4250",
                                            "signedOrder": {
                                                "makerAddress": "0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149",
                                                "makerAssetData": "0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
                                                "makerAssetAmount": "1000000000000000000",
                                                "makerFee": "0",
                                                "takerAddress": "0x0000000000000000000000000000000000000000",
                                                "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                                                "takerAssetAmount": "10000000000000000000000",
                                                "takerFee": "0",
                                                "senderAddress": "0x0000000000000000000000000000000000000000",
                                                "exchangeAddress": "0x080bf510FCbF18b91105470639e9561022937712",
                                                "feeRecipientAddress": "0x0000000000000000000000000000000000000000",
                                                "expirationTimeSeconds": "1586340602",
                                                "salt": "41253767178111694375645046549067933145709740457131351457334397888365956743955",
                                                "signature": "0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02"
                                            },
                                            "fillableTakerAssetAmount": "10000000000000000000000"
                                        }
                                    ],
                                    "rejected": []
                                }
                            }
                        `,
                    ];
                    connection.sendUTF(responses[requestNum]);
                    requestNum++;
                }) as any);
            }) as any);

            const signedOrders = [
                {
                    makerAddress: '0xa3ece5d5b6319fa785efc10d3112769a46c6e149',
                    takerAddress: '0x0000000000000000000000000000000000000000',
                    makerAssetAmount: new BigNumber('1000000000000000000'),
                    takerAssetAmount: new BigNumber('10000000000000000000000'),
                    expirationTimeSeconds: new BigNumber('1586340602'),
                    makerFee: new BigNumber('0'),
                    takerFee: new BigNumber('0'),
                    feeRecipientAddress: '0x0000000000000000000000000000000000000000',
                    senderAddress: '0x0000000000000000000000000000000000000000',
                    salt: new BigNumber(
                        '41253767178111694375645046549067933145709740457131351457334397888365956743955',
                    ),
                    makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
                    takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
                    exchangeAddress: '0x080bf510fcbf18b91105470639e9561022937712',
                    signature:
                        '0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02',
                },
            ];
            const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
            const validationResults = await client.addOrdersAsync(signedOrders);
            expect(typeof validationResults === 'object').to.equal(true);
            expect(validationResults.accepted).to.have.length(1);
            expect(validationResults.rejected).to.have.length(0);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.makerAssetAmount)).to.equal(true);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.takerAssetAmount)).to.equal(true);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.makerFee)).to.equal(true);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.takerFee)).to.equal(true);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.salt)).to.equal(true);
            expect(BigNumber.isBigNumber(validationResults.accepted[0].signedOrder.expirationTimeSeconds)).to.equal(
                true,
            );
            expect(BigNumber.isBigNumber(validationResults.accepted[0].fillableTakerAssetAmount)).to.equal(true);

            client.destroy();
        });
    });
    describe('#getStatsAsync', async () => {
        /* FIXME */
    });
    describe('#getOrdersAsync', async () => {
        it('properly makes multiple paginated requests under-the-hood and returns all signedOrders', async () => {
            const wsServer = await setupServerAsync();
            wsServer.on('connect', ((connection: WebSocket.connection) => {
                let requestNum = 0;
                connection.on('message', ((message: WSMessage) => {
                    const jsonRpcRequest = JSON.parse(message.utf8Data);
                    const responses = [
                        // Heartbeat subscription (under-the-hood)
                        `
                            {
                                "id": "${jsonRpcRequest.id}",
                                "jsonrpc": "2.0",
                                "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                            }
                        `,
                        // First paginated request
                        `
                            {
                                "id": "${jsonRpcRequest.id}",
                                "jsonrpc": "2.0",
                                "result": {
                                    "snapshotID": "123",
                                    "ordersInfos": [
                                        {
                                            "orderHash": "0xa0fcb54919f0b3823aa14b3f511146f6ac087ab333a70f9b24bbb1ba657a4250",
                                            "signedOrder": {
                                                "makerAddress": "0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149",
                                                "makerAssetData": "0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
                                                "makerAssetAmount": "1000000000000000000",
                                                "makerFee": "0",
                                                "takerAddress": "0x0000000000000000000000000000000000000000",
                                                "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                                                "takerAssetAmount": "10000000000000000000000",
                                                "takerFee": "0",
                                                "senderAddress": "0x0000000000000000000000000000000000000000",
                                                "exchangeAddress": "0x080bf510FCbF18b91105470639e9561022937712",
                                                "feeRecipientAddress": "0x0000000000000000000000000000000000000000",
                                                "expirationTimeSeconds": "1586340602",
                                                "salt": "41253767178111694375645046549067933145709740457131351457334397888365956743955",
                                                "signature": "0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02"
                                            },
                                            "fillableTakerAssetAmount": "10000000000000000000000"
                                        }
                                    ]
                                }
                            }
                            `,
                        // Second paginated request
                        `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": {
                                        "snapshotID": "123",
                                        "ordersInfos": []
                                    }
                                }
                            `,
                    ];
                    connection.sendUTF(responses[requestNum]);
                    requestNum++;
                }) as any);
            }) as any);

            const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
            const perPage = 1;
            const { ordersInfos } = await client.getOrdersAsync(perPage);

            expect(ordersInfos).to.have.length(1);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.makerAssetAmount)).to.equal(true);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.takerAssetAmount)).to.equal(true);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.makerFee)).to.equal(true);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.takerFee)).to.equal(true);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.salt)).to.equal(true);
            expect(BigNumber.isBigNumber(ordersInfos[0].signedOrder.expirationTimeSeconds)).to.equal(true);

            client.destroy();
        });
    });
    describe('#subscribeToOrdersAsync', async () => {
        it('should receive subscription updates', (done: DoneCallback) => {
            (async () => {
                const wsServer = await setupServerAsync();
                wsServer.on('connect', ((connection: WebSocket.connection) => {
                    let requestNum = 0;
                    connection.on('message', (async (message: WSMessage) => {
                        const jsonRpcRequest = JSON.parse(message.utf8Data);
                        const responses = [
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                                }
                            `,
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xc2ba3e8af590364c09d0fa6a12103adb"
                                }
                            `,
                            // Response to unsubscribe
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc":"2.0",
                                    "result":true,
                                }
                            `,
                        ];
                        connection.sendUTF(responses[requestNum]);
                        requestNum++;

                        if (requestNum === 2) {
                            // tslint:disable-next-line:custom-no-magic-numbers
                            await sleepAsync(100);

                            const eventResponse = `
                                {
                                    "jsonrpc":"2.0",
                                    "method":"mesh_subscription",
                                    "params": {
                                        "subscription":"0xc2ba3e8af590364c09d0fa6a12103adb",
                                        "result": [
                                            {
                                                "orderHash": "0x96e6eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ecc13fd4",
                                                "signedOrder": {
                                                    "makerAddress": "0x50f84bbee6fb250d6f49e854fa280445369d64d9",
                                                    "makerAssetData": "0xf47261b00000000000000000000000000f5d2fb29fb7d3cfee444a200298f468908cc942",
                                                    "makerAssetAmount": "4424020538752105500000",
                                                    "makerFee": "0",
                                                    "takerAddress": "0x0000000000000000000000000000000000000000",
                                                    "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                                                    "takerAssetAmount": "1000000000000000061",
                                                    "takerFee": "0",
                                                    "senderAddress": "0x0000000000000000000000000000000000000000",
                                                    "exchangeAddress": "0x4f833a24e1f95d70f028921e27040ca56e09ab0b",
                                                    "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
                                                    "expirationTimeSeconds": "1559422407",
                                                    "salt": "1559422141994",
                                                    "signature": "0x1cf16c2f3a210965b5e17f51b57b869ba4ddda33df92b0017b4d8da9dacd3152b122a73844eaf50ccde29a42950239ba36a525ed7f1698a8a5e1896cf7d651aed203"
                                                },
                                                "endState": "CANCELLED",
                                                "fillableTakerAssetAmount": 0,
                                                "contractEvents": [
                                                    {
                                                        "blockHash": "0x1be2eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ec11a4d2",
                                                        "txHash": "0xbcce172374dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ec232e3a",
                                                        "txIndex": 23,
                                                        "logIndex": 0,
                                                        "isRemoved": false,
                                                        "address": "0x4f833a24e1f95d70f028921e27040ca56e09ab0b",
                                                        "kind": "ExchangeCancelEvent",
                                                        "parameters": {
                                                            "makerAddress": "0x50f84bbee6fb250d6f49e854fa280445369d64d9",
                                                            "senderAddress": "0x0000000000000000000000000000000000000000",
                                                            "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
                                                            "orderHash": "0x96e6eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ecc13fd4",
                                                            "makerAssetData": "0xf47261b00000000000000000000000000f5d2fb29fb7d3cfee444a200298f468908cc942",
                                                            "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
                                                        }
                                                    }
                                                ]
                                            }
                                        ]
                                    }
                                }
                            `;
                            connection.sendUTF(eventResponse);
                        }
                    }) as any);
                }) as any);

                const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
                const expectToBeCalledOnce = true;
                const callback = callbackErrorReporter.reportNoErrorCallbackErrors(done, expectToBeCalledOnce)(
                    async (orderEvents: OrderEvent[]) => {
                        expect(orderEvents).to.have.length(1);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.makerAssetAmount)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.takerAssetAmount)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.makerFee)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.takerFee)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.salt)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].signedOrder.expirationTimeSeconds)).to.equal(true);
                        expect(BigNumber.isBigNumber(orderEvents[0].fillableTakerAssetAmount)).to.equal(true);

                        client.destroy();
                    },
                );
                await client.subscribeToOrdersAsync(callback);
            })().catch(done);
        });
    });
    describe('#_subscribeToHeartbeatAsync', async () => {
        it('should receive subscription updates', (done: DoneCallback) => {
            (async () => {
                const wsServer = await setupServerAsync();
                wsServer.on('connect', ((connection: WebSocket.connection) => {
                    let requestNum = 0;
                    connection.on('message', (async (message: WSMessage) => {
                        const jsonRpcRequest = JSON.parse(message.utf8Data);
                        const responses = [
                            // Heartbeat subscription (under-the-hood)
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                                }
                            `,
                            // Requested heartbeat subscription
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                                }
                            `,
                            // Response to unsubscribe
                            `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc":"2.0",
                                    "result":true,
                                }
                            `,
                        ];
                        connection.sendUTF(responses[requestNum]);
                        requestNum++;

                        if (requestNum === 2) {
                            // tslint:disable-next-line:custom-no-magic-numbers
                            await sleepAsync(100);

                            const eventResponse = `
                                {
                                    "jsonrpc":"2.0",
                                    "method":"mesh_subscription",
                                    "params": {
                                        "subscription":"0xab1a3e8af590364c09d0fa6a12103ada",
                                        "result":"tick"
                                    }
                                }
                            `;
                            connection.sendUTF(eventResponse);
                        }
                    }) as any);
                }) as any);

                const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
                const expectToBeCalledOnce = true;
                const callback = callbackErrorReporter.reportNoErrorCallbackErrors(done, expectToBeCalledOnce)(
                    async (ack: string) => {
                        expect(ack).to.be.equal('tick');
                        client.destroy();
                    },
                );
                await (client as any)._subscribeToHeartbeatAsync(callback);
            })().catch(done);
        });
    });
    describe('#onClose', () => {
        it('should trigger when connection is closed', (done: DoneCallback) => {
            // tslint:disable-next-line:no-floating-promises
            (async () => {
                const wsServer = await setupServerAsync();
                wsServer.on('connect', async (connection: WebSocket.connection) => {
                    // tslint:disable-next-line:custom-no-magic-numbers
                    await sleepAsync(100);
                    connection.close();
                });

                const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
                client.onClose(() => {
                    client.destroy();
                    done();
                });
            })();
        });
    });
    describe('#onReconnected', async () => {
        it('should trigger the callback when reconnected', (done: DoneCallback) => {
            // tslint:disable-next-line:no-floating-promises
            (async () => {
                const wsServer = await setupServerAsync();
                let connectionNum = 0;
                wsServer.on('connect', async (connection: WebSocket.connection) => {
                    let requestNum = 0;
                    connectionNum++;
                    connection.on('message', (async (message: WSMessage) => {
                        const jsonRpcRequest = JSON.parse(message.utf8Data);
                        if (requestNum === 0) {
                            const response = `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                                }
                            `;
                            connection.sendUTF(response);
                            if (connectionNum === 1) {
                                // tslint:disable-next-line:custom-no-magic-numbers
                                await sleepAsync(100);
                                const reasonCode = WebSocket.connection.CLOSE_REASON_PROTOCOL_ERROR;
                                const description = (WebSocket.connection as any).CLOSE_DESCRIPTIONS[reasonCode];
                                connection.drop(reasonCode, description);
                            }
                        }
                        requestNum++;
                    }) as any);
                });

                const client = new WSClient(`ws://localhost:${SERVER_PORT}`, { reconnectAfter: 100 });
                client.onReconnected(async () => {
                    // We need to add a sleep here so that we leave time for the client
                    // to get connected before destroying it.
                    // tslint:disable-next-line:custom-no-magic-numbers
                    await sleepAsync(100);
                    client.destroy();
                    done();
                });
            })();
        });
    });
    describe('#destroy', async () => {
        it('should unsubscribe and trigger onClose when close() is called', (done: DoneCallback) => {
            // tslint:disable-next-line:no-floating-promises
            (async () => {
                const wsServer = await setupServerAsync();
                let numMessages = 0;
                wsServer.on('connect', ((connection: WebSocket.connection) => {
                    connection.on('message', (async (message: WSMessage) => {
                        const jsonRpcRequest = JSON.parse(message.utf8Data);
                        if (numMessages === 0) {
                            expect(jsonRpcRequest.method).to.be.equal('mesh_subscribe');
                            const response = `
                                {
                                    "id": "${jsonRpcRequest.id}",
                                    "jsonrpc": "2.0",
                                    "result": "0xab1a3e8af590364c09d0fa6a12103ada"
                                }
                            `;
                            connection.sendUTF(response);
                            numMessages++;
                            return;
                        }
                        numMessages++;
                    }) as any);
                }) as any);

                const client = new WSClient(`ws://localhost:${SERVER_PORT}`);
                client.onClose(() => {
                    expect(numMessages).to.be.equal(2);
                    done();
                });
                // We need to add a sleep here so that we leave time for the client
                // to get connected before destroying it.
                // tslint:disable-next-line:custom-no-magic-numbers
                await sleepAsync(100);
                client.destroy();
            })().catch(done);
        });
    });
});

async function sleepAsync(ms: number): Promise<NodeJS.Timer> {
    return new Promise<NodeJS.Timer>(resolve => setTimeout(resolve, ms));
}
