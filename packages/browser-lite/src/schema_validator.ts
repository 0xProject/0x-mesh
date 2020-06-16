/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
import * as ajv from 'ajv';

interface AsynchronousValidationFunction {
    (data: any): Promise<boolean>;
    schema?: object | boolean;
    errors?: null | ajv.ErrorObject[];
}

export interface SchemaValidationResult {
    success: boolean;
    errors: string[];
}

export interface SchemaValidator {
    orderValidator: (input: string) => Promise<SchemaValidationResult>;
    messageValidator: (input: string) => Promise<SchemaValidationResult>;
}

/**
 * Creates a schema validator object for the provided schemas.
 * @param customOrderSchema The custom filter that will be used by the root
 *        schemas.
 * @param schemas These are all of the schemas that can be compiled before the
 *        customOrderSchema.
 * @param rootSchemas The root schemas. These must be compiled last.
 */
export function createSchemaValidator(
    customOrderSchemaString: string,
    schemas: string[],
    rootSchemas: string[],
): SchemaValidator {
    const AJV = new ajv();
    for (const schema of schemas) {
        AJV.addSchema({
            ...JSON.parse(schema),
            async: true,
        });
    }
    const customOrderSchema = JSON.parse(customOrderSchemaString);
    AJV.addSchema({
        ...customOrderSchema,
        async: true,
        $id: '/customOrder',
    });
    for (const schema of rootSchemas) {
        AJV.addSchema({
            ...JSON.parse(schema),
            async: true,
        });
    }

    const orderValidate = AJV.getSchema('/rootOrder');
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrder" schema in AJV');
    }
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrder" schema in AJV');
    }
    const messageValidate = AJV.getSchema('/rootOrderMessage');
    if (messageValidate === undefined) {
        throw new Error('Cannot find "rootOrderMessage" schema in AJV');
    }
    return {
        orderValidator: constructValidationFunctionWrapper(orderValidate as AsynchronousValidationFunction),
        messageValidator: constructValidationFunctionWrapper(messageValidate as AsynchronousValidationFunction),
    };
}

function constructValidationFunctionWrapper(
    validationFunction: AsynchronousValidationFunction,
): (input: string) => Promise<SchemaValidationResult> {
    return async (input: string) => {
        const result: SchemaValidationResult = { success: false, errors: [] };
        result.success = await validationFunction(JSON.parse(input));
        if (validationFunction.errors) {
            result.errors = validationFunction.errors.map(error => JSON.stringify(error));
        }
        return result;
    };
}
