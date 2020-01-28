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
import { wasmBuffer } from '../ts/generated/test_wasm_buffer';
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
} from '../ts/types';
import '../ts/wasm_exec';

interface ConversionTestCase {
    contractEventsAsync: () => ContractEvent[];
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

const NULL_BYTES = '0x0000000000000000000000000000000000000000000000000000000000000000';

(async () => {
    await waitForLoadAsync();
    const contractEvents = await conversionTestCases.contractEventsAsync();
    testContractEvents(contractEvents);
})();

function testContractEvents(contractEvents: ContractEvent[]): void {
    // ERC20ApprovalEvent
    let printer = prettyPrintTestCase('contractEventTest', 0);
    testContractEventPrelude(printer, contractEvents[0]);
    printer('kind', contractEvents[0].kind === 'ERC20ApprovalEvent');
    const erc20ApprovalParams = contractEvents[0].parameters as ERC20ApprovalEvent;
    printer('parameter | owner', erc20ApprovalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | spender', erc20ApprovalParams.spender === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', new BigNumber(1000).isEqualTo(erc20ApprovalParams.value));

    // ERC20TransferEvent
    printer = prettyPrintTestCase('contractEventTest', 1);
    testContractEventPrelude(printer, contractEvents[1]);
    printer('kind', contractEvents[1].kind === 'ERC20TransferEvent');
    const erc20TransferParams = contractEvents[1].parameters as ERC20TransferEvent;
    printer('parameter | from', erc20TransferParams.from === hexUtils.leftPad('0x4', 20));
    printer('parameter | to', erc20TransferParams.to === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', new BigNumber(1000).isEqualTo(erc20TransferParams.value));

    // ERC721ApprovalEvent
    printer = prettyPrintTestCase('contractEventTest', 2);
    testContractEventPrelude(printer, contractEvents[2]);
    printer('kind', contractEvents[2].kind === 'ERC721ApprovalEvent');
    const erc721ApprovalParams = contractEvents[2].parameters as ERC721ApprovalEvent;
    printer('parameter | owner', erc721ApprovalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | approved', erc721ApprovalParams.approved === hexUtils.leftPad('0x5', 20));
    printer('parameter | tokenId', new BigNumber(1).isEqualTo(erc721ApprovalParams.tokenId));

    // ERC721ApprovalForAllEvent
    printer = prettyPrintTestCase('contractEventTest', 3);
    testContractEventPrelude(printer, contractEvents[3]);
    printer('kind', contractEvents[3].kind === 'ERC721ApprovalForAllEvent');
    const erc721ApprovalForAllParams = contractEvents[3].parameters as ERC721ApprovalForAllEvent;
    printer('parameter | owner', erc721ApprovalForAllParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | operator', erc721ApprovalForAllParams.operator === hexUtils.leftPad('0x5', 20));
    printer('parameter | approved', erc721ApprovalForAllParams.approved);

    // ERC721TransferEvent
    printer = prettyPrintTestCase('contractEventTest', 4);
    testContractEventPrelude(printer, contractEvents[4]);
    printer('kind', contractEvents[4].kind === 'ERC721TransferEvent');
    const erc721TransferParams = contractEvents[4].parameters as ERC721TransferEvent;
    printer('parameter | from', erc721TransferParams.from === hexUtils.leftPad('0x4', 20));
    printer('parameter | to', erc721TransferParams.to === hexUtils.leftPad('0x5', 20));
    printer('parameter | tokenId', new BigNumber(1).isEqualTo(erc721TransferParams.tokenId));

    // ERC1155ApprovalForAllEvent
    printer = prettyPrintTestCase('contractEventTest', 5);
    testContractEventPrelude(printer, contractEvents[5]);
    printer('kind', contractEvents[5].kind === 'ERC1155ApprovalForAllEvent');
    const erc1155ApprovalForAllParams = contractEvents[5].parameters as ERC1155ApprovalForAllEvent;
    printer('parameter | owner', erc1155ApprovalForAllParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | operator', erc1155ApprovalForAllParams.operator === hexUtils.leftPad('0x5', 20));
    printer('parameter | approved', !erc1155ApprovalForAllParams.approved);

    // ERC1155TransferSingleEvent
    printer = prettyPrintTestCase('contractEventTest', 6);
    testContractEventPrelude(printer, contractEvents[6]);
    printer('kind', contractEvents[6].kind === 'ERC1155TransferSingleEvent');
    const erc1155TransferSingleParams = contractEvents[6].parameters as ERC1155TransferSingleEvent;
    printer('parameter | operator', erc1155TransferSingleParams.operator === hexUtils.leftPad('0x4', 20));
    printer('parameter | from', erc1155TransferSingleParams.from === hexUtils.leftPad('0x5', 20));
    printer('parameter | to', erc1155TransferSingleParams.to === hexUtils.leftPad('0x6', 20));
    printer('parameter | id', new BigNumber(1).isEqualTo(erc1155TransferSingleParams.id));
    printer('parameter | value', new BigNumber(100).isEqualTo(erc1155TransferSingleParams.value));

    // ERC1155TransferBatchEvent
    printer = prettyPrintTestCase('contractEventTest', 7);
    testContractEventPrelude(printer, contractEvents[7]);
    printer('kind', contractEvents[7].kind === 'ERC1155TransferBatchEvent');
    const erc1155TransferBatchParams = contractEvents[7].parameters as ERC1155TransferBatchEvent;
    printer('parameter | operator', erc1155TransferBatchParams.operator === hexUtils.leftPad('0x4', 20));
    printer('parameter | from', erc1155TransferBatchParams.from === hexUtils.leftPad('0x5', 20));
    printer('parameter | to', erc1155TransferBatchParams.to === hexUtils.leftPad('0x6', 20));
    printer(
        'parameter | ids',
        erc1155TransferBatchParams.ids.length === 1 && new BigNumber(1).isEqualTo(erc1155TransferBatchParams.ids[0]),
    );
    printer(
        'parameter | values',
        erc1155TransferBatchParams.values.length === 1 &&
            new BigNumber(100).isEqualTo(erc1155TransferBatchParams.values[0]),
    );
}

function testContractEventPrelude(
    printer: (section: string, value: boolean) => void,
    contractEvent: ContractEvent,
): void {
    printer('blockHash', contractEvent.blockHash === hexUtils.leftPad(1, 32));
    printer('txHash', contractEvent.txHash === hexUtils.leftPad(2, 32));
    printer('txIndex', contractEvent.txIndex === 123);
    printer('logIndex', contractEvent.logIndex === 321);
    printer('isRemoved', !contractEvent.isRemoved);
    printer('address', contractEvent.address === hexUtils.leftPad(3, 20));
}

function prettyPrintTestCase(name: string, idx: number): (section: string, value: boolean) => void {
    return (section: string, value: boolean) => {
        console.log(`(${name} | ${idx} | ${section}): ${value}`);
    };
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
