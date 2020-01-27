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
import { ContractEvent, ERC20ApprovalEvent, ERC20TransferEvent, WrapperContractEvent } from '../ts/types';
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
    const approvalParams = contractEvents[0].parameters as ERC20ApprovalEvent;
    printer('parameter | owner', approvalParams.owner === hexUtils.leftPad('0x4', 20));
    printer('parameter | spender', approvalParams.spender === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', new BigNumber(1000).isEqualTo(approvalParams.value));

    // ERC20TransferEvent
    printer = prettyPrintTestCase('contractEventTest', 1);
    testContractEventPrelude(printer, contractEvents[1]);
    printer('kind', contractEvents[1].kind === 'ERC20TransferEvent');
    const transferParams = contractEvents[1].parameters as ERC20TransferEvent;
    printer('parameter | from', transferParams.from === hexUtils.leftPad('0x4', 20));
    printer('parameter | to', transferParams.to === hexUtils.leftPad('0x5', 20));
    printer('parameter | value', new BigNumber(1000).isEqualTo(transferParams.value));
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
