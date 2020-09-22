/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
import { SignedOrder } from '@0x/order-utils';
import { BigNumber, providerUtils } from '@0x/utils';

import {
    AcceptedOrderInfo,
    Config,
    ContractEvent,
    ContractEventKind,
    ContractEventParameters,
    ERC1155ApprovalForAllEvent,
    ERC721ApprovalForAllEvent,
    ExchangeCancelEvent,
    GetOrdersResponse,
    LatestBlock,
    OrderEvent,
    OrderInfo,
    RejectedOrderInfo,
    Stats,
    ValidationResults,
    WrapperAcceptedOrderInfo,
    WrapperConfig,
    WrapperContractEvent,
    WrapperERC1155TransferBatchEvent,
    WrapperERC1155TransferSingleEvent,
    WrapperERC20ApprovalEvent,
    WrapperERC20TransferEvent,
    WrapperERC721ApprovalEvent,
    WrapperERC721TransferEvent,
    WrapperExchangeCancelUpToEvent,
    WrapperExchangeFillEvent,
    WrapperGetOrdersResponse,
    WrapperLatestBlock,
    WrapperOrderEvent,
    WrapperOrderInfo,
    WrapperRejectedOrderInfo,
    WrapperSignedOrder,
    WrapperStats,
    WrapperValidationResults,
    WrapperWethDepositEvent,
    WrapperWethWithdrawalEvent,
} from './types';

// NOTE(jalextowle): These functions are only exported so that it's easier to share code with
// the conversion tests. They should not be used outside of `0x-mesh/browser/ts/index.ts`.
// tslint:disable:completed-docs
export function configToWrapperConfig(config: Config): WrapperConfig {
    const bootstrapList = config.bootstrapList == null ? undefined : config.bootstrapList.join(',');
    const customContractAddresses =
        config.customContractAddresses == null ? undefined : JSON.stringify(config.customContractAddresses);
    const customOrderFilter = config.customOrderFilter == null ? undefined : JSON.stringify(config.customOrderFilter);
    const standardizedProvider =
        config.web3Provider == null ? undefined : providerUtils.standardizeOrThrow(config.web3Provider);
    return {
        ...config,
        bootstrapList,
        customContractAddresses,
        customOrderFilter,
        web3Provider: standardizedProvider,
    };
}

export function wrapperSignedOrderToSignedOrder(wrapperSignedOrder: WrapperSignedOrder): SignedOrder {
    return {
        ...wrapperSignedOrder,
        makerFee: new BigNumber(wrapperSignedOrder.makerFee),
        takerFee: new BigNumber(wrapperSignedOrder.takerFee),
        makerAssetAmount: new BigNumber(wrapperSignedOrder.makerAssetAmount),
        takerAssetAmount: new BigNumber(wrapperSignedOrder.takerAssetAmount),
        salt: new BigNumber(wrapperSignedOrder.salt),
        expirationTimeSeconds: new BigNumber(wrapperSignedOrder.expirationTimeSeconds),
    };
}

export function wrapperContractEventsToContractEvents(wrapperContractEvents: WrapperContractEvent[]): ContractEvent[] {
    const contractEvents: ContractEvent[] = [];
    if (wrapperContractEvents === null) {
        return contractEvents;
    }
    wrapperContractEvents.forEach(wrapperContractEvent => {
        const kind = wrapperContractEvent.kind as ContractEventKind;
        const rawParameters = wrapperContractEvent.parameters;
        let parameters: ContractEventParameters;
        switch (kind) {
            case ContractEventKind.ERC20TransferEvent:
                const erc20TransferEvent = rawParameters as WrapperERC20TransferEvent;
                parameters = {
                    from: erc20TransferEvent.from,
                    to: erc20TransferEvent.to,
                    value: new BigNumber(erc20TransferEvent.value),
                };
                break;
            case ContractEventKind.ERC20ApprovalEvent:
                const erc20ApprovalEvent = rawParameters as WrapperERC20ApprovalEvent;
                parameters = {
                    owner: erc20ApprovalEvent.owner,
                    spender: erc20ApprovalEvent.spender,
                    value: new BigNumber(erc20ApprovalEvent.value),
                };
                break;
            case ContractEventKind.ERC721TransferEvent:
                const erc721TransferEvent = rawParameters as WrapperERC721TransferEvent;
                parameters = {
                    from: erc721TransferEvent.from,
                    to: erc721TransferEvent.to,
                    tokenId: new BigNumber(erc721TransferEvent.tokenId),
                };
                break;
            case ContractEventKind.ERC721ApprovalEvent:
                const erc721ApprovalEvent = rawParameters as WrapperERC721ApprovalEvent;
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
                const erc1155TransferSingleEvent = rawParameters as WrapperERC1155TransferSingleEvent;
                parameters = {
                    operator: erc1155TransferSingleEvent.operator,
                    from: erc1155TransferSingleEvent.from,
                    to: erc1155TransferSingleEvent.to,
                    id: new BigNumber(erc1155TransferSingleEvent.id),
                    value: new BigNumber(erc1155TransferSingleEvent.value),
                };
                break;
            case ContractEventKind.ERC1155TransferBatchEvent:
                const erc1155TransferBatchEvent = rawParameters as WrapperERC1155TransferBatchEvent;
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
                const exchangeFillEvent = rawParameters as WrapperExchangeFillEvent;
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
                const exchangeCancelUpToEvent = rawParameters as WrapperExchangeCancelUpToEvent;
                parameters = {
                    makerAddress: exchangeCancelUpToEvent.makerAddress,
                    orderSenderAddress: exchangeCancelUpToEvent.orderSenderAddress,
                    orderEpoch: new BigNumber(exchangeCancelUpToEvent.orderEpoch),
                };
                break;
            case ContractEventKind.WethDepositEvent:
                const wethDepositEvent = rawParameters as WrapperWethDepositEvent;
                parameters = {
                    owner: wethDepositEvent.owner,
                    value: new BigNumber(wethDepositEvent.value),
                };
                break;
            case ContractEventKind.WethWithdrawalEvent:
                const wethWithdrawalEvent = rawParameters as WrapperWethWithdrawalEvent;
                parameters = {
                    owner: wethWithdrawalEvent.owner,
                    value: new BigNumber(wethWithdrawalEvent.value),
                };
                break;
            default:
                throw new Error(`Unrecognized ContractEventKind: ${kind}`);
        }
        const contractEvent: ContractEvent = {
            blockHash: wrapperContractEvent.blockHash,
            txHash: wrapperContractEvent.txHash,
            txIndex: wrapperContractEvent.txIndex,
            logIndex: wrapperContractEvent.logIndex,
            isRemoved: wrapperContractEvent.isRemoved,
            address: wrapperContractEvent.address,
            kind,
            parameters,
        };
        contractEvents.push(contractEvent);
    });
    return contractEvents;
}

