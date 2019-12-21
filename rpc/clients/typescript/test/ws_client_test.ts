import {getContractAddressesForChainOrThrow} from '@0x/contract-addresses';
import {artifacts, DummyERC20TokenContract} from '@0x/contracts-erc20';
import {blockchainTests, constants, expect, OrderFactory, orderHashUtils} from '@0x/contracts-test-utils';
import {callbackErrorReporter, Web3Config, web3Factory} from '@0x/dev-utils';
import {assetDataUtils} from '@0x/order-utils';
import {Web3ProviderEngine} from '@0x/subproviders';
import {DoneCallback, SignedOrder} from '@0x/types';
import {BigNumber} from '@0x/utils';
import 'mocha';
import * as WebSocket from 'websocket';

import {OrderEvent, OrderEventEndState, WSClient} from '../src/index';
import {WSMessage} from '../src/types';

import {SERVER_PORT, setupServerAsync, stopServer} from './utils/mock_ws_server';
import {MeshDeployment, startServerAndClientAsync} from './utils/ws_server';

blockchainTests.resets('WSClient', env => {
    describe('integration tests', () => {
        let orderFactory: OrderFactory;
        let provider: Web3ProviderEngine;
        let deployment: MeshDeployment;

        const oneSecondInMs = 1000;

        beforeEach(async () => {
            deployment = await startServerAndClientAsync();
        });

        afterEach(async () => {
            deployment.client.destroy();
            deployment.mesh.stopMesh();
        });

        async function deployErc20TokenAsync(name: string, symbol: string): Promise<DummyERC20TokenContract> {
            return DummyERC20TokenContract.deployFrom0xArtifactAsync(
                artifacts.DummyERC20Token,
                provider,
                env.txDefaults,
                artifacts,
                name,
                symbol,
                new BigNumber(18),
                new BigNumber('100e18'),
            );
        }

        before(async () => {
            const chainId = await env.getChainIdAsync();
            const accounts = await env.getAccountAddressesAsync();
            const [makerAddress] = accounts;

            // Create a new provider so that the ganache instance running on port
            // 8545 will be used instead of the in-process ganache instance.
            const providerConfigs: Web3Config = {
                total_accounts: constants.NUM_TEST_ACCOUNTS,
                shouldUseInProcessGanache: false,
                shouldAllowUnlimitedContractSize: true,
                unlocked_accounts: [makerAddress],
            };
            provider = web3Factory.getRpcProvider(providerConfigs);

            const exchangeAddress = getContractAddressesForChainOrThrow(chainId).exchange;
            const erc20ProxyAddress = getContractAddressesForChainOrThrow(chainId).erc20Proxy;

            // Configure two tokens and an order factory with a maker address so
            // that valid orders can be created easily in the tests.
            const makerToken = await deployErc20TokenAsync('MakerToken', 'MKT');
            const feeToken = await deployErc20TokenAsync('FeeToken', 'FEE');
            await makerToken
                .approve(erc20ProxyAddress, new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({from: makerAddress});
            await feeToken
                .approve(erc20ProxyAddress, new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({from: makerAddress});
            orderFactory = new OrderFactory(constants.TESTRPC_PRIVATE_KEYS[accounts.indexOf(makerAddress)], {
                ...constants.STATIC_ORDER_PARAMS,
                feeRecipientAddress: constants.NULL_ADDRESS,
                makerAddress,
                exchangeAddress,
                chainId: 1337,
                makerAssetData: assetDataUtils.encodeERC20AssetData(makerToken.address),
                takerAssetData: assetDataUtils.encodeERC20AssetData(makerToken.address),
                makerFeeAssetData: assetDataUtils.encodeERC20AssetData(feeToken.address),
                takerFeeAssetData: assetDataUtils.encodeERC20AssetData(feeToken.address),
            });
        });

        describe('#addOrdersAsync', async () => {
            it('correctly validates a v3 order signed by OrderFactory `addOrdersAsync`', async () => {
                const order = await orderFactory.newSignedOrderAsync({});
                const validationResults = await deployment.client.addOrdersAsync([order]);
                expect(validationResults).to.be.deep.eq({
                    accepted: [{
                        fillableTakerAssetAmount: order.takerAssetAmount,
                        isNew: true,
                        orderHash: orderHashUtils.getOrderHashHex(order),
                        signedOrder: order,
                    }],
                    rejected: [],
                });
            });
        });

        describe('#getStats', () => {
            it('Ensure that the stats are correct when no orders have been added', async () => {
                const stats = await deployment.client.getStatsAsync();

                // NOTE(jalextowle): Ensure that the latest block of the returned
                // stats is valid and then clear the field since we don't know
                // the block number of the stats in this test a priori.
                expect(stats.latestBlock).to.not.be.undefined();
                expect(stats.latestBlock.number).to.be.greaterThan(0);
                stats.version = '';
                stats.latestBlock = {
                    number: 0,
                    hash: '',
                };

                const now = new Date(Date.now());
                const expectedStartOfCurrentUTCDay = `${now.getUTCFullYear()}-${now.getUTCMonth() +
                    1}-${now.getUTCDate()}T00:00:00Z`;
                const expectedStats = {
                    version: '',
                    pubSubTopic: '/0x-orders/network/1337/version/2',
                    rendezvous: '/0x-mesh/network/1337/version/2',
                    peerID: deployment.peerID,
                    ethereumChainID: 1337,
                    latestBlock: {
                        number: 0,
                        hash: '',
                    },
                    numPeers: 0,
                    numOrders: 0,
                    numOrdersIncludingRemoved: 0,
                    numPinnedOrders: 0,
                    maxExpirationTime: constants.MAX_UINT256.toString(),
                    startOfCurrentUTCDay: expectedStartOfCurrentUTCDay,
                    ethRPCRequestsSentInCurrentUTCDay: 0,
                    ethRPCRateLimitExpiredRequests: 0,
                };
                expect(stats).to.be.deep.eq(expectedStats);
            });
        });

        // This pattern will only match strings with the following pattern:
        // [ 8 digits or letters ]-[ 4 digits or letters ]-[ 4 digits or letters ]-[ 4 digits or letters ]-[ 12 digits or letters ]
        const ganacheSnapshotIdPattern = /^([0-9]|[a-z]){8}-([0-9]|[a-z]){4}-([0-9]|[a-z]){4}-([0-9]|[a-z]){4}-([0-9]|[a-z]){12}$/;

        describe('#getOrdersAsync', async () => {
            it('properly makes multiple paginated requests under-the-hood and returns all signedOrders', async () => {
                const ordersLength = 10;
                const orders = [];
                for (let i = 0; i < ordersLength; i++) {
                    orders[i] = await orderFactory.newSignedOrderAsync({});
                }
                const validationResults = await deployment.client.addOrdersAsync(orders);
                expect(validationResults.accepted.length).to.be.eq(ordersLength);

                // NOTE(jalextowle): The time returned by Date uses milliseconds, but
                // the mesh timestamp only uses second. Multiplying the seconds timestamp
                // by 1000 gives us a comparable value. We only try to ensure that this
                // timestamp is approximately equal (within 1 second) because the server
                // will receive the request slightly after it is sent.
                const now = new Date(Date.now()).getTime();
                const perPage = 10;
                const response = await deployment.client.getOrdersAsync(perPage);
                assertRoughlyEquals(now, response.snapshotTimestamp * oneSecondInMs, oneSecondInMs);

                // Verify that all of the orders that were added to the mesh node
                // were returned in the `getOrders` rpc response, and that the
                // ganache snapshot ID in the response meets the expected schema.
                expect(ganacheSnapshotIdPattern.test(response.snapshotID));
                for (const order of orders) {
                    let hasSeenMatch = false;
                    for (const responseOrder of response.ordersInfos) {
                        if (orderHashUtils.getOrderHashHex(order) === responseOrder.orderHash) {
                            hasSeenMatch = true;
                            expect(order).to.be.deep.eq(responseOrder.signedOrder);
                            break;
                        }
                    }
                    expect(hasSeenMatch).to.be.true();
                }
            });
        });

        describe('#_subscribeToHeartbeatAsync', async () => {
            it('should receive subscription updates', (done: DoneCallback) => {
                (async () => {
                    const expectToBeCalledOnce = true;
                    const callback = callbackErrorReporter.reportNoErrorCallbackErrors(done, expectToBeCalledOnce)(
                        async (ack: string) => {
                            expect(ack).to.be.equal('tick');
                        },
                    );
                    await (deployment.client as any)._subscribeToHeartbeatAsync(callback);
                })().catch(done);
            });
        });

        describe('#subscribeToOrdersAsync', async () => {
            it('should receive subscription updates', (done: DoneCallback) => {
                (async () => {
                    // Create the subscription with a testing callback.
                    const orders = [] as SignedOrder[];
                    const now = new Date(Date.now()).getTime();
                    const callback = (orderEvents: OrderEvent[]) => {
                        expect(orderEvents.length).to.be.eq(orders.length);
                        for (const orderEvent of orderEvents) {
                            expect(orderEvent.endState).to.be.eq(OrderEventEndState.Added);
                            assertRoughlyEquals(now, orderEvent.timestampMs, oneSecondInMs);
                        }

                        // Ensure that all of the orders that were added had an associated order event emitted.
                        for (const order of orders) {
                            const orderHash = orderHashUtils.getOrderHashHex(order);
                            let hasSeenMatch = false;
                            for (const orderEvent of orderEvents) {
                                if (orderHash === orderEvent.orderHash) {
                                    hasSeenMatch = true;
                                    expect(orderEvent.signedOrder).to.be.deep.eq(order);
                                    expect(orderEvent.fillableTakerAssetAmount).to.be.bignumber.eq(order.takerAssetAmount);
                                    break;
                                }
                            }
                            expect(hasSeenMatch).to.be.true();
                        }

                        done();
                    };
                    const subscription = deployment.client.subscribeToOrdersAsync(callback);

                    // Add orders to the mesh node.
                    const ordersLength = 10;
                    for (let i = 0; i < ordersLength; i++) {
                        orders[i] = await orderFactory.newSignedOrderAsync({});
                    }
                    const validationResults = await deployment.client.addOrdersAsync(orders);
                    expect(validationResults.accepted.length).to.be.eq(ordersLength);

                    await subscription;
                })().catch(done);
            });
        });
    });

    describe('unit tests', () => {
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
                        stopServer();
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

                    const client = new WSClient(`ws://localhost:${SERVER_PORT}`, {reconnectDelay: 100});
                    client.onReconnected(async () => {
                        // We need to add a sleep here so that we leave time for the client
                        // to get connected before destroying it.
                        // tslint:disable-next-line:custom-no-magic-numbers
                        await sleepAsync(100);
                        client.destroy();
                        stopServer();
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
});

function assertRoughlyEquals(a: number, b: number, delta: number): void {
    expect(Math.abs(a - b)).to.be.lessThan(delta);
}

async function sleepAsync(ms: number): Promise<NodeJS.Timer> {
    return new Promise<NodeJS.Timer>(resolve => setTimeout(resolve, ms));
}
// tslint:disable-line:max-file-line-count
