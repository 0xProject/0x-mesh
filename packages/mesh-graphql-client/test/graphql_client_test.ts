import { getContractAddressesForChainOrThrow } from '@0x/contract-addresses';
import { DummyERC20TokenContract } from '@0x/contracts-erc20';
import { ExchangeContract } from '@0x/contracts-exchange';
import {
    blockchainTests,
    constants,
    expect,
    OrderFactory,
    orderHashUtils,
    randomAddress,
    getRandomInteger,
} from '@0x/contracts-test-utils';
import { LimitOrder, LimitOrderFields, Signature } from '@0x/protocol-utils';
import { BlockchainLifecycle, Web3Config, web3Factory } from '@0x/dev-utils';
import { assetDataUtils } from '@0x/order-utils';
import { Web3ProviderEngine } from '@0x/subproviders';
import { DoneCallback, SignedOrder } from '@0x/types';
import { BigNumber, hexUtils } from '@0x/utils';
import { Web3Wrapper } from '@0x/web3-wrapper';
import { gql } from '@apollo/client';
import 'mocha';

import {
    FilterKind,
    MeshGraphQLClient,
    OrderEvent,
    OrderEventEndState,
    OrderWithMetadata,
    RejectedOrderCode,
    SortDirection,
} from '../src/index';

const exchangeProxyV4Address = '0x5315e44798395d4a952530d131249fe00f554565';

function getRandomLimitOrder(fields: Partial<LimitOrderFields> = {}): LimitOrder {
    return new LimitOrder({
        chainId: 1,
        verifyingContract: exchangeProxyV4Address,
        makerToken: '0x34d402f14d58e001d8efbe6585051bf9706aa064',
        takerToken: '0xcdb594a32b1cc3479d8746279712c39d18a07fc0',
        makerAmount: new BigNumber(100),
        takerAmount: new BigNumber(42),
        takerTokenFeeAmount: new BigNumber(0),
        maker: '0x6ecbe1db9ef729cbe972c83fb886247691fb6beb',
        taker: '0x0000000000000000000000000000000000000000',
        sender: '0x0000000000000000000000000000000000000000',
        feeRecipient: '0x0000000000000000000000000000000000000000',
        pool: '0x0000000000000000000000000000000000000000',
        expiry: new BigNumber(Math.floor(Date.now() / 1000 + 60)),
        salt: getRandomInteger('1e5', '1e6'),
        ...fields,
    });
}

import { MeshDeployment, startServerAndClientAsync } from './utils/graphql_server';
import { SignedOrderV4 } from '../src/types';
import { client } from 'websocket';