export function signedOrderToWrapperSignedOrder(signedOrder: SignedOrder): WrapperSignedOrder {
    return {
        ...signedOrder,
        makerFee: signedOrder.makerFee.toString(),
        takerFee: signedOrder.takerFee.toString(),
        makerAssetAmount: signedOrder.makerAssetAmount.toString(),
        takerAssetAmount: signedOrder.takerAssetAmount.toString(),
        salt: signedOrder.salt.toString(),
        expirationTimeSeconds: signedOrder.expirationTimeSeconds.toString(),
    };
}

export function wrapperOrderEventToOrderEvent(wrapperOrderEvent: WrapperOrderEvent): OrderEvent {
    return {
        ...wrapperOrderEvent,
        timestampMs: new Date(wrapperOrderEvent.timestamp).getTime(),
        signedOrder: wrapperSignedOrderToSignedOrder(wrapperOrderEvent.signedOrder),
        fillableTakerAssetAmount: new BigNumber(wrapperOrderEvent.fillableTakerAssetAmount),
        contractEvents: wrapperContractEventsToContractEvents(wrapperOrderEvent.contractEvents),
    };
}

export function orderEventsHandlerToWrapperOrderEventsHandler(
    orderEventsHandler: (events: OrderEvent[]) => void,
): (events: WrapperOrderEvent[]) => void {
    return (wrapperOrderEvents: WrapperOrderEvent[]) => {
        const orderEvents = wrapperOrderEvents.map(wrapperOrderEventToOrderEvent);
        orderEventsHandler(orderEvents);
    };
}

export function wrapperValidationResultsToValidationResults(
    wrapperValidationResults: WrapperValidationResults,
): ValidationResults {
    return {
        accepted: wrapperValidationResults.accepted.map(wrapperAcceptedOrderInfoToAcceptedOrderInfo),
        rejected: wrapperValidationResults.rejected.map(wrapperRejectedOrderInfoToRejectedOrderInfo),
    };
}

export function wrapperAcceptedOrderInfoToAcceptedOrderInfo(
    wrapperAcceptedOrderInfo: WrapperAcceptedOrderInfo,
): AcceptedOrderInfo {
    return {
        ...wrapperAcceptedOrderInfo,
        signedOrder: wrapperSignedOrderToSignedOrder(wrapperAcceptedOrderInfo.signedOrder),
        fillableTakerAssetAmount: new BigNumber(wrapperAcceptedOrderInfo.fillableTakerAssetAmount),
    };
}

export function wrapperRejectedOrderInfoToRejectedOrderInfo(
    wrapperRejectedOrderInfo: WrapperRejectedOrderInfo,
): RejectedOrderInfo {
    return {
        ...wrapperRejectedOrderInfo,
        signedOrder: wrapperSignedOrderToSignedOrder(wrapperRejectedOrderInfo.signedOrder),
    };
}

export function wrapperStatsToStats(wrapperStats: WrapperStats): Stats {
    return {
        ...wrapperStats,
        latestBlock: wrapperLatestBlockToLatestBlock(wrapperStats.latestBlock),
        startOfCurrentUTCDay: new Date(wrapperStats.startOfCurrentUTCDay),
        maxExpirationTime: new BigNumber(wrapperStats.maxExpirationTime),
    };
}

export function wrapperLatestBlockToLatestBlock(wrapperLatestBlock: WrapperLatestBlock): LatestBlock {
    return {
        ...wrapperLatestBlock,
        number: new BigNumber(wrapperLatestBlock.number),
    };
}

export function wrapperGetOrdersResponseToGetOrdersResponse(
    wrapperGetOrdersResponse: WrapperGetOrdersResponse,
): GetOrdersResponse {
    return {
        ...wrapperGetOrdersResponse,
        timestamp: new Date(wrapperGetOrdersResponse.timestamp).getTime(),
        ordersInfos: wrapperGetOrdersResponse.ordersInfos.map(wrapperOrderInfoToOrderInfo),
    };
}

export function wrapperOrderInfoToOrderInfo(wrapperOrderInfo: WrapperOrderInfo): OrderInfo {
    return {
        ...wrapperOrderInfo,
        fillableTakerAssetAmount: new BigNumber(wrapperOrderInfo.fillableTakerAssetAmount),
        signedOrder: wrapperSignedOrderToSignedOrder(wrapperOrderInfo.signedOrder),
    };
}
// tslint:enable:completed-docs
