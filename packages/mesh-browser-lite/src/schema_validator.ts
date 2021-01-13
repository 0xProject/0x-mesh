/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
import * as ajv from 'ajv';

interface SynchronousValidationFunction {
    (data: any): boolean;
    schema?: object | boolean;
    errors?: null | ajv.ErrorObject[];
}

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
        AJV.addSchema(JSON.parse(schema));
    }
    const customOrderSchema = JSON.parse(customOrderSchemaString);
    AJV.addSchema({
        ...customOrderSchema,
        $id: '/customOrder',
    });
    for (const schema of rootSchemas) {
        AJV.addSchema(JSON.parse(schema));
    }

    const orderValidate = AJV.getSchema('/rootOrderV3');
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrderV3" schema in AJV');
    }
    const messageValidate = AJV.getSchema('/rootOrderMessageV3');
    if (messageValidate === undefined) {
        throw new Error('Cannot find "rootOrderMessageV3" schema in AJV');
    }
    return {
        orderValidator: constructValidationFunctionWrapper(orderValidate as SynchronousValidationFunction),
        messageValidator: constructValidationFunctionWrapper(messageValidate as SynchronousValidationFunction),
    };
}

function constructValidationFunctionWrapper(
    validationFunction: SynchronousValidationFunction,
): (input: string) => SchemaValidationResult {
    return (input: string) => {
        const result: SchemaValidationResult = { success: false, errors: [] };
        try {
            result.success = validationFunction(JSON.parse(input));
            if (validationFunction.errors) {
                result.errors = validationFunction.errors.map((error) => JSON.stringify(error));
            }
        } catch (error) {
            result.fatal = JSON.stringify(error);
        }
        return result;
    };
}