blockchainTests.resets('GraphQLClient', (env) => {
    describe('integration tests', () => {
        let deployment: MeshDeployment;
        let exchange: ExchangeContract;
        let exchangeAddress: string;
        let makerAddress: string;
        let orderFactory: OrderFactory;
        let provider: Web3ProviderEngine;
        let makerTokenAddress: string;
        let feeTokenAddress: string;

        beforeEach(async () => {
            deployment = await startServerAndClientAsync();
        });

        afterEach(async () => {
            deployment.mesh.stopMesh();
        });

        before(async () => {
            const chainId = await env.getChainIdAsync();
            const accounts = await env.getAccountAddressesAsync();
            [makerAddress] = accounts;

            // Create a new provider so that the ganache instance running on port
            // 8545 will be used instead of the in-process ganache instance.
            const providerConfigs: Web3Config = {
                total_accounts: constants.NUM_TEST_ACCOUNTS,
                shouldUseInProcessGanache: false,
                shouldAllowUnlimitedContractSize: true,
                unlocked_accounts: [makerAddress],
            };
            provider = web3Factory.getRpcProvider(providerConfigs);

            // HACK(jalextowle): We can't currently specify an out of process provider for a blockchainTests
            // suite, so we need to update env.blockchainLifecycle so that the resets suite works as expected.
            env.blockchainLifecycle = new BlockchainLifecycle(new Web3Wrapper(provider));

            exchangeAddress = getContractAddressesForChainOrThrow(chainId).exchange;
            exchange = new ExchangeContract(exchangeAddress, provider);
            const erc20ProxyAddress = getContractAddressesForChainOrThrow(chainId).erc20Proxy;

            // Configure two tokens and an order factory with a maker address so
            // that valid orders can be created easily in the tests. The two dummy tokens are
            // used in the makerToken and feeToken fields.
            const makerToken = new DummyERC20TokenContract('0x34d402f14d58e001d8efbe6585051bf9706aa064', provider);
            const feeToken = new DummyERC20TokenContract('0xcdb594a32b1cc3479d8746279712c39d18a07fc0', provider);
            const mintAmount = new BigNumber('100e18');
            // tslint:disable-next-line: await-promise
            await makerToken.mint(mintAmount).awaitTransactionSuccessAsync({ from: makerAddress });
            // tslint:disable-next-line: await-promise
            await feeToken.mint(mintAmount).awaitTransactionSuccessAsync({ from: makerAddress });
            // tslint:disable-next-line: await-promise
            await makerToken
                .approve(erc20ProxyAddress, new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({ from: makerAddress });
            // tslint:disable-next-line: await-promise
            await feeToken
                .approve(erc20ProxyAddress, new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({ from: makerAddress });
            // tslint:disable-next-line: await-promise
            await makerToken
                .approve('0x5315e44798395d4a952530d131249fe00f554565', new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({ from: makerAddress });
            // tslint:disable-next-line: await-promise
            await feeToken
                .approve('0x5315e44798395d4a952530d131249fe00f554565', new BigNumber('100e18'))
                .awaitTransactionSuccessAsync({ from: makerAddress });

            makerTokenAddress = makerToken.address;
            feeTokenAddress = feeToken.address;
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

        describe('#LimitOrderV4', async () => {
            it('is being signed properly', async () => {
                const testOrder = new LimitOrder({
                    chainId: 1,
                    verifyingContract: '0x5315e44798395d4a952530d131249fe00f554565',
                    makerToken: '0x0b1ba0af832d7c05fd64161e0db78e85978e8082',
                    takerToken: '0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
                    makerAmount: new BigNumber(100),
                    takerAmount: new BigNumber(42),
                    takerTokenFeeAmount: new BigNumber(0),
                    maker: '0x05cac48d17ecc4d8a9db09dde766a03959b98367',
                    taker: '0x0000000000000000000000000000000000000000',
                    sender: '0x0000000000000000000000000000000000000000',
                    feeRecipient: '0x0000000000000000000000000000000000000000',
                    pool: '0x0000000000000000000000000000000000000000000000000000000000000000',
                    expiry: new BigNumber(1611614113),
                    salt: new BigNumber(464085317),
                });
                const expectedHash = '0x4de42b085ed2ce2ef4f8e64df51b7d3ee8586f8ab781ffe45fa45ce884cd278d';
                const hash = testOrder.getHash();
                expect(hash).to.be.eq(expectedHash);
                const expectedSignature: Signature = {
                    signatureType: 3,
                    v: 28,
                    r: '0x575a304e2fe23038a0e4561097b8e47b75270dc28c6bda49bc9a431a6aa9d980',
                    s: '0x43f976de7c59159b0ea2e5594dadc2ec1c7ff4088f76160ab4aa49999996ef7f',
                };
                const signature = testOrder.getSignatureWithKey(
                    '0xee094b79aa0315914955f2f09be9abe541dcdc51f0aae5bec5453e9f73a471a6',
                );

                expect(signature).to.be.deep.eq(expectedSignature);
            });
        });
        describe('#findOrdersV4Async', async () => {
            it('returns all orders when no options are provided', async () => {
                const ordersLength = 1;
                const orders = [];

                console.log('token addresses', makerTokenAddress, feeTokenAddress);
                for (let i = 0; i < ordersLength; i++) {
                    const order = getRandomLimitOrder().clone({
                        maker: makerAddress,
                        makerToken: makerTokenAddress,
                        takerToken: feeTokenAddress,
                        chainId: 1337,
                    });
                    const signature = await order.getSignatureWithProviderAsync(provider);
                    const signedOrder: SignedOrderV4 = {
                        chainId: order.chainId,
                        exchangeAddress: order.verifyingContract,
                        makerToken: order.makerToken,
                        takerToken: order.takerToken,
                        makerAmount: order.makerAmount,
                        takerAmount: order.takerAmount,
                        takerTokenFeeAmount: order.takerTokenFeeAmount,
                        maker: order.maker,
                        taker: order.taker,
                        sender: order.sender,
                        feeRecipient: order.feeRecipient,
                        pool: order.pool,
                        expiry: order.expiry,
                        salt: order.salt,
                        signatureType: signature.signatureType,
                        signatureV: signature.v,
                        signatureR: signature.r,
                        signatureS: signature.s,
                    };
                    orders[i] = signedOrder;
                    console.log(signedOrder);
                }
                console.log('adding orders');
                // const newClient = new MeshGraphQLClient({
                //     httpUrl: `http://localhost:60567/graphql`,
                //     webSocketUrl: `ws://localhost:60567/graphql`,
                // });

                const validationResults = await deployment.client.addOrdersV4Async(orders);
                console.log(validationResults.rejected);
                expect(validationResults.accepted.length).to.be.eq(ordersLength);

                // Verify that all of the orders that were added to the mesh node
                // were returned in the response.
                console.log('checking orders');
                const gotOrders = await deployment.client.findOrdersV4Async();
                console.log(gotOrders);
                // const expectedOrders = orders.map((order) => ({
                //     ...order,
                //     hash: orderHashUtils.getOrderHashHex(order),
                //     fillableTakerAssetAmount: order.takerAssetAmount,
                // }));
                // expectContainsOrders(gotOrders, expectedOrders);
            });
            it('returns orders that match a given query', async () => {
                const ordersLength = 10;
                const orders = [];
                // Create some orders with makerAssetAmount = 1, 2, 3, etc.
                for (let i = 0; i < ordersLength; i++) {
                    orders[i] = await orderFactory.newSignedOrderAsync({
                        makerAssetAmount: new BigNumber(i + 1),
                    });
                }
                const validationResults = await deployment.client.addOrdersAsync(orders);
                expect(validationResults.accepted.length).to.be.eq(ordersLength);

                // Verify that all of the orders that were added to the mesh node
                // were returned in the response.
                const gotOrders = await deployment.client.findOrdersAsync({
                    filters: [{ field: 'makerAssetAmount', kind: FilterKind.LessOrEqual, value: new BigNumber(7) }],
                    sort: [{ field: 'makerAssetAmount', direction: SortDirection.Desc }],
                    limit: 5,
                });
                // We expect 5 orders sorted in descending order by makerAssetAmount starting at 7.
                // I.e. orders with makerAmounts of 7, 6, 5, 4, and 3.
                const expectedOrders = orders.map((order) => ({
                    ...order,
                    hash: orderHashUtils.getOrderHashHex(order),
                    fillableTakerAssetAmount: order.takerAssetAmount,
                }));
                const sortedExpectedOrders = sortOrdersByMakerAssetAmount(expectedOrders).reverse();
                // tslint:disable-next-line: custom-no-magic-numbers
                expect(gotOrders).to.be.deep.eq(sortedExpectedOrders.slice(3, 8));
            });
        });
        describe('#addOrdersAsync', async () => {
            it('accepts valid order', async () => {
                const order = await orderFactory.newSignedOrderAsync({});
                const validationResults = await deployment.client.addOrdersAsync([order]);
                expect(validationResults).to.be.deep.eq({
                    accepted: [
                        {
                            isNew: true,
                            order: {
                                ...order,
                                hash: orderHashUtils.getOrderHashHex(order),
                                fillableTakerAssetAmount: order.takerAssetAmount,
                            },
                        },
                    ],
                    rejected: [],
                });
            });

            it('accepts expired order with "keepExpired"', async () => {
                const order = await orderFactory.newSignedOrderAsync({
                    expirationTimeSeconds: new BigNumber(0),
                });
                const hash = orderHashUtils.getOrderHashHex(order);
                const validationResults = await deployment.client.addOrdersAsync([order], false, { keepExpired: true });
                expect(validationResults).to.be.deep.eq({
                    accepted: [],
                    rejected: [
                        {
                            hash,
                            order,
                            code: RejectedOrderCode.OrderExpired,
                            message: 'order expired according to latest block timestamp',
                        },
                    ],
                });
                const responseOrder = await deployment.client.getOrderAsync(hash);
                expect(responseOrder).to.be.deep.eq({
                    ...order,
                    fillableTakerAssetAmount: new BigNumber(0),
                    hash,
                });
            });

            it('rejects order with invalid signature', async () => {
                const invalidOrder = {
                    ...(await orderFactory.newSignedOrderAsync({})),
                    signature: hexUtils.hash('0x0'),
                };
                const validationResults = await deployment.client.addOrdersAsync([invalidOrder]);
                expect(validationResults).to.be.deep.eq({
                    accepted: [],
                    rejected: [
                        {
                            hash: orderHashUtils.getOrderHashHex(invalidOrder),
                            order: invalidOrder,
                            code: RejectedOrderCode.OrderHasInvalidSignature,
                            message: 'order signature must be valid',
                        },
                    ],
                });
            });
        });

        describe('#getStatsAsync', () => {
            it('Ensure that the stats are correct when no orders have been added', async () => {
                const stats = await deployment.client.getStatsAsync();

                // NOTE(jalextowle): Ensure that the latest block of the returned
                // stats is valid and then clear the field since we don't know
                // the block number of the stats in this test a priori.
                expect(stats.latestBlock).to.not.be.undefined();
                expect(stats.latestBlock.number).to.be.bignumber.greaterThan(0);
                stats.latestBlock = {
                    number: new BigNumber(0),
                    hash: '',
                };
                expect(stats.version).to.not.be.eq('');
                stats.version = '';

                const now = new Date(Date.now());
                const expectedStartOfCurrentUTCDay = `${now.getUTCFullYear()}-${leftPad(
                    now.getUTCMonth() + 1,
                )}-${leftPad(now.getUTCDate())}T00:00:00Z`;
                const expectedStats = {
                    version: '',
                    pubSubTopic: '/0x-orders/version/3/chain/1337/schema/e30=',
                    rendezvous: '/0x-mesh/network/1337/version/2',
                    peerID: deployment.peerID,
                    ethereumChainID: 1337,
                    latestBlock: {
                        number: new BigNumber(0),
                        hash: '',
                    },
                    numPeers: 0,
                    numOrders: 0,
                    numOrdersIncludingRemoved: 0,
                    maxExpirationTime: constants.MAX_UINT256,
                    startOfCurrentUTCDay: new Date(expectedStartOfCurrentUTCDay),
                    ethRPCRequestsSentInCurrentUTCDay: 0,
                    ethRPCRateLimitExpiredRequests: 0,
                };
                expect(stats).to.be.deep.eq(expectedStats);
            });
        });

        describe('#getOrderAsync', async () => {
            it('gets an order by its hash', async () => {
                const order = await orderFactory.newSignedOrderAsync({});
                const validationResults = await deployment.client.addOrdersAsync([order]);
                expect(validationResults.accepted.length).to.be.eq(1);

                const orderHash = orderHashUtils.getOrderHashHex(order);
                const foundOrder = await deployment.client.getOrderAsync(orderHash);
                const expectedOrder = {
                    ...order,
                    hash: orderHash,
                    fillableTakerAssetAmount: order.takerAssetAmount,
                };
                expect(foundOrder).to.be.deep.eq(expectedOrder);
            });
            it('returns null when the order does not exist', async () => {
                const nonExistentOrderHash = '0xabcd46910c6a8a4730878e6e8a4abb328844c0b58f0cdfbb5b6ad28ee0bae347';
                const foundOrder = await deployment.client.getOrderAsync(nonExistentOrderHash);
                expect(foundOrder).to.be.null();
            });
        });

        describe('#findOrdersAsync', async () => {
            it('returns all orders when no options are provided', async () => {
                const ordersLength = 10;
                const orders = [];
                for (let i = 0; i < ordersLength; i++) {
                    orders[i] = await orderFactory.newSignedOrderAsync({});
                }
                const validationResults = await deployment.client.addOrdersAsync(orders);
                expect(validationResults.accepted.length).to.be.eq(ordersLength);

                // Verify that all of the orders that were added to the mesh node
                // were returned in the response.
                const gotOrders = await deployment.client.findOrdersAsync();
                const expectedOrders = orders.map((order) => ({
                    ...order,
                    hash: orderHashUtils.getOrderHashHex(order),
                    fillableTakerAssetAmount: order.takerAssetAmount,
                }));
                expectContainsOrders(gotOrders, expectedOrders);
            });
            it('returns orders that match a given query', async () => {
                const ordersLength = 10;
                const orders = [];
                // Create some orders with makerAssetAmount = 1, 2, 3, etc.
                for (let i = 0; i < ordersLength; i++) {
                    orders[i] = await orderFactory.newSignedOrderAsync({
                        makerAssetAmount: new BigNumber(i + 1),
                    });
                }
                const validationResults = await deployment.client.addOrdersAsync(orders);
                expect(validationResults.accepted.length).to.be.eq(ordersLength);

                // Verify that all of the orders that were added to the mesh node
                // were returned in the response.
                const gotOrders = await deployment.client.findOrdersAsync({
                    filters: [{ field: 'makerAssetAmount', kind: FilterKind.LessOrEqual, value: new BigNumber(7) }],
                    sort: [{ field: 'makerAssetAmount', direction: SortDirection.Desc }],
                    limit: 5,
                });
                // We expect 5 orders sorted in descending order by makerAssetAmount starting at 7.
                // I.e. orders with makerAmounts of 7, 6, 5, 4, and 3.
                const expectedOrders = orders.map((order) => ({
                    ...order,
                    hash: orderHashUtils.getOrderHashHex(order),
                    fillableTakerAssetAmount: order.takerAssetAmount,
                }));
                const sortedExpectedOrders = sortOrdersByMakerAssetAmount(expectedOrders).reverse();
                // tslint:disable-next-line: custom-no-magic-numbers
                expect(gotOrders).to.be.deep.eq(sortedExpectedOrders.slice(3, 8));
            });
        });

        describe('#rawQueryAsync', async () => {
            it('runs a raw query and returns raw results', async () => {
                const response = await deployment.client.rawQueryAsync({
                    query: gql`
                        {
                            stats {
                                numOrders
                            }
                        }
                    `,
                });
                const expectedResponse = {
                    data: {
                        stats: {
                            numOrders: 0,
                        },
                    },
                    loading: false,
                    networkStatus: 7,
                };
                expect(response).to.be.deep.eq(expectedResponse);
            });
        });

        describe('#subscribeToOrdersAsync', async () => {
            it('should receive subscription updates about added orders', (done: DoneCallback) => {
                (async () => {
                    // Keep track of whether or not the test is complete. Used to determine
                    // whether WebSocket errors should be considered test failures.
                    let isDone = false;
                    // Create orders to add to the mesh node.
                    const ordersLength = 10;
                    const orders = [] as SignedOrder[];
                    for (let i = 0; i < ordersLength; i++) {
                        orders[i] = await orderFactory.newSignedOrderAsync({});
                    }

                    // Subscribe to orders and wait for order events.
                    const orderEvents = deployment.client.onOrderEvents();
                    orderEvents.subscribe({
                        error: (err) => {
                            if (isDone && err.message === 'WebSocket connection lost') {
                                // This error is expected to happen after the server is shut down.
                            } else {
                                // Other errors are not expected.
                                throw err;
                            }
                        },
                        next: (events: OrderEvent[]) => {
                            expect(events.length).to.be.eq(orders.length);
                            for (const orderEvent of events) {
                                expect(orderEvent.endState).to.be.eq(OrderEventEndState.Added);
                                const now = new Date().getUTCMilliseconds();
                                // tslint:disable-next-line:custom-no-magic-numbers
                                assertRoughlyEquals(now, orderEvent.timestampMs, secondsToMs(10));
                            }

                            // Ensure that all of the orders that were added had an associated order event emitted.
                            for (const order of orders) {
                                const orderHash = orderHashUtils.getOrderHashHex(order);
                                let hasSeenMatch = false;
                                for (const event of events) {
                                    if (orderHash === event.order.hash) {
                                        hasSeenMatch = true;
                                        const expectedOrder = {
                                            ...order,
                                            hash: orderHash,
                                            fillableTakerAssetAmount: order.takerAssetAmount,
                                        };
                                        expect(event.order).to.be.deep.eq(expectedOrder);
                                        break;
                                    }
                                }
                                expect(hasSeenMatch).to.be.true();
                            }
                            isDone = true;
                            done();
                        },
                    });
                    const validationResults = await deployment.client.addOrdersAsync(orders);
                    expect(validationResults.accepted.length).to.be.eq(ordersLength);
                })().catch(done);
            });

            it('should receive subscription updates about cancelled orders', (done: DoneCallback) => {
                (async () => {
                    // Keep track of whether or not the test is complete. Used to determine
                    // whether WebSocket errors should be considered test failures.
                    let isDone = false;
                    // Add an order and then cancel it.
                    const order = await orderFactory.newSignedOrderAsync({});
                    const validationResults = await deployment.client.addOrdersAsync([order]);
                    expect(validationResults.accepted.length).to.be.eq(1);

                    // Subscribe to order events and assert that only a single cancel event was received.
                    const orderEvents = deployment.client.onOrderEvents();
                    orderEvents.subscribe({
                        error: (err) => {
                            if (isDone && err.message === 'WebSocket connection lost') {
                                // This error is expected to happen after the server is shut down.
                            } else {
                                // Other errors are not expected.
                                throw err;
                            }
                        },
                        next: (events: OrderEvent[]) => {
                            // Ensure that the correct cancel event was logged.
                            expect(events.length).to.be.eq(1);
                            const [orderEvent] = events;
                            expect(orderEvent.endState).to.be.eq(OrderEventEndState.Cancelled);
                            const expectedOrder = {
                                ...order,
                                hash: orderHashUtils.getOrderHashHex(order),
                                fillableTakerAssetAmount: constants.ZERO_AMOUNT,
                            };
                            expect(orderEvent.order).to.be.deep.eq(expectedOrder);
                            const now = new Date().getUTCMilliseconds();
                            assertRoughlyEquals(orderEvent.timestampMs, now, secondsToMs(2));
                            expect(orderEvent.contractEvents.length).to.be.eq(1);

                            // Ensure that the contract event is correct.
                            const [contractEvent] = orderEvent.contractEvents;
                            expect(contractEvent.address).to.be.eq(exchangeAddress);
                            expect(contractEvent.kind).to.be.equal('ExchangeCancelEvent');
                            expect(contractEvent.logIndex).to.be.eq(0);
                            expect(contractEvent.isRemoved).to.be.false();
                            expect(contractEvent.txIndex).to.be.eq(0);
                            const hashLength = 66;
                            expect(contractEvent.blockHash.length).to.be.eq(hashLength);
                            expect(contractEvent.blockHash).to.not.be.eq(constants.NULL_BYTES32);
                            expect(contractEvent.txHash.length).to.be.eq(hashLength);
                            const parameters = contractEvent.parameters;
                            parameters.makerAddress = parameters.makerAddress.toLowerCase();
                            parameters.senderAddress = parameters.makerAddress;
                            expect(parameters.feeRecipientAddress.toLowerCase()).to.be.eq(order.feeRecipientAddress);
                            expect(parameters.makerAddress.toLowerCase()).to.be.eq(makerAddress);
                            expect(parameters.makerAssetData).to.be.eq(order.makerAssetData);
                            expect(parameters.orderHash).to.be.eq(orderHashUtils.getOrderHashHex(order));
                            expect(parameters.senderAddress.toLowerCase()).to.be.eq(makerAddress);
                            expect(parameters.takerAssetData).to.be.eq(order.takerAssetData);
                            isDone = true;
                            done();
                        },
                    });

                    // Cancel an order and then wait for the emitted order event.
                    // tslint:disable-next-line: await-promise
                    await exchange.cancelOrder(order).awaitTransactionSuccessAsync({ from: makerAddress });
                })().catch(done);
            });
        });
    });
});

function assertRoughlyEquals(a: number, b: number, delta: number): void {
    expect(Math.abs(a - b)).to.be.lessThan(delta);
}

function leftPad(a: number, paddingDigits: number = 2): string {
    return `${'0'.repeat(paddingDigits - a.toString().length)}${a.toString()}`;
}

function secondsToMs(seconds: number): number {
    const msPerSecond = 1000;
    return seconds * msPerSecond;
}

function sortOrdersByMakerAssetAmount(orders: OrderWithMetadata[]): OrderWithMetadata[] {
    return orders.sort((a, b) => {
        if (a.makerAssetAmount.gt(b.makerAssetAmount)) {
            return 1;
        } else if (a.makerAssetAmount.lt(b.makerAssetAmount)) {
            return -1;
        }
        return 0;
    });
}

// Verify that all of the orders that were added to the mesh node
// were returned in the `getOrders` rpc response
function expectContainsOrders(gotOrders: OrderWithMetadata[], expectedOrders: OrderWithMetadata[]): void {
    for (const expectedOrder of expectedOrders) {
        let hasSeenMatch = false;
        for (const gotOrder of gotOrders) {
            if (expectedOrder.hash === gotOrder.hash) {
                hasSeenMatch = true;
                expect(gotOrder).to.be.deep.eq(expectedOrder);
                break;
            }
        }
        expect(hasSeenMatch).to.be.true();
    }
}
// tslint:disable-line:max-file-line-count
