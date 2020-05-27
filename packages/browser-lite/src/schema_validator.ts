import * as ajv from 'ajv';

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
 * Sets up a schema validator on the window object.
 * @param customOrderSchema The custom filter that will be used by the root
 *        schemas.
 * @param schemas These are all of the schemas that can be compiled before the
 *        customOrderSchema.
 * @param rootSchemas The root schemas. These must be compiled last.
 */
export function setSchemaValidator(customOrderSchemaString: string, schemas: string[], rootSchemas: string[]): void {
    // NOTE(jalextowle): This causes `ajv` to look at `id` fields rather than
    // `$id` fields.
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

    const orderValidate = AJV.getSchema('/rootOrder');
    if (orderValidate === undefined) {
        throw new Error('Cannot find "/rootOrder" schema in AJV');
    }
    const messageValidate = AJV.getSchema('/rootOrderMessage');
    if (messageValidate === undefined) {
        throw new Error('Cannot find "rootOrderMessage" schema in AJV');
    }
    (window as any).schemaValidator = {
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
