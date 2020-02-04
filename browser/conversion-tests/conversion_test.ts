import { BigNumber, hexUtils, logUtils } from '@0x/utils';

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
} from '../ts/encoding';
import { wasmBuffer } from './generated/test_wasm_buffer';
import {
    ContractEvent,
    ERC1155ApprovalForAllEvent,
    ERC1155TransferBatchEvent,
    ERC1155TransferSingleEvent,
    ERC20ApprovalEvent,
    ERC20TransferEvent,
    ERC721ApprovalEvent,
    ERC721ApprovalForAllEvent,
    ERC721TransferEvent,
    ExchangeCancelEvent,
    ExchangeCancelUpToEvent,
    ExchangeFillEvent,
    WethDepositEvent,
    WethWithdrawalEvent,
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
    WrapperOrderEvent,
    WrapperSignedOrder,
    WrapperStats,
    WrapperValidationResults,
    WrapperWethDepositEvent,
    WrapperWethWithdrawalEvent,
} from '../ts/types';
import '../ts/wasm_exec';

interface ConversionTestCase {
    contractEventsAsync: () => Promise<WrapperContractEvent[]>;
    getOrdersResponseAsync: () => Promise<WrapperGetOrdersResponse[]>;
    orderEventsAsync: () => Promise<WrapperOrderEvent[]>;
    signedOrdersAsync: () => Promise<WrapperSignedOrder[]>;
    statsAsync: () => Promise<WrapperStats[]>;
    validationResultsAsync: () => Promise<WrapperValidationResults[]>;
}

// The Go code sets certain global values and this is our only way of
// interacting with it. Define those values and their types here.
declare global {
    // Defined in wasm_exec.ts
    class Go {
        public importObject: any;
        public run(instance: WebAssembly.Instance): void;
    }

    // Define variables that are defined in `browser/go/conversion-test.go`
    const conversionTestCases: ConversionTestCase;
}

// The interval (in milliseconds) to check whether Wasm is done loading.
const wasmLoadCheckIntervalMs = 100;

// We use a global variable to track whether the Wasm code has finished loading.
let isWasmLoaded = false;
const loadEventName = '0xmeshtest';
window.addEventListener(loadEventName, () => {
    isWasmLoaded = true;
});

// Start compiling the WebAssembly as soon as the script is loaded. This lets
// us initialize as quickly as possible.
const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject)
    .then(module => {
        go.run(module.instance);
    })
    .catch(err => {
        // tslint:disable-next-line no-console
        console.error('Could not load Wasm');
        // tslint:disable-next-line no-console
        console.error(err);
        // If the Wasm bytecode didn't compile, Mesh won't work. We have no
        // choice but to throw an error.
        setImmediate(() => {
            throw err;
        });
    });

/*********************** Tests ***********************/
// tslint:disable:custom-no-magic-numbers
// tslint:disable:no-console

(async () => {
    await waitForLoadAsync();
    const contractEvents = await conversionTestCases.contractEventsAsync();
    testContractEvents(contractEvents);
    const getOrdersResponse = await conversionTestCases.getOrdersResponseAsync();
    testGetOrdersResponse(getOrdersResponse);
    const orderEvents = await conversionTestCases.orderEventsAsync();
    testOrderEvents(orderEvents);
    const signedOrders = await conversionTestCases.signedOrdersAsync();
    testSignedOrders(signedOrders);
    const stats = await conversionTestCases.statsAsync();
    testStats(stats);
    const validationResults = await conversionTestCases.validationResultsAsync();
    testValidationResults(validationResults);

    // This special #jsFinished div is used to signal the headless Chrome driver
    // that the JavaScript code is done running.
    const finishedDiv = document.createElement('div');
    finishedDiv.setAttribute('id', 'jsFinished');
    document.querySelector('body')!.appendChild(finishedDiv); // tslint:disable-line:no-non-null-assertion
})().catch(err => {
    throw err;
});

function prettyPrintTestCase(name: string, description: string): (section: string, value: boolean) => void {
    return (section: string, value: boolean) => {
        console.log(`(${name} | ${description} | ${section}): ${value}`);
    };
}

