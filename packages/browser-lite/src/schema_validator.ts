import * as ajv from 'ajv';

import { JsonSchema } from './types';

export interface SchemaValidationResult {
    success: boolean;
    errors: string[];
    fatal?: string;
}

export interface SchemaValidator {
    orderValidator: (input: string) => SchemaValidationResult;
    messageValidator: (input: string) => SchemaValidationResult;
}

export function getSchemaValidator(
    chainId: number,
    exchangeAddress: string,
    customOrderFilter?: JsonSchema,
): SchemaValidator {
    const chainIdSchema = {
        $id: 'http://example.com/chainId',
        const: chainId,
    };

    // FIXME(jalextowle): Add the checksummed address to this.
    const exchangeAddressSchema = {
        $id: 'http://example.com/exchangeAddress',
        enum: [exchangeAddress],
    };

    const AJV = new ajv({
        schemas: [
            {
                ...customOrderFilter,
                $id: 'http://example.com/customOrder',
            },
            addressSchema,
            wholeNumberSchema,
            hexSchema,
            chainIdSchema,
            exchangeAddressSchema,
            orderSchema,
            signedOrderSchema,
            rootOrderSchema,
            rootOrderMessageSchema,
        ],
    });
    const orderValidate = AJV.getSchema('http://example.com/rootOrder');
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrder" schema in AJV');
    }
    const messageValidate = AJV.getSchema('http://example.com/rootOrderMessage');
    if (messageValidate === undefined) {
        throw new Error('Cannot find "rootOrderMessage" schema in AJV');
    }
    return {
        orderValidator: constructValidationFunctionWrapper(orderValidate),
        messageValidator: constructValidationFunctionWrapper(messageValidate),
    };
}

function constructValidationFunctionWrapper(
    validationFunction: ajv.ValidateFunction,
): (input: string) => SchemaValidationResult {
    return (input: string) => {
        const result: any = { success: false, errors: [] };
        try {
            result.success = validationFunction(JSON.parse(input));
            if (validationFunction.errors) {
                result.errors = validationFunction.errors.map(error => JSON.stringify(error));
            }
        } catch (error) {
            result.fatal = JSON.stringify(error);
        }
        return result;
    };
}

const addressSchema = {
    $id: 'http://example.com/address',
    type: 'string',
    pattern: '^0x[0-9a-fA-F]{40}$',
};
const wholeNumberSchema = {
    $id: 'http://example.com/wholeNumber',
    anyOf: [{ type: 'string', pattern: '^\\d+$' }, { type: 'integer' }],
};
const hexSchema = { $id: 'http://example.com/hex', type: 'string', pattern: '^0x(([0-9a-fA-F][0-9a-fA-F])+)?$' };
const orderSchema = {
    $id: 'http://example.com/order',
    properties: {
        makerAddress: { $ref: 'http://example.com/address' },
        takerAddress: { $ref: 'http://example.com/address' },
        makerFee: { $ref: 'http://example.com/wholeNumber' },
        takerFee: { $ref: 'http://example.com/wholeNumber' },
        senderAddress: { $ref: 'http://example.com/address' },
        makerAssetAmount: { $ref: 'http://example.com/wholeNumber' },
        takerAssetAmount: { $ref: 'http://example.com/wholeNumber' },
        makerAssetData: { $ref: 'http://example.com/hex' },
        takerAssetData: { $ref: 'http://example.com/hex' },
        makerFeeAssetData: { $ref: 'http://example.com/hex' },
        takerFeeAssetData: { $ref: 'http://example.com/hex' },
        salt: { $ref: 'http://example.com/wholeNumber' },
        feeRecipientAddress: { $ref: 'http://example.com/address' },
        expirationTimeSeconds: { $ref: 'http://example.com/wholeNumber' },
        exchangeAddress: { $ref: 'http://example.com/exchangeAddress' },
        chainId: { $ref: 'http://example.com/chainId' },
    },
    required: [
        'makerAddress',
        'takerAddress',
        'makerFee',
        'takerFee',
        'senderAddress',
        'makerAssetAmount',
        'takerAssetAmount',
        'makerAssetData',
        'takerAssetData',
        'makerFeeAssetData',
        'takerFeeAssetData',
        'salt',
        'feeRecipientAddress',
        'expirationTimeSeconds',
        'exchangeAddress',
        'chainId',
    ],
    type: 'object',
};
const signedOrderSchema = {
    $id: 'http://example.com/signedOrder',
    allOf: [
        { $ref: 'http://example.com/order' },
        { properties: { signature: { $ref: 'http://example.com/hex' } }, required: ['signature'] },
    ],
};
const rootOrderSchema = {
    $id: 'http://example.com/rootOrder',
    allOf: [{ $ref: 'http://example.com/customOrder' }, { $ref: 'http://example.com/signedOrder' }],
};
const rootOrderMessageSchema = {
    $id: 'http://example.com/rootOrderMessage',
    properties: {
        messageType: { type: 'string', pattern: 'order' },
        order: { $ref: 'http://example.com/rootOrder' },
        topics: { type: 'array', minItems: 1, items: { type: 'string' } },
    },
    required: ['messageType', 'order', 'topics'],
};
