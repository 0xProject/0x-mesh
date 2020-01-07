import * as chai from 'chai';

import {
    configToWrapperConfig,
    orderEventsHandlerToWrapperOrderEventsHandler,
    signedOrderToWrapperSignedOrder,
    wrapperAcceptedOrderInfoToAcceptedOrderInfo,
    wrapperContractEventsToContractEvents,
    wrapperOrderEventToOrderEvent,
    wrapperRejectedOrderInfoToRejectedOrderInfo,
    wrapperSignedOrderToSignedOrder,
    wrapperValidationResultsToValidationResults,
} from '../ts/type_conversion';
import { Verbosity } from '../ts/types';

const { expect } = chai;

describe('Type Conversion Unit Tests', () => {
    const ethereumRPCURL = 'https://localhost:8545';
    const ethereumChainID = 1337;

    describe('configToWrapperConfig', () => {
        it('encodes a config with a single bootstrap peer specified', async () => {
            const bootstrapList = [
                '/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF',
            ];
            const customContractAddresses = {
                exchange: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
                devUtils: '0x38ef19fdf8e8415f18c307ed71967e19aac28ba1',
                erc20Proxy: '0x1dc4c1cefef38a777b15aa20260a54e584b16c48',
                erc721Proxy: '0x1d7022f5b17d2f8b695918fb48fa1089c9f85401',
                erc1155Proxy: '0x1d7022f5b17d2f8b695918fb48fa1089c9f85401',
            };
            const wrapperConfig = configToWrapperConfig({
                verbosity: Verbosity.Debug,
                ethereumRPCURL,
                ethereumChainID,
                useBootstrapList: true,
                bootstrapList,
                blockPollingIntervalSeconds: 10,
                ethereumRPCMaxContentLength: 524288,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 100000,
                ethereumRPCMaxRequestsPerSecond: 30,
                customContractAddresses,
                maxOrdersInStorage: 100000,
            });
            expect(wrapperConfig).to.be.deep.eq({
                verbosity: Verbosity.Debug as number,
                ethereumRPCURL,
                ethereumChainID,
                useBootstrapList: true,
                bootstrapList: bootstrapList.join(','),
                blockPollingIntervalSeconds: 10,
                ethereumRPCMaxContentLength: 524288,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 100000,
                ethereumRPCMaxRequestsPerSecond: 30,
                customContractAddresses: JSON.stringify(customContractAddresses),
                maxOrdersInStorage: 100000,
            });
        });

        it('encodes a config with multiple bootstrap peers specified', async () => {
            const bootstrapList = [
                '/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF',
                '/ip4/3.214.190.68/tcp/60559/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumG',
            ];
            const customContractAddresses = {
                exchange: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
                devUtils: '0x38ef19fdf8e8415f18c307ed71967e19aac28ba1',
                erc20Proxy: '0x1dc4c1cefef38a777b15aa20260a54e584b16c48',
                erc721Proxy: '0x1d7022f5b17d2f8b695918fb48fa1089c9f85401',
                erc1155Proxy: '0x1d7022f5b17d2f8b695918fb48fa1089c9f85401',
            };
            const wrapperConfig = configToWrapperConfig({
                verbosity: Verbosity.Debug,
                ethereumRPCURL,
                ethereumChainID,
                useBootstrapList: true,
                bootstrapList,
                blockPollingIntervalSeconds: 10,
                ethereumRPCMaxContentLength: 524288,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 100000,
                ethereumRPCMaxRequestsPerSecond: 30,
                customContractAddresses,
                maxOrdersInStorage: 100000,
            });
            expect(wrapperConfig).to.be.deep.eq({
                verbosity: Verbosity.Debug as number,
                ethereumRPCURL,
                ethereumChainID,
                useBootstrapList: true,
                bootstrapList: bootstrapList.join(','),
                blockPollingIntervalSeconds: 10,
                ethereumRPCMaxContentLength: 524288,
                enableEthereumRPCRateLimiting: true,
                ethereumRPCMaxRequestsPer24HrUTC: 100000,
                ethereumRPCMaxRequestsPerSecond: 30,
                customContractAddresses: JSON.stringify(customContractAddresses),
                maxOrdersInStorage: 100000,
            });
        });

        it('encodes a config file with none of the optional fields', async () => {
            const wrapperConfig = configToWrapperConfig({
                ethereumRPCURL,
                ethereumChainID,
            });
            // NOTE(jalextowle): Chai distinguishes betwen a value that has been
            // explicitly marked as undefined and a value that was never defined.
            // The only values that will be explicitly marked
            // as undefined are `bootstrapList` and `customContractAddresses`.
            // The `to.be.undefined` assertion is how we will check that the other
            // values are in fact undefined.
            // tslint:disable:no-unused-expression
            expect(wrapperConfig.verbosity).to.be.undefined;
            expect(wrapperConfig.useBootstrapList).to.be.undefined;
            expect(wrapperConfig.blockPollingIntervalSeconds).to.be.undefined;
            expect(wrapperConfig.ethereumRPCMaxContentLength).to.be.undefined;
            expect(wrapperConfig.enableEthereumRPCRateLimiting).to.be.undefined;
            expect(wrapperConfig.ethereumRPCMaxRequestsPer24HrUTC).to.be.undefined;
            expect(wrapperConfig.ethereumRPCMaxRequestsPerSecond).to.be.undefined;
            expect(wrapperConfig.maxOrdersInStorage).to.be.undefined;
            // tslint:enable:no-unused-expression
            expect(wrapperConfig).to.be.deep.eq({
                ethereumRPCURL,
                ethereumChainID,
                bootstrapList: undefined,
                customContractAddresses: undefined,
            });
        });
    });
});
