import { RPCSubprovider, Web3ProviderEngine } from '@0x/subproviders';
import * as chai from 'chai';
import { ZeroExProvider } from 'ethereum-types';
import 'mocha';

import { WrapperConfig } from '../src/types';
import { configToWrapperConfig } from '../src/wrapper_conversion';

import { chaiSetup } from './utils/chai_setup';
chaiSetup.configure();
const expect = chai.expect;

// FIXME(jalextowle): Delete these linting rules once the tests are written.
// tslint:disable:no-empty
// NOTE(jalextowle): Throughtout these unit tests, some fields will be explicitly
// marked as 'undefined' while others will be implicitley marked as 'undefined' by
// being left out entirely. This issue https://github.com/chaijs/chai/issues/735
// explains why these need to be here. In short, we declare only some of the fields
// of some of the objects in the conversion functions, and these are precisely the
// values that must be explicitley marked as 'undefined' in the tests.
describe('Wrapper Conversion Unit Tests', () => {
    describe('#configToWrapperConfig', () => {
        it('minimal config', () => {
            expect(
                configToWrapperConfig({
                    ethereumChainID: 1337,
                }),
            ).to.be.deep.equal({
                bootstrapList: undefined,
                customContractAddresses: undefined,
                customOrderFilter: undefined,
                ethereumChainID: 1337,
                web3Provider: undefined,
            });
        });
        it('full config', () => {
            const ethereumRPCURL = 'http://localhost:8545';
            const provider = new Web3ProviderEngine();
            provider.addProvider(new RPCSubprovider(ethereumRPCURL));
            expect((provider as any).isStandardizedProvider).to.be.undefined('');
            const config = configToWrapperConfig({
                verbosity: 5,
                ethereumRPCURL,
                useBootstrapList: true,
                bootstrapList: [
                    '/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF',
                    '/ip4/3.214.190.67/tcp/60559/ipfs/16Uiu2HAmGx8Z6gdq6T5AQE54GMtqDhDFhizywTy1o28NJbAMMumG',
                    '/ip4/3.214.190.67/tcp/60560/ipfs/16Uiu2HAmGx8Z6gdq7T5AQE54GMtqDhDFhizywTy1o28NJbAMMumH',
                ],
                blockPollingIntervalSeconds: 1,
                ethereumRPCMaxContentLength: 1024,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 30000,
                customContractAddresses: {
                    exchange: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
                    devUtils: '0x38ef19fdf8e8415f18c307ed71967e19aac28ba1',
                    erc20Proxy: '0x1dc4c1cefef38a777b15aa20260a54e584b16c48',
                    erc721Proxy: '0x1d7022f5b17d2f8b695918fb48fa1089c9f85401',
                    erc1155Proxy: '0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f',
                },
                maxOrdersInStorage: 30000,
                customOrderFilter: {
                    properties: {
                        makerAssetData: {
                            const: '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
                        },
                    },
                },
                ethereumChainID: 1337,
                web3Provider: provider,
            });
            // NOTE(jalextowle): It's hard to verify that the provider is correct because it is a very large object,
            // so we use this as a rough proxy. This check in combination with the check that `isStandardizedProvider`
            // was initially undefined provides a basic check to make sure that `standerizeOrThrow` was called on the
            // provider.
            expect((config.web3Provider! as any).isStandardizedProvider).to.be.true(''); // tslint:disable-line:no-non-null-assertion
            expect({
                ...config,
                web3Provider: undefined,
            }).to.be.deep.equal({
                verbosity: 5,
                ethereumRPCURL,
                useBootstrapList: true,
                bootstrapList:
                    '/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF,' +
                    '/ip4/3.214.190.67/tcp/60559/ipfs/16Uiu2HAmGx8Z6gdq6T5AQE54GMtqDhDFhizywTy1o28NJbAMMumG,' +
                    '/ip4/3.214.190.67/tcp/60560/ipfs/16Uiu2HAmGx8Z6gdq7T5AQE54GMtqDhDFhizywTy1o28NJbAMMumH',
                blockPollingIntervalSeconds: 1,
                ethereumRPCMaxContentLength: 1024,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 30000,
                customContractAddresses:
                    '{"exchange":"0x48bacb9266a570d521063ef5dd96e61686dbe788",' +
                    '"devUtils":"0x38ef19fdf8e8415f18c307ed71967e19aac28ba1",' +
                    '"erc20Proxy":"0x1dc4c1cefef38a777b15aa20260a54e584b16c48",' +
                    '"erc721Proxy":"0x1d7022f5b17d2f8b695918fb48fa1089c9f85401",' +
                    '"erc1155Proxy":"0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"}',
                maxOrdersInStorage: 30000,
                customOrderFilter:
                    '{"properties":{"makerAssetData":{"const":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"}}}',
                ethereumChainID: 1337,
                web3Provider: undefined,
            });
        });
    });
    describe('#orderEventsHandlerToWrapperOrderEventsHandler', () => {});
    describe('#signedOrderToWrapperSignedOrder', () => {});
    describe('#wrapperAcceptedOrderInfoToAcceptedOrderInfo', () => {});
    describe('#wrapperContractEventsToContractEvents', () => {});
    describe('#wrapperGetOrdersResponse', () => {});
    describe('#wrapperOrderEventToOrderEvent', () => {});
    describe('#wrapperOrderInfoToOrderInfo', () => {});
    describe('#wrapperRejectedOrderInfoToRejectedOrderInfo', () => {});
    describe('#wrapperSignedOrderToSignedOrder', () => {});
    describe('#wrapperStatsToStats', () => {});
    describe('#wrapperValidationResultsToValidationResults', () => {});
});
// tslint:enable:no-empty