function testContractEvents(contractEvents: WrapperContractEvent[]): void {
    // ERC20ApprovalEvent
    let printer = prettyPrintTestCase('contractEventTest', 'ERC20ApprovalEvent');
    testContractEventPrelude(printer, contractEvents[0]);
    printer('kind', contractEvents[0].kind === 'ERC20ApprovalEvent');
    const erc20ApprovalParams = contractEvents[0].parameters as WrapperERC20ApprovalEvent;
    printer('parameter | owner', erc20ApprovalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | spender', erc20ApprovalParams.spender === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', erc20ApprovalParams.value === '1000');

    // ERC20TransferEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC20TransferEvent');
    testContractEventPrelude(printer, contractEvents[1]);
    printer('kind', contractEvents[1].kind === 'ERC20TransferEvent');
    const erc20TransferParams = contractEvents[1].parameters as WrapperERC20TransferEvent;
    printer('parameter | from', erc20TransferParams.from === hexUtils.leftPad('0x4', 20));
    printer('parameter | to', erc20TransferParams.to === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', erc20TransferParams.value === '1000');

    // ERC721ApprovalEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC721ApprovalEvent');
    testContractEventPrelude(printer, contractEvents[2]);
    printer('kind', contractEvents[2].kind === 'ERC721ApprovalEvent');
    const erc721ApprovalParams = contractEvents[2].parameters as WrapperERC721ApprovalEvent;
    printer('parameter | owner', erc721ApprovalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | approved', erc721ApprovalParams.approved === hexUtils.leftPad('0x5', 20));
    printer('parameter | tokenId', erc721ApprovalParams.tokenId === '1');

    // ERC721ApprovalForAllEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC721ApprovalForAllEvent');
    testContractEventPrelude(printer, contractEvents[3]);
    printer('kind', contractEvents[3].kind === 'ERC721ApprovalForAllEvent');
    const erc721ApprovalForAllParams = contractEvents[3].parameters as ERC721ApprovalForAllEvent;
    printer('parameter | owner', erc721ApprovalForAllParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | operator', erc721ApprovalForAllParams.operator === hexUtils.leftPad('0x5', 20));
    printer('parameter | approved', erc721ApprovalForAllParams.approved);

    // ERC721TransferEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC721TransferEvent');
    testContractEventPrelude(printer, contractEvents[4]);
    printer('kind', contractEvents[4].kind === 'ERC721TransferEvent');
    const erc721TransferParams = contractEvents[4].parameters as WrapperERC721TransferEvent;
    printer('parameter | from', erc721TransferParams.from === hexUtils.leftPad('0x4', 20));
    printer('parameter | to', erc721TransferParams.to === hexUtils.leftPad('0x5', 20));
    printer('parameter | tokenId', erc721TransferParams.tokenId === '1');

    // ERC1155ApprovalForAllEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC1155ApprovalForAllEvent');
    testContractEventPrelude(printer, contractEvents[5]);
    printer('kind', contractEvents[5].kind === 'ERC1155ApprovalForAllEvent');
    const erc1155ApprovalForAllParams = contractEvents[5].parameters as ERC1155ApprovalForAllEvent;
    printer('parameter | owner', erc1155ApprovalForAllParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | operator', erc1155ApprovalForAllParams.operator === hexUtils.leftPad('0x5', 20));
    printer('parameter | approved', !erc1155ApprovalForAllParams.approved);

    // ERC1155TransferSingleEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC1155TransferSingleEvent');
    testContractEventPrelude(printer, contractEvents[6]);
    printer('kind', contractEvents[6].kind === 'ERC1155TransferSingleEvent');
    const erc1155TransferSingleParams = contractEvents[6].parameters as WrapperERC1155TransferSingleEvent;
    printer('parameter | operator', erc1155TransferSingleParams.operator === hexUtils.leftPad('0x4', 20));
    printer('parameter | from', erc1155TransferSingleParams.from === hexUtils.leftPad('0x5', 20));
    printer('parameter | to', erc1155TransferSingleParams.to === hexUtils.leftPad('0x6', 20));
    // FIXME(jalextowle): Investigate whether or not this style of encoding also occurs with large numbers
    printer('parameter | id', erc1155TransferSingleParams.id === '1');
    printer('parameter | value', erc1155TransferSingleParams.value === '100');

    // ERC1155TransferBatchEvent
    printer = prettyPrintTestCase('contractEventTest', 'ERC1155TransferBatchEvent');
    testContractEventPrelude(printer, contractEvents[7]);
    printer('kind', contractEvents[7].kind === 'ERC1155TransferBatchEvent');
    const erc1155TransferBatchParams = contractEvents[7].parameters as WrapperERC1155TransferBatchEvent;
    printer('parameter | operator', erc1155TransferBatchParams.operator === hexUtils.leftPad('0x4', 20));
    printer('parameter | from', erc1155TransferBatchParams.from === hexUtils.leftPad('0x5', 20));
    printer('parameter | to', erc1155TransferBatchParams.to === hexUtils.leftPad('0x6', 20));
    printer(
        'parameter | ids',
        erc1155TransferBatchParams.ids.length === 1 && erc1155TransferBatchParams.ids[0] === '1',
    );
    printer(
        'parameter | values',
        erc1155TransferBatchParams.values.length === 1 && erc1155TransferBatchParams.values[0] === '100',
    );

    // ExchangeFillEvent
    printer = prettyPrintTestCase('contractEventTest', 'ExchangeFillEvent');
    testContractEventPrelude(printer, contractEvents[8]);
    printer('kind', contractEvents[8].kind === 'ExchangeFillEvent');
    const exchangeFillParams = contractEvents[8].parameters as WrapperExchangeFillEvent;
    printer('parameter | makerAddress', exchangeFillParams.makerAddress === hexUtils.leftPad('0x4', 20));
    printer('parameter | takerAddress', exchangeFillParams.takerAddress === hexUtils.leftPad('0x0', 20));
    printer('parameter | senderAddress', exchangeFillParams.senderAddress === hexUtils.leftPad('0x5', 20));
    printer('parameter | feeRecipientAddress', exchangeFillParams.feeRecipientAddress === hexUtils.leftPad('0x6', 20));
    printer('parameter | makerAssetFilledAmount', exchangeFillParams.makerAssetFilledAmount === '456');
    printer('parameter | takerAssetFilledAmount', exchangeFillParams.takerAssetFilledAmount === '654');
    printer('parameter | makerFeePaid', exchangeFillParams.makerFeePaid === '12');
    printer('parameter | takerFeePaid', exchangeFillParams.takerFeePaid === '21');
    printer('parameter | protocolFeePaid', exchangeFillParams.protocolFeePaid === '150000');
    printer('parameter | orderHash', exchangeFillParams.orderHash === hexUtils.leftPad('0x7', 32));
    printer('parameter | makerAssetData', exchangeFillParams.makerAssetData === '0x');
    printer('parameter | takerAssetData', exchangeFillParams.takerAssetData === '0x');
    printer('parameter | makerFeeAssetData', exchangeFillParams.makerFeeAssetData === '0x');
    printer('parameter | takerFeeAssetData', exchangeFillParams.takerFeeAssetData === '0x');

    // ExchangeCancelEvent
    printer = prettyPrintTestCase('contractEventTest', 'ExchangeCancelEvent');
    testContractEventPrelude(printer, contractEvents[9]);
    printer('kind', contractEvents[9].kind === 'ExchangeCancelEvent');
    const exchangeCancelParams = contractEvents[9].parameters as ExchangeCancelEvent;
    printer('parameter | makerAddress', exchangeCancelParams.makerAddress === hexUtils.leftPad('0x4', 20));
    printer('parameter | senderAddress', exchangeCancelParams.senderAddress === hexUtils.leftPad('0x5', 20));
    printer(
        'parameter | feeRecipientAddress',
        exchangeCancelParams.feeRecipientAddress === hexUtils.leftPad('0x6', 20),
    );
    printer('parameter | orderHash', exchangeCancelParams.orderHash === hexUtils.leftPad('0x7', 32));
    printer('parameter | makerAssetData', exchangeCancelParams.makerAssetData === '0x');
    printer('parameter | takerAssetData', exchangeCancelParams.takerAssetData === '0x');

    // ExchangeCancelUpToEvent
    printer = prettyPrintTestCase('contractEventTest', 'ExchangeCancelUpToEvent');
    testContractEventPrelude(printer, contractEvents[10]);
    printer('kind', contractEvents[10].kind === 'ExchangeCancelUpToEvent');
    const exchangeCancelUpToParams = contractEvents[10].parameters as WrapperExchangeCancelUpToEvent;
    printer('parameter | makerAddress', exchangeCancelUpToParams.makerAddress === hexUtils.leftPad('0x4', 20));
    printer(
        'parameter | orderSenderAddress',
        exchangeCancelUpToParams.orderSenderAddress === hexUtils.leftPad('0x5', 20),
    );
    printer('parameter | orderEpoch', exchangeCancelUpToParams.orderEpoch === '50');

    // WethDepositEvent
    printer = prettyPrintTestCase('contractEventTest', 'WethDepositEvent');
    testContractEventPrelude(printer, contractEvents[11]);
    printer('kind', contractEvents[11].kind === 'WethDepositEvent');
    const wethDepositParams = contractEvents[11].parameters as WrapperWethDepositEvent;
    printer('parameter | owner', wethDepositParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | value', wethDepositParams.value === '150000');

    // WethWithdrawalEvent
    printer = prettyPrintTestCase('contractEventTest', 'WethWithdrawalEvent');
    testContractEventPrelude(printer, contractEvents[12]);
    printer('kind', contractEvents[12].kind === 'WethWithdrawalEvent');
    const wethWithdrawalParams = contractEvents[12].parameters as WrapperWethWithdrawalEvent;
    printer('parameter | owner', wethWithdrawalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | value', wethWithdrawalParams.value === '150000');

    // FooBarBaz
    printer = prettyPrintTestCase('contractEventTest', 'FooBarBazEvent');
    testContractEventPrelude(printer, contractEvents[13]);
    printer('kind', contractEvents[13].kind === 'FooBarBazEvent');
    const fooBarBazParams = contractEvents[13].parameters as any;
    printer('parameter | owner', fooBarBazParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | spender', fooBarBazParams.spender === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', fooBarBazParams.value === '1');
}

function testContractEventPrelude(
    printer: (section: string, value: boolean) => void,
    contractEvent: WrapperContractEvent,
): void {
    printer('blockHash', contractEvent.blockHash === hexUtils.leftPad(1, 32));
    printer('txHash', contractEvent.txHash === hexUtils.leftPad(2, 32));
    printer('txIndex', contractEvent.txIndex === 123);
    printer('logIndex', contractEvent.logIndex === 321);
    printer('isRemoved', !contractEvent.isRemoved);
    printer('address', contractEvent.address === hexUtils.leftPad(3, 20));
}

function testGetOrdersResponse(getOrdersResponse: WrapperGetOrdersResponse[]): void {
    let printer = prettyPrintTestCase('getOrdersResponseTest', 'emptyOrderInfo');
    printer('snapshotID', getOrdersResponse[0].snapshotID === '208c81f9-6f8d-44aa-b6ea-0a3276ec7318');
    printer('snapshotTimestamp', getOrdersResponse[0].snapshotTimestamp === '2006-01-01T00:00:00Z');
    printer('orderInfo | length', getOrdersResponse[0].ordersInfos.length === 0);

    printer = prettyPrintTestCase('getOrdersResponseTest', 'oneOrderInfo');
    printer('snapshotID', getOrdersResponse[1].snapshotID === '208c81f9-6f8d-44aa-b6ea-0a3276ec7318');
    printer('snapshotTimestamp', getOrdersResponse[1].snapshotTimestamp === '2006-01-01T00:00:00Z');
    printer('orderInfo | length', getOrdersResponse[1].ordersInfos.length === 1);
    printer('orderInfo | orderHash', getOrdersResponse[1].ordersInfos[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('orderInfo | signedOrder | chainId', getOrdersResponse[1].ordersInfos[0].signedOrder.chainId === 1337);
    printer(
        'orderInfo | signedOrder | makerAddress',
        getOrdersResponse[1].ordersInfos[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'orderInfo | signedOrder | takerAddress',
        getOrdersResponse[1].ordersInfos[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'orderInfo | signedOrder | senderAddress',
        getOrdersResponse[1].ordersInfos[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'orderInfo | signedOrder | feeRecipientAddress',
        getOrdersResponse[1].ordersInfos[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'orderInfo | signedOrder | exchangeAddress',
        getOrdersResponse[1].ordersInfos[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'orderInfo | signedOrder | makerAssetData',
        getOrdersResponse[1].ordersInfos[0].signedOrder.makerAssetData ===
            '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer(
        'orderInfo | signedOrder | makerAssetAmount',
        getOrdersResponse[1].ordersInfos[0].signedOrder.makerAssetAmount === '123456789',
    );
    printer(
        'orderInfo | signedOrder | makerFeeAssetData',
        getOrdersResponse[1].ordersInfos[0].signedOrder.makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('orderInfo | signedOrder | makerFee', getOrdersResponse[1].ordersInfos[0].signedOrder.makerFee === '89');
    printer(
        'orderInfo | signedOrder | takerAssetData',
        getOrdersResponse[1].ordersInfos[0].signedOrder.takerAssetData ===
            '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer(
        'orderInfo | signedOrder | takerAssetAmount',
        getOrdersResponse[1].ordersInfos[0].signedOrder.takerAssetAmount === '987654321',
    );
    printer(
        'orderInfo | signedOrder | takerFeeAssetData',
        getOrdersResponse[1].ordersInfos[0].signedOrder.takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('orderInfo | signedOrder | takerFee', getOrdersResponse[1].ordersInfos[0].signedOrder.takerFee === '12');
    printer(
        'orderInfo | signedOrder | expirationTimeSeconds',
        getOrdersResponse[1].ordersInfos[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('orderInfo | signedOrder | salt', getOrdersResponse[1].ordersInfos[0].signedOrder.salt === '1532559225');
    printer(
        'orderInfo | fillableTakerAssetAmount',
        getOrdersResponse[1].ordersInfos[0].fillableTakerAssetAmount === '987654321',
    );

    printer = prettyPrintTestCase('getOrdersResponseTest', 'twoOrderInfos');
    printer('snapshotID', getOrdersResponse[2].snapshotID === '208c81f9-6f8d-44aa-b6ea-0a3276ec7318');
    printer('snapshotTimestamp', getOrdersResponse[2].snapshotTimestamp === '2006-01-01T00:00:00Z');
    printer('orderInfo | length', getOrdersResponse[2].ordersInfos.length === 2);
    printer('orderInfo | orderHash', getOrdersResponse[2].ordersInfos[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('orderInfo | signedOrder | chainId', getOrdersResponse[2].ordersInfos[0].signedOrder.chainId === 1337);
    printer(
        'orderInfo | signedOrder | makerAddress',
        getOrdersResponse[2].ordersInfos[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'orderInfo | signedOrder | takerAddress',
        getOrdersResponse[2].ordersInfos[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'orderInfo | signedOrder | senderAddress',
        getOrdersResponse[2].ordersInfos[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'orderInfo | signedOrder | feeRecipientAddress',
        getOrdersResponse[2].ordersInfos[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'orderInfo | signedOrder | exchangeAddress',
        getOrdersResponse[2].ordersInfos[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'orderInfo | signedOrder | makerAssetData',
        getOrdersResponse[2].ordersInfos[0].signedOrder.makerAssetData === '0x',
    );
    printer(
        'orderInfo | signedOrder | makerAssetAmount',
        getOrdersResponse[2].ordersInfos[0].signedOrder.makerAssetAmount === '0',
    );
    printer(
        'orderInfo | signedOrder | makerFeeAssetData',
        getOrdersResponse[2].ordersInfos[0].signedOrder.makerFeeAssetData === '0x',
    );
    printer('orderInfo | signedOrder | makerFee', getOrdersResponse[2].ordersInfos[0].signedOrder.makerFee === '0');
    printer(
        'orderInfo | signedOrder | takerAssetData',
        getOrdersResponse[2].ordersInfos[0].signedOrder.takerAssetData === '0x',
    );
    printer(
        'orderInfo | signedOrder | takerAssetAmount',
        getOrdersResponse[2].ordersInfos[0].signedOrder.takerAssetAmount === '0',
    );
    printer(
        'orderInfo | signedOrder | takerFeeAssetData',
        getOrdersResponse[2].ordersInfos[0].signedOrder.takerFeeAssetData === '0x',
    );
    printer('orderInfo | signedOrder | takerFee', getOrdersResponse[2].ordersInfos[0].signedOrder.takerFee === '0');
    printer(
        'orderInfo | signedOrder | expirationTimeSeconds',
        getOrdersResponse[2].ordersInfos[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('orderInfo | signedOrder | salt', getOrdersResponse[2].ordersInfos[0].signedOrder.salt === '1532559225');
    printer(
        'orderInfo | fillableTakerAssetAmount',
        getOrdersResponse[2].ordersInfos[0].fillableTakerAssetAmount === '0',
    );
    printer('orderInfo | orderHash', getOrdersResponse[2].ordersInfos[1].orderHash === hexUtils.leftPad('0x1', 32));
    printer('orderInfo | signedOrder | chainId', getOrdersResponse[2].ordersInfos[1].signedOrder.chainId === 1337);
    printer(
        'orderInfo | signedOrder | makerAddress',
        getOrdersResponse[2].ordersInfos[1].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'orderInfo | signedOrder | takerAddress',
        getOrdersResponse[2].ordersInfos[1].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'orderInfo | signedOrder | senderAddress',
        getOrdersResponse[2].ordersInfos[1].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'orderInfo | signedOrder | feeRecipientAddress',
        getOrdersResponse[2].ordersInfos[1].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'orderInfo | signedOrder | exchangeAddress',
        getOrdersResponse[2].ordersInfos[1].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'orderInfo | signedOrder | makerAssetData',
        getOrdersResponse[2].ordersInfos[1].signedOrder.makerAssetData ===
            '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer(
        'orderInfo | signedOrder | makerAssetAmount',
        getOrdersResponse[2].ordersInfos[1].signedOrder.makerAssetAmount === '123456789',
    );
    printer(
        'orderInfo | signedOrder | makerFeeAssetData',
        getOrdersResponse[2].ordersInfos[1].signedOrder.makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('orderInfo | signedOrder | makerFee', getOrdersResponse[2].ordersInfos[1].signedOrder.makerFee === '89');
    printer(
        'orderInfo | signedOrder | takerAssetData',
        getOrdersResponse[2].ordersInfos[1].signedOrder.takerAssetData ===
            '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer(
        'orderInfo | signedOrder | takerAssetAmount',
        getOrdersResponse[2].ordersInfos[1].signedOrder.takerAssetAmount === '987654321',
    );
    printer(
        'orderInfo | signedOrder | takerFeeAssetData',
        getOrdersResponse[2].ordersInfos[1].signedOrder.takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('orderInfo | signedOrder | takerFee', getOrdersResponse[2].ordersInfos[1].signedOrder.takerFee === '12');
    printer(
        'orderInfo | signedOrder | expirationTimeSeconds',
        getOrdersResponse[2].ordersInfos[1].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('orderInfo | signedOrder | salt', getOrdersResponse[2].ordersInfos[1].signedOrder.salt === '1532559225');
    printer(
        'orderInfo | fillableTakerAssetAmount',
        getOrdersResponse[2].ordersInfos[1].fillableTakerAssetAmount === '987654321',
    );
}

function testOrderEvents(orderEvents: WrapperOrderEvent[]): void {
    let printer = prettyPrintTestCase('orderEventTest', 'EmptyContractEvents');
    printer('timestamp', orderEvents[0].timestamp === '2006-01-01T00:00:00Z');
    printer('orderHash', orderEvents[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('endState', orderEvents[0].endState === 'ADDED');
    printer('fillableTakerAssetAmount', orderEvents[0].fillableTakerAssetAmount === '1');
    printer = prettyPrintTestCase('orderEventTest', 'EmptyContractEvents | signedOrder');
    printer('chainId', orderEvents[0].signedOrder.chainId === 1337);
    printer('makerAddress', orderEvents[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20));
    printer('takerAddress', orderEvents[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20));
    printer('senderAddress', orderEvents[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20));
    printer('feeRecipientAddress', orderEvents[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20));
    printer('exchangeAddress', orderEvents[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20));
    printer('makerAssetData', orderEvents[0].signedOrder.makerAssetData === '0x');
    printer('makerAssetAmount', orderEvents[0].signedOrder.makerAssetAmount === '0');
    printer('makerFeeAssetData', orderEvents[0].signedOrder.makerFeeAssetData === '0x');
    printer('makerFee', orderEvents[0].signedOrder.makerFee === '0');
    printer('takerAssetData', orderEvents[0].signedOrder.takerAssetData === '0x');
    printer('takerAssetAmount', orderEvents[0].signedOrder.takerAssetAmount === '0');
    printer('takerFeeAssetData', orderEvents[0].signedOrder.takerFeeAssetData === '0x');
    printer('takerFee', orderEvents[0].signedOrder.takerFee === '0');
    printer('expirationTimeSeconds', orderEvents[0].signedOrder.expirationTimeSeconds === '10000000000');
    printer('salt', orderEvents[0].signedOrder.salt === '1532559225');
    printer = prettyPrintTestCase('orderEventTest', 'EmptyContractEvents | contractEvents');
    printer('length', orderEvents[0].contractEvents.length === 0);

    printer = prettyPrintTestCase('orderEventTest', 'ExchangeFillContractEvent');
    printer('timestamp', orderEvents[1].timestamp === '2006-01-01T01:01:01Z');
    printer('orderHash', orderEvents[1].orderHash === hexUtils.leftPad('0x1', 32));
    printer('endState', orderEvents[1].endState === 'FILLED');
    printer('fillableTakerAssetAmount', orderEvents[1].fillableTakerAssetAmount === '0');
    printer = prettyPrintTestCase('orderEventTest', 'ExchangeFillContractEvent | signedOrder');
    printer('chainId', orderEvents[1].signedOrder.chainId === 1337);
    printer('makerAddress', orderEvents[1].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20));
    printer('takerAddress', orderEvents[1].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20));
    printer('senderAddress', orderEvents[1].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20));
    printer('feeRecipientAddress', orderEvents[1].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20));
    printer('exchangeAddress', orderEvents[1].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20));
    printer(
        'makerAssetData',
        orderEvents[1].signedOrder.makerAssetData ===
            '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer('makerAssetAmount', orderEvents[1].signedOrder.makerAssetAmount === '123456789');
    printer(
        'makerFeeAssetData',
        orderEvents[1].signedOrder.makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('makerFee', orderEvents[1].signedOrder.makerFee === '89');
    printer(
        'takerAssetData',
        orderEvents[1].signedOrder.takerAssetData ===
            '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer('takerAssetAmount', orderEvents[1].signedOrder.takerAssetAmount === '987654321');
    printer(
        'takerFeeAssetData',
        orderEvents[1].signedOrder.takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('takerFee', orderEvents[1].signedOrder.takerFee === '12');
    printer('expirationTimeSeconds', orderEvents[1].signedOrder.expirationTimeSeconds === '10000000000');
    printer('salt', orderEvents[1].signedOrder.salt === '1532559225');
    printer = prettyPrintTestCase('orderEventTest', 'ExchangeFillContractEvent | contractEvents');
    printer('length', orderEvents[1].contractEvents.length === 1);
    printer = prettyPrintTestCase('orderEventTest', 'ExchangeFillContractEvent | contractEvents');
    printer('blockHash', orderEvents[1].contractEvents[0].blockHash === hexUtils.leftPad('0x1', 32));
    printer('txHash', orderEvents[1].contractEvents[0].txHash === hexUtils.leftPad('0x2', 32));
    printer('txIndex', orderEvents[1].contractEvents[0].txIndex === 123);
    printer('logIndex', orderEvents[1].contractEvents[0].logIndex === 321);
    printer('isRemoved', orderEvents[1].contractEvents[0].isRemoved === false);
    printer('address', orderEvents[1].contractEvents[0].address === hexUtils.leftPad('0x5', 20));
    printer('kind', orderEvents[1].contractEvents[0].kind === 'ExchangeFillEvent');
}

function testSignedOrders(signedOrders: WrapperSignedOrder[]): void {
    let printer = prettyPrintTestCase('signedOrderTest', 'NullAssetData');
    printer('chainId', signedOrders[0].chainId === 1337);
    printer('makerAddress', signedOrders[0].makerAddress === hexUtils.leftPad('0x1', 20));
    printer('takerAddress', signedOrders[0].takerAddress === hexUtils.leftPad('0x2', 20));
    printer('senderAddress', signedOrders[0].senderAddress === hexUtils.leftPad('0x3', 20));
    printer('feeRecipientAddress', signedOrders[0].feeRecipientAddress === hexUtils.leftPad('0x4', 20));
    printer('exchangeAddress', signedOrders[0].exchangeAddress === hexUtils.leftPad('0x5', 20));
    printer('makerAssetData', signedOrders[0].makerAssetData === '0x');
    printer('makerAssetAmount', signedOrders[0].makerAssetAmount === '0');
    printer('makerFeeAssetData', signedOrders[0].makerFeeAssetData === '0x');
    printer('makerFee', signedOrders[0].makerFee === '0');
    printer('takerAssetData', signedOrders[0].takerAssetData === '0x');
    printer('takerAssetAmount', signedOrders[0].takerAssetAmount === '0');
    printer('takerFeeAssetData', signedOrders[0].takerFeeAssetData === '0x');
    printer('takerFee', signedOrders[0].takerFee === '0');
    printer('expirationTimeSeconds', signedOrders[0].expirationTimeSeconds === '10000000000');
    printer('salt', signedOrders[0].salt === '1532559225');
    printer('signature', signedOrders[0].signature === '0x');

    printer = prettyPrintTestCase('signedOrderTest', 'NonNullAssetData');
    printer('chainId', signedOrders[1].chainId === 1337);
    printer('makerAddress', signedOrders[1].makerAddress === hexUtils.leftPad('0x1', 20));
    printer('takerAddress', signedOrders[1].takerAddress === hexUtils.leftPad('0x2', 20));
    printer('senderAddress', signedOrders[1].senderAddress === hexUtils.leftPad('0x3', 20));
    printer('feeRecipientAddress', signedOrders[1].feeRecipientAddress === hexUtils.leftPad('0x4', 20));
    printer('exchangeAddress', signedOrders[1].exchangeAddress === hexUtils.leftPad('0x5', 20));
    printer(
        'makerAssetData',
        signedOrders[1].makerAssetData === '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer('makerAssetAmount', signedOrders[1].makerAssetAmount === '123456789');
    printer(
        'makerFeeAssetData',
        signedOrders[1].makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('makerFee', signedOrders[1].makerFee === '89');
    printer(
        'takerAssetData',
        signedOrders[1].takerAssetData === '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer('takerAssetAmount', signedOrders[1].takerAssetAmount === '987654321');
    printer(
        'takerFeeAssetData',
        signedOrders[1].takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('takerFee', signedOrders[1].takerFee === '12');
    printer('expirationTimeSeconds', signedOrders[1].expirationTimeSeconds === '10000000000');
    printer('salt', signedOrders[1].salt === '1532559225');
    printer(
        'signature',
        signedOrders[1].signature === '0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33',
    );
}

function testStats(stats: WrapperStats[]): void {
    const printer = prettyPrintTestCase('statsTest', 'realisticStats');
    printer('version', stats[0].version === 'development');
    printer('pubSubTopic', stats[0].pubSubTopic === 'someTopic');
    printer('rendezvous', stats[0].rendezvous === '/0x-mesh/network/1337/version/2');
    printer('peerID', stats[0].peerID === '16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7');
    printer('ethereumChainID', stats[0].ethereumChainID === 1337);
    printer('latestBlock | hash', stats[0].latestBlock.hash === hexUtils.leftPad('0x1', 32));
    printer('latestBlock | number', stats[0].latestBlock.number === 1500);
    printer('numOrders', stats[0].numOrders === 100000);
    printer('numPeers', stats[0].numPeers === 200);
    printer('numOrdersIncludingRemoved', stats[0].numOrdersIncludingRemoved === 200000);
    printer('numPinnedOrders', stats[0].numPinnedOrders === 400);
    printer(
        'maxExpirationTime',
        stats[0].maxExpirationTime === '115792089237316195423570985008687907853269984665640564039457584007913129639935',
    );
    printer('startOfCurrentUTCDay', stats[0].startOfCurrentUTCDay === '2006-01-01 00:00:00 +0000 UTC');
    printer('ethRPCRequestsSentInCurrentUTCDay', stats[0].ethRPCRequestsSentInCurrentUTCDay === 100000);
    printer('ethRPCRateLimitExpiredRequests', stats[0].ethRPCRateLimitExpiredRequests === 5000);
}

function testValidationResults(validationResults: WrapperValidationResults[]): void {
    let printer = prettyPrintTestCase('validationResultsTest', 'emptyValidationResults');
    printer('accepted | length', validationResults[0].accepted.length === 0);
    printer('rejected | length', validationResults[0].rejected.length === 0);

    printer = prettyPrintTestCase('validationResultsTest', 'oneAcceptedResult');
    printer('accepted | length', validationResults[1].accepted.length === 1);
    printer('accepted | orderHash', validationResults[1].accepted[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('accepted | signedOrder | chainId', validationResults[1].accepted[0].signedOrder.chainId === 1337);
    printer(
        'accepted | signedOrder | makerAddress',
        validationResults[1].accepted[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'accepted | signedOrder | takerAddress',
        validationResults[1].accepted[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'accepted | signedOrder | senderAddress',
        validationResults[1].accepted[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'accepted | signedOrder | feeRecipientAddress',
        validationResults[1].accepted[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'accepted | signedOrder | exchangeAddress',
        validationResults[1].accepted[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'accepted | signedOrder | makerAssetData',
        validationResults[1].accepted[0].signedOrder.makerAssetData === '0x',
    );
    printer(
        'accepted | signedOrder | makerAssetAmount',
        validationResults[1].accepted[0].signedOrder.makerAssetAmount === '0',
    );
    printer(
        'accepted | signedOrder | makerFeeAssetData',
        validationResults[1].accepted[0].signedOrder.makerFeeAssetData === '0x',
    );
    printer('accepted | signedOrder | makerFee', validationResults[1].accepted[0].signedOrder.makerFee === '0');
    printer(
        'accepted | signedOrder | takerAssetData',
        validationResults[1].accepted[0].signedOrder.takerAssetData === '0x',
    );
    printer(
        'accepted | signedOrder | takerAssetAmount',
        validationResults[1].accepted[0].signedOrder.takerAssetAmount === '0',
    );
    printer(
        'accepted | signedOrder | takerFeeAssetData',
        validationResults[1].accepted[0].signedOrder.takerFeeAssetData === '0x',
    );
    printer('accepted | signedOrder | takerFee', validationResults[1].accepted[0].signedOrder.takerFee === '0');
    printer(
        'accepted | signedOrder | expirationTimeSeconds',
        validationResults[1].accepted[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('accepted | signedOrder | salt', validationResults[1].accepted[0].signedOrder.salt === '1532559225');
    printer('accepted | signedOrder | signature', validationResults[1].accepted[0].signedOrder.signature === '0x');
    printer('accepted | fillableTakerAssetAmount', validationResults[1].accepted[0].fillableTakerAssetAmount === '0');
    printer('accepted | isNew', validationResults[1].accepted[0].isNew);
    printer('rejected | length', validationResults[1].rejected.length === 0);

    printer = prettyPrintTestCase('validationResultsTest', 'oneRejectedResult');
    printer('accepted | length', validationResults[2].accepted.length === 0);
    printer('rejected | length', validationResults[2].rejected.length === 1);
    printer('rejected | orderHash', validationResults[2].rejected[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('rejected | signedOrder | chainId', validationResults[2].rejected[0].signedOrder.chainId === 1337);
    printer(
        'rejected | signedOrder | makerAddress',
        validationResults[2].rejected[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'rejected | signedOrder | takerAddress',
        validationResults[2].rejected[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'rejected | signedOrder | senderAddress',
        validationResults[2].rejected[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'rejected | signedOrder | feeRecipientAddress',
        validationResults[2].rejected[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'rejected | signedOrder | exchangeAddress',
        validationResults[2].rejected[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'rejected | signedOrder | makerAssetData',
        validationResults[2].rejected[0].signedOrder.makerAssetData === '0x',
    );
    printer(
        'rejected | signedOrder | makerAssetAmount',
        validationResults[2].rejected[0].signedOrder.makerAssetAmount === '0',
    );
    printer(
        'rejected | signedOrder | makerFeeAssetData',
        validationResults[2].rejected[0].signedOrder.makerFeeAssetData === '0x',
    );
    printer('rejected | signedOrder | makerFee', validationResults[1].accepted[0].signedOrder.makerFee === '0');
    printer(
        'rejected | signedOrder | takerAssetData',
        validationResults[2].rejected[0].signedOrder.takerAssetData === '0x',
    );
    printer(
        'rejected | signedOrder | takerAssetAmount',
        validationResults[2].rejected[0].signedOrder.takerAssetAmount === '0',
    );
    printer(
        'rejected | signedOrder | takerFeeAssetData',
        validationResults[2].rejected[0].signedOrder.takerFeeAssetData === '0x',
    );
    printer('rejected | signedOrder | takerFee', validationResults[1].accepted[0].signedOrder.takerFee === '0');
    printer(
        'rejected | signedOrder | expirationTimeSeconds',
        validationResults[2].rejected[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('rejected | signedOrder | salt', validationResults[2].rejected[0].signedOrder.salt === '1532559225');
    printer('rejected | signedOrder | signature', validationResults[2].rejected[0].signedOrder.signature === '0x');
    printer('rejected | kind', validationResults[2].rejected[0].kind === 'ZEROEX_VALIDATION');
    printer(
        'rejected | status | code',
        validationResults[2].rejected[0].status.code === 'OrderHasInvalidMakerAssetData',
    );
    printer(
        'rejected | status | message',
        validationResults[2].rejected[0].status.message ===
            'order makerAssetData must encode a supported assetData type',
    );

    printer = prettyPrintTestCase('validationResultsTest', 'realisticValidationResults');
    // Accepted 1
    printer('accepted | length', validationResults[3].accepted.length === 2);
    printer('accepted | orderHash', validationResults[3].accepted[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('accepted | signedOrder | chainId', validationResults[3].accepted[0].signedOrder.chainId === 1337);
    printer(
        'accepted | signedOrder | makerAddress',
        validationResults[3].accepted[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'accepted | signedOrder | takerAddress',
        validationResults[3].accepted[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'accepted | signedOrder | senderAddress',
        validationResults[3].accepted[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'accepted | signedOrder | feeRecipientAddress',
        validationResults[3].accepted[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'accepted | signedOrder | exchangeAddress',
        validationResults[3].accepted[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'accepted | signedOrder | makerAssetData',
        validationResults[3].accepted[0].signedOrder.makerAssetData === '0x',
    );
    printer(
        'accepted | signedOrder | makerAssetAmount',
        validationResults[3].accepted[0].signedOrder.makerAssetAmount === '0',
    );
    printer(
        'accepted | signedOrder | makerFeeAssetData',
        validationResults[3].accepted[0].signedOrder.makerFeeAssetData === '0x',
    );
    printer('accepted | signedOrder | makerFee', validationResults[3].accepted[0].signedOrder.makerFee === '0');
    printer(
        'accepted | signedOrder | takerAssetData',
        validationResults[3].accepted[0].signedOrder.takerAssetData === '0x',
    );
    printer(
        'accepted | signedOrder | takerAssetAmount',
        validationResults[3].accepted[0].signedOrder.takerAssetAmount === '0',
    );
    printer(
        'accepted | signedOrder | takerFeeAssetData',
        validationResults[3].accepted[0].signedOrder.takerFeeAssetData === '0x',
    );
    printer('accepted | signedOrder | takerFee', validationResults[3].accepted[0].signedOrder.takerFee === '0');
    printer(
        'accepted | signedOrder | expirationTimeSeconds',
        validationResults[3].accepted[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('accepted | signedOrder | salt', validationResults[3].accepted[0].signedOrder.salt === '1532559225');
    printer('accepted | signedOrder | signature', validationResults[3].accepted[0].signedOrder.signature === '0x');
    printer('accepted | fillableTakerAssetAmount', validationResults[3].accepted[0].fillableTakerAssetAmount === '0');
    printer('accepted | isNew', validationResults[3].accepted[0].isNew);
    // Accepted 2
    printer('accepted | orderHash', validationResults[3].accepted[1].orderHash === hexUtils.leftPad('0x1', 32));
    printer('accepted | signedOrder | chainId', validationResults[3].accepted[1].signedOrder.chainId === 1337);
    printer(
        'accepted | signedOrder | makerAddress',
        validationResults[3].accepted[1].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'accepted | signedOrder | takerAddress',
        validationResults[3].accepted[1].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'accepted | signedOrder | senderAddress',
        validationResults[3].accepted[1].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'accepted | signedOrder | feeRecipientAddress',
        validationResults[3].accepted[1].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'accepted | signedOrder | exchangeAddress',
        validationResults[3].accepted[1].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'accepted | signedOrder | makerAssetData',
        validationResults[3].accepted[1].signedOrder.makerAssetData ===
            '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer(
        'accepted | signedOrder | makerAssetAmount',
        validationResults[3].accepted[1].signedOrder.makerAssetAmount === '123456789',
    );
    printer(
        'accepted | signedOrder | makerFeeAssetData',
        validationResults[3].accepted[1].signedOrder.makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('accepted | signedOrder | makerFee', validationResults[3].accepted[1].signedOrder.makerFee === '89');
    printer(
        'accepted | signedOrder | takerAssetData',
        validationResults[3].accepted[1].signedOrder.takerAssetData ===
            '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer(
        'accepted | signedOrder | takerAssetAmount',
        validationResults[3].accepted[1].signedOrder.takerAssetAmount === '987654321',
    );
    printer(
        'accepted | signedOrder | takerFeeAssetData',
        validationResults[3].accepted[1].signedOrder.takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('accepted | signedOrder | takerFee', validationResults[3].accepted[1].signedOrder.takerFee === '12');
    printer(
        'accepted | signedOrder | expirationTimeSeconds',
        validationResults[3].accepted[1].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('accepted | signedOrder | salt', validationResults[3].accepted[1].signedOrder.salt === '1532559225');
    printer(
        'accepted | signedOrder | signature',
        validationResults[3].accepted[1].signedOrder.signature ===
            '0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33',
    );
    printer(
        'accepted | fillableTakerAssetAmount',
        validationResults[3].accepted[1].fillableTakerAssetAmount === '987654321',
    );
    printer('accepted | isNew', validationResults[3].accepted[1].isNew);
    // Rejected 1
    printer('rejected | length', validationResults[3].rejected.length === 1);
    printer('rejected | orderHash', validationResults[3].rejected[0].orderHash === hexUtils.leftPad('0x1', 32));
    printer('rejected | signedOrder | chainId', validationResults[3].rejected[0].signedOrder.chainId === 1337);
    printer(
        'rejected | signedOrder | makerAddress',
        validationResults[3].rejected[0].signedOrder.makerAddress === hexUtils.leftPad('0x1', 20),
    );
    printer(
        'rejected | signedOrder | takerAddress',
        validationResults[3].rejected[0].signedOrder.takerAddress === hexUtils.leftPad('0x2', 20),
    );
    printer(
        'rejected | signedOrder | senderAddress',
        validationResults[3].rejected[0].signedOrder.senderAddress === hexUtils.leftPad('0x3', 20),
    );
    printer(
        'rejected | signedOrder | feeRecipientAddress',
        validationResults[3].rejected[0].signedOrder.feeRecipientAddress === hexUtils.leftPad('0x4', 20),
    );
    printer(
        'rejected | signedOrder | exchangeAddress',
        validationResults[3].rejected[0].signedOrder.exchangeAddress === hexUtils.leftPad('0x5', 20),
    );
    printer(
        'rejected | signedOrder | makerAssetData',
        validationResults[3].rejected[0].signedOrder.makerAssetData ===
            '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
    );
    printer(
        'rejected | signedOrder | makerAssetAmount',
        validationResults[3].rejected[0].signedOrder.makerAssetAmount === '123456789',
    );
    printer(
        'rejected | signedOrder | makerFeeAssetData',
        validationResults[3].rejected[0].signedOrder.makerFeeAssetData ===
            '0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064',
    );
    printer('rejected | signedOrder | makerFee', validationResults[3].rejected[0].signedOrder.makerFee === '89');
    printer(
        'rejected | signedOrder | takerAssetData',
        validationResults[3].rejected[0].signedOrder.takerAssetData ===
            '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    );
    printer(
        'rejected | signedOrder | takerAssetAmount',
        validationResults[3].rejected[0].signedOrder.takerAssetAmount === '987654321',
    );
    printer(
        'rejected | signedOrder | takerFeeAssetData',
        validationResults[3].rejected[0].signedOrder.takerFeeAssetData ===
            '0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3',
    );
    printer('rejected | signedOrder | takerFee', validationResults[3].rejected[0].signedOrder.takerFee === '12');
    printer(
        'rejected | signedOrder | expirationTimeSeconds',
        validationResults[3].rejected[0].signedOrder.expirationTimeSeconds === '10000000000',
    );
    printer('rejected | signedOrder | salt', validationResults[3].rejected[0].signedOrder.salt === '1532559225');
    printer(
        'rejected | signedOrder | signature',
        validationResults[3].rejected[0].signedOrder.signature ===
            '0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33',
    );
    printer('rejected | kind', validationResults[3].rejected[0].kind === 'MESH_ERROR');
    printer('rejected | status | code', validationResults[3].rejected[0].status.code === 'CoordinatorEndpointNotFound');
    printer(
        'rejected | status | message',
        validationResults[3].rejected[0].status.message ===
            'corresponding coordinator endpoint not found in CoordinatorRegistry contract',
    );
}

// tslint:enable:no-console
// tslint:enable:custom-no-magic-numbers
/*********************** Utils ***********************/

async function waitForLoadAsync(): Promise<void> {
    // Note: this approach is not CPU efficient but it avoids race
    // conditions and has the advantage of returning instantaneously if the
    // Wasm code has already loaded.
    while (!isWasmLoaded) {
        await sleepAsync(wasmLoadCheckIntervalMs);
    }
}

async function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
// tslint:disable-line:max-file-line-count
