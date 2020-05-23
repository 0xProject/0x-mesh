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

/**
 * Constructs a `SchemaValidator` object that is compatible with the `orderfilter` package.
 * @param chainId The chainId of the current ethereum network
 * @param exchangeAddress The exchange address that Mesh should recognize on the current network
 * @param customOrderFilter An extension of the order schema that adds specific requirements to
 *        orders.
 */
export function getSchemaValidator(
    chainId: number,
    exchangeAddress: string,
    customOrderFilter?: JsonSchema,
): SchemaValidator {
    const chainIdSchema = {
        $async: false,
        $id: '/chainId',
        const: chainId,
    };

    // NOTE(jalextowle): The current implementation of `zeroex.SignedOrder.MarhsalJSON`
    // lowercases all addresses. This eliminates the need to add the checksummed address.
    const exchangeAddressSchema = {
        $async: false,
        $id: '/exchangeAddress',
        enum: [exchangeAddress],
    };

    const AJV = new ajv();
    AJV.addSchema(addressSchema);
    AJV.addSchema(wholeNumberSchema);
    AJV.addSchema(hexSchema);
    AJV.addSchema(chainIdSchema);
    AJV.addSchema(exchangeAddressSchema);
    AJV.addSchema({
        ...customOrderFilter,
        $async: false,
        $id: '/customOrder',
    });
    AJV.addSchema(orderSchema);
    AJV.addSchema(signedOrderSchema);
    AJV.addSchema(rootOrderSchema);
    AJV.addSchema(rootOrderMessageSchema);

    const orderValidate = AJV.getSchema('/rootOrder');
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrder" schema in AJV');
    }
    const messageValidate = AJV.getSchema('/rootOrderMessage');
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
    $async: false,
    $id: '/address',
    type: 'string',
    pattern: '^0x[0-9a-fA-F]{40}$',
};
const wholeNumberSchema = {
    $async: false,
    $id: '/wholeNumber',
    anyOf: [{ type: 'string', pattern: '^\\d+$' }, { type: 'integer' }],
};
const hexSchema = {
    $async: false,
    $id: '/hex',
    type: 'string',
    pattern: '^0x(([0-9a-fA-F][0-9a-fA-F])+)?$',
};
const orderSchema = {
    $async: false,
    $id: '/order',
    properties: {
        makerAddress: { $ref: '/address' },
        takerAddress: { $ref: '/address' },
        makerFee: { $ref: '/wholeNumber' },
        takerFee: { $ref: '/wholeNumber' },
        senderAddress: { $ref: '/address' },
        makerAssetAmount: { $ref: '/wholeNumber' },
        takerAssetAmount: { $ref: '/wholeNumber' },
        makerAssetData: { $ref: '/hex' },
        takerAssetData: { $ref: '/hex' },
        makerFeeAssetData: { $ref: '/hex' },
        takerFeeAssetData: { $ref: '/hex' },
        salt: { $ref: '/wholeNumber' },
        feeRecipientAddress: { $ref: '/address' },
        expirationTimeSeconds: { $ref: '/wholeNumber' },
        exchangeAddress: { $ref: '/exchangeAddress' },
        chainId: { $ref: '/chainId' },
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
    $async: false,
    $id: '/signedOrder',
    allOf: [{ $ref: '/order' }, { properties: { signature: { $ref: '/hex' } }, required: ['signature'] }],
};
const rootOrderSchema = {
    $async: false,
    $id: '/rootOrder',
    allOf: [{ $ref: '/customOrder' }, { $ref: '/signedOrder' }],
};
const rootOrderMessageSchema = {
    $async: false,
    $id: '/rootOrderMessage',
    properties: {
        messageType: { type: 'string', pattern: 'order' },
        order: { $ref: '/rootOrder' },
        topics: { type: 'array', minItems: 1, items: { type: 'string' } },
    },
    required: ['messageType', 'order', 'topics'],
};
